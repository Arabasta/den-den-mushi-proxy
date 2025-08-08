package pty_token

import (
	request2 "den-den-mushi-Go/internal/control/app/pty_token/request"
	oapi "den-den-mushi-Go/openapi/control"
	dtopkg "den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/httpx"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/middleware/wrapper"
	"den-den-mushi-Go/pkg/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}

func (h *Handler) PostApiV1PtyTokenStart(c *gin.Context) {
	var raw oapi.StartRequest
	if err := c.ShouldBindJSON(&raw); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "Failed to bind JSON", err, h.Log)
		return
	}
	h.Log.Debug("Start request received", zap.Any("request", raw))

	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return
	}

	r := wrapper.WithAuth[request2.StartRequest]{
		Body: request2.StartRequest{
			Purpose: types.ConnectionPurpose(raw.Purpose),
			ChangeID: func() string {
				if raw.ChangeId != nil {
					return *raw.ChangeId
				}
				return ""
			}(),
			Server: dtopkg.ServerInfo{
				OSUser: raw.Server.OsUser,
				IP:     raw.Server.ServerIp,
			},
		},
		AuthCtx: *authCtx,
	}

	token, proxy, err := h.Service.mintStartToken(r)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "Failed to mint start token. Reason: "+err.Error(), err, h.Log)
		return
	}

	h.Log.Debug("Start token minted", zap.String("jwt", token), zap.String("userID", r.AuthCtx.UserID))
	c.JSON(200, oapi.TokenResponse{
		Token:    token,
		ProxyUrl: proxy,
	})
}

func (h *Handler) PostApiV1PtyTokenJoin(c *gin.Context) {
	var raw oapi.JoinRequest
	if err := c.ShouldBindJSON(&raw); err != nil {
		httpx.RespondError(c, http.StatusBadRequest, "Failed to bind JSON", err, h.Log)
		return
	}
	h.Log.Debug("Join request received", zap.Any("request", raw))

	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return
	}

	r := wrapper.WithAuth[request2.JoinRequest]{
		Body: request2.JoinRequest{
			PtySessionId: raw.PtySessionId,
			StartRole:    types.StartRole(raw.StartRole),
		},
		AuthCtx: *authCtx,
	}

	token, proxy, err := h.Service.mintJoinToken(r)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "Failed to mint join token. Reason: "+err.Error(), err, h.Log)
		return
	}

	h.Log.Debug("Join token minted", zap.String("jwt", token), zap.String("userID", r.AuthCtx.UserID))
	c.JSON(200, oapi.TokenResponse{
		Token:    token,
		ProxyUrl: proxy,
	})
}

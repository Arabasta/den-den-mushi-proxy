package pty_sessions

import (
	oapi "den-den-mushi-Go/openapi/llm_external"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}

func (h *Handler) GetApiV1PtySessions(c *gin.Context, params oapi.GetApiV1PtySessionsParams) {
	changeRequestId := params.ChangeRequestId
	if changeRequestId == "" {
		h.Log.Warn("Missing required query parameter: change_request_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing change_request_id query parameter"})
		return
	}

	sessions, err := h.Service.FindAllByChangeRequestID(changeRequestId)
	if err != nil {
		h.Log.Error("Failed to fetch PTY sessions by change request ID",
			zap.String("change_request_id", changeRequestId),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch PTY sessions"})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

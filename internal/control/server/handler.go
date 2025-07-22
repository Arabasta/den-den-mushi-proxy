package server

import (
	"den-den-mushi-Go/internal/control/make_change"
	"den-den-mushi-Go/internal/control/pty_token"
	oapi "den-den-mushi-Go/openapi/control"
	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	PtyHandler        *pty_token.Handler
	MakeChangeHandler *make_change.Handler
}

// Forwarding methods (required by oapi.ServerInterface)

func (h *MasterHandler) PostApiV1PtyTokenStart(c *gin.Context) {
	h.PtyHandler.PostApiV1PtyTokenStart(c)
}

func (h *MasterHandler) PostApiV1PtyTokenJoin(c *gin.Context) {
	h.PtyHandler.PostApiV1PtyTokenJoin(c)
}

func (h *MasterHandler) GetApiV1ChangeRequests(c *gin.Context, params oapi.GetApiV1ChangeRequestsParams) {
	h.MakeChangeHandler.GetApiV1ChangeRequests(c, params)
}

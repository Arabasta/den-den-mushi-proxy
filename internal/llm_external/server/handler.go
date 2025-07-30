package server

import (
	"den-den-mushi-Go/internal/llm_external/pty_sessions"
	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	PtySessions *pty_sessions.Handler
}

// Forwarding methods (required by oapi.ServerInterface)

func (h *MasterHandler) GetApiV1PtySessionChangeRequestId(c *gin.Context, changeRequestId string) {
	h.PtySessions.GetApiV1PtySessionChangeRequestId(c, changeRequestId)
}

package server

import (
	"den-den-mushi-Go/internal/llm_external/pty_sessions"
	oapi "den-den-mushi-Go/openapi/llm_external"
	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	PtySessions *pty_sessions.Handler
}

// Forwarding methods (required by oapi.ServerInterface)

func (h *MasterHandler) GetApiV1PtySessions(c *gin.Context, params oapi.GetApiV1PtySessionsParams) {
	h.PtySessions.GetApiV1PtySessions(c, params)
}

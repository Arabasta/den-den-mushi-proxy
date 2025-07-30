package pty_sessions

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}

func (h *Handler) GetApiV1PtySessionChangeRequestId(c *gin.Context, changeRequestId string) {
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

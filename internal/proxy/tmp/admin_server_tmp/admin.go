package admin_server_tmp

import (
	"den-den-mushi-Go/internal/proxy/core/session_manager"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterAdminRoutes(r *gin.RouterGroup, sm *session_manager.Service, log *zap.Logger) {
	admin := r.Group("/api/v1")
	admin.GET("/pty_sessions", getSessions(sm, log))
}

func getSessions(sm *session_manager.Service, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessions := sm.GetPtySessionsTmpForAdminServer()
		log.Info("Retrieved PTY sessions", zap.Int("count", len(sessions)))
		c.JSON(200, gin.H{"sessions": sessions})
	}
}

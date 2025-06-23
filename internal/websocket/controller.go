package websocket

import (
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/pkg/token"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterWebsocketRoutes(r *gin.RouterGroup, cfg *config.Config, log *zap.Logger, svc *Service) {
	r.GET("/ws", websocketHandler(svc, cfg, log))
}

func websocketHandler(svc *Service, cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, _ := c.Get("claims")
		claims := val.(*token.Claims)

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade connection", zap.Error(err))
			return
		}

		svc.Run(c.Request.Context(), ws, claims)
	}
}

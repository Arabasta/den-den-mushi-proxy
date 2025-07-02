package websocket

import (
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/pkg/token"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterWebsocketRoutes(r *gin.RouterGroup, cfg *config.Config, log *zap.Logger, svc *Service) {
	r.GET("/ws", websocketHandler(svc, log))
}

func websocketHandler(svc *Service, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("websocketHandler called", zap.String("url", c.Request.URL.String()))

		val, _ := c.Get("claims")
		claims := val.(*token.Claims)

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade websocket connection", zap.Error(err))
			return
		}

		svc.run(c.Request.Context(), ws, claims)
	}
}

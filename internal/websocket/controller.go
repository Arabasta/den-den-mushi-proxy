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
		log.Info("websocketHandler called", zap.String("url", c.Request.URL.String()))

		val, _ := c.Get("claims")
		claims := val.(*token.Claims)

		log.Info("WebSocket claims",
			zap.String("Connection OSUser", claims.Connection.OSUser),
			zap.String("Connection ServerIP", claims.Connection.ServerIP),
			zap.String("Connection Port", claims.Connection.Port),
			zap.String("Connection Purpose", string(claims.Connection.Purpose)),
			zap.String("ConnectionType", string(claims.Connection.Type)))

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade connection", zap.Error(err))
			return
		}

		svc.Run(c.Request.Context(), ws, claims)
	}
}

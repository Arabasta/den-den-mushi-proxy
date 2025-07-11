package middleware

import (
	"den-den-mushi-Go/internal/proxy/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// todo: secure and make configurable, move to pkg
func Cors(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		method := c.Request.Method
		reqHeaders := c.GetHeader("Access-Control-Request-Headers")

		log.Debug("CORS Middleware triggered", zap.String("origin", origin), zap.String("method", method), zap.String("requestedHeaders", reqHeaders))

		c.Header("Access-Control-Allow-Origin", "*") // allow all for now
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Vary", "Origin")

		// preflight requests
		if method == http.MethodOptions {
			log.Debug("CORS preflight request handled", zap.String("path", c.FullPath()))
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

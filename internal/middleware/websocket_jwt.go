package middleware

import (
	"context"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

// todo: secrure
func WsAuthMiddleware(v *validator.Validator, cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Sec-WebSocket-Protocol")
		parts := strings.SplitN(h, ",", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(400, gin.H{"error": "missing token"})
			return
		}
		jwtStr := strings.TrimSpace(parts[1])
		claims, err := v.Validate(context.Background(), jwtStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		c.Header("Sec-WebSocket-Protocol", "jwt")
		c.Set("claims", claims)
		c.Next()
	}
}

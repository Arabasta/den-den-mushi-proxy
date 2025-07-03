package middleware

import (
	"den-den-mushi-Go/internal/proxy/jwt_service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WsJwtMiddleware(v *jwt_service.Validator, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Sec-WebSocket-Protocol")
		rawToken, err := v.ExtractTokenFromHeader(h)
		if err != nil {
			log.Error("Failed to extract JWT from WebSocket header", zap.String("header", h), zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"error": "JWT validation failed"}) // don't return error details to client
			return
		}

		token, claims, err := v.GetTokenAndClaims(rawToken)
		if err != nil {
			log.Error("Failed to parse JWT", zap.String("rawToken", rawToken), zap.Error(err))
			c.AbortWithStatusJSON(401, gin.H{"error": "JWT validation failed"})
			return
		}

		err = v.ValidateClaims(claims, token)
		if err != nil {
			log.Error("Failed to validate claims", zap.Any("claims", claims), zap.Any("token", token), zap.Error(err))
			c.AbortWithStatusJSON(401, gin.H{"error": "JWT validation failed"})
			return
		}

		log.Info("WebSocket JWT validated, setting claims in ctx", zap.String("user", claims.Subject), zap.String("jti", claims.ID))

		c.Header("Sec-WebSocket-Protocol", "jwt")
		c.Set("claims", claims)
		c.Next()
	}
}

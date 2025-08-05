package websocket_jwt

import (
	"den-den-mushi-Go/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func WsJwtMiddleware(v *Validator, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx, ok := middleware.GetAuthContext(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
			return
		}

		h := c.GetHeader("Sec-WebSocket-Protocol")
		rawToken, err := v.ExtractProxyTokenFromHeader(h)
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

		err = v.ValidateClaims(claims, token, authCtx)
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

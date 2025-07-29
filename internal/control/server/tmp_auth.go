package server

import (
	"den-den-mushi-Go/internal/control/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
)

func TmpAuth(log *zap.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Development.IsTmpAuthCookieEnabled {
			var tokenStr string
			var err error

			tokenStr = c.GetHeader("Authorization")
			if tokenStr == "" {
				log.Error("Authorization header missing")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			jwtSecret := []byte(cfg.CookieTmp.Secret)

			// validate JWT with secret
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})
			if err != nil || !token.Valid {
				log.Warn("invalid jwt token", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				subject, _ := claims[cfg.CookieTmp.UserIdKey].(string)
				c.Set("user_id", subject)
			} else {
				log.Warn("invalid jwt claims")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
		} else {
			log.Debug("TmpAuth middleware is disabled")
			c.Set("user_id", "tmp_user")
		}
		c.Next()
	}
}

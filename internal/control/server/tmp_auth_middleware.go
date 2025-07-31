package server

import (
	"den-den-mushi-Go/internal/control/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func TmpAuth(log *zap.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.Development.IsTmpAuthCookieEnabled {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				log.Error("Authorization header missing")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			// extract token from Authorization header, Bearer optional for now for backward compatibility
			var tokenStr string
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimSpace(authHeader[7:])
			} else {
				tokenStr = strings.TrimSpace(authHeader)
			}

			// validate JWT with secret
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.CookieTmp.Secret), nil
			})
			if err != nil || !token.Valid {
				log.Warn("invalid jwt token", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			// set the cl
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				subject, _ := claims[cfg.CookieTmp.UserIdKey].(string)
				c.Set("user_id", subject)

				// each user should only have 1 ou group, this is handled by the auth service
				ouGroup, _ := claims[cfg.CookieTmp.OuGroupKey].(string)
				if ouGroup == "" {
					log.Debug("ou_group claim is empty")
				}
				c.Set("ou_group", ouGroup)
			} else {
				log.Warn("invalid jwt claims")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
		} else {
			log.Debug("TmpAuth middleware is disabled")
			c.Set("user_id", "ddmtest")
			c.Set("ou_group", "ddm_L1_admin_group")
		}

		c.Next()
	}
}

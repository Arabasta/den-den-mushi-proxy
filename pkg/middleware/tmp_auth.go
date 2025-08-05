package middleware

import (
	"den-den-mushi-Go/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func TmpAuth(log *zap.Logger, cfg *config.Tmpauth) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.IsEnabled {
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
				return []byte(cfg.Secret), nil
			})
			if err != nil || !token.Valid {
				log.Warn("invalid jwt token", zap.Error(err))
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization failed. Please try logging in again."})
				return
			}

			// set the claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Warn("invalid jwt claims")
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}

			subject, _ := claims[cfg.UserIdKey].(string)
			c.Set("user_id", subject)

			// each user should only have 1 ou group, this is handled by the auth service
			ouGroup, _ := claims[cfg.OuGroupKey].(string)
			if ouGroup == "" {
				log.Debug("ou_group claim is empty")
			}
			c.Set("ou_group", ouGroup)
		} else {
			log.Debug("TmpAuth middleware is disabled")
			c.Set("user_id", "ddmtest")
			c.Set("ou_group", "ddm_L1_admin_group")
		}

		c.Next()
	}
}

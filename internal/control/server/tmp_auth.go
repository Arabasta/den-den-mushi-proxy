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
		tokenStr, err := c.Cookie(cfg.CookieTmp.Name)
		if err != nil {
			log.Error("cookie not found", zap.Error(err))
			c.Redirect(http.StatusFound, cfg.CookieTmp.Redirect)
			c.Abort()
			return
		}

		var jwtSecret = []byte(cfg.CookieTmp.Secret)

		// validate JWT with secret
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Warn("invalid jwt token", zap.Error(err))
			c.Redirect(http.StatusFound, cfg.CookieTmp.Redirect)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			subject := claims[cfg.CookieTmp.UserIdKey].(string)
			c.Set("user_id", subject)
		} else {
			log.Warn("invalid jwt claims")
			c.Redirect(http.StatusFound, cfg.CookieTmp.Redirect)
			c.Abort()
			return
		}

		// inject username into Gin context
		c.Set("user_id", "subject_from_token")
		//c.Set("ou_groups", "ou_group1", "ou_group2")
		c.Next()
	}
}

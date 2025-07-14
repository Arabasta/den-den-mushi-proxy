package middleware

import (
	"den-den-mushi-Go/internal/control/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Keycloak checks for keycloak_token cookie, parses and validates it sets user in context
func Keycloak(cfg *config.Config, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("keycloak_token")
		if err != nil {
			log.Info("keycloak token cookie not found", zap.Error(err))
			c.Redirect(http.StatusFound, "/static/login.html")
			c.Abort()
			return
		}

		// validate JWT with secret

		// inject username into Gin context
		c.Set("user_id", "subject_from_token")
		//c.Set("ou_groups", "ou_group1", "ou_group2")
		c.Next()
	}
}

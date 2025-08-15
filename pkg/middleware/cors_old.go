package middleware

import (
	"den-den-mushi-Go/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func Cors2(cfg *config.Cors, log *zap.Logger) gin.HandlerFunc {
	allowedMethods := strings.Join(cfg.AllowMethods, ",")
	allowedHeaders := strings.Join(cfg.AllowHeaders, ",")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		method := c.Request.Method
		reqHeaders := c.GetHeader("Access-Control-Request-Headers")

		log.Debug("CORS Middleware triggered", zap.String("origin", origin), zap.String("method", method), zap.String("requestedHeaders", reqHeaders))

		c.Header("Access-Control-Allow-Origin", "*") // allow all for now
		c.Header("Access-Control-Allow-Methods", allowedMethods)
		c.Header("Access-Control-Allow-Headers", allowedHeaders)
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(cfg.AllowCredentials))
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

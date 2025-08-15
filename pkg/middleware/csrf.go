package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func CsrfGuardNoCookies(log *zap.Logger) gin.HandlerFunc {
	isUnsafe := func(m string) bool {
		switch m {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			return false
		default:
			return true
		}
	}

	hasCookieHeader := func(h http.Header) bool {
		return len(h.Values("Cookie")) > 0
	}

	return func(c *gin.Context) {
		if !isUnsafe(c.Request.Method) {
			c.Next()
			return
		}

		if hasCookieHeader(c.Request.Header) {
			names := make([]string, 0, 4)
			for _, ck := range c.Request.Cookies() {
				names = append(names, ck.Name)
			}
			// don't log cookie values
			log.Warn("Blocked unsafe request with cookies",
				zap.String("method", c.Request.Method),
				zap.String("path", c.FullPath()),
				zap.Strings("cookie_names", names),
			)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}

func StripSetCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // run handlers
		h := c.Writer.Header()
		if len(h.Values("Set-Cookie")) > 0 {
			h.Del("Set-Cookie")
		}
	}
}

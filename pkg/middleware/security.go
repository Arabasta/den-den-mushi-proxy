package middleware

import (
	"github.com/gin-gonic/gin"
	"regexp"
)

var allowedOriginPattern = regexp.MustCompile(`^https://x([a-zA-Z0-9-]+\.)*corp\.com(:[0-9]+)?$`)

func Security(sslEnabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// have to allow iframe, even though its not secure.
		// idk why frontend is using it
		// fine grained control annoying to setup cause we don't know where the frontend will be
		// permanently hosted for now
		// have to rely on csp
		// c.Writer.Header().Set("X-Frame-Options", "Self https://*/com")

		// todo scp
		//origin := c.GetHeader("Origin")
		//if origin == "" {
		//	origin = c.GetHeader("Referer") // fallback
		//}
		//
		//frameAncestors := "'self'" // default
		//
		//// add frontend origin if match
		//if origin != "" {
		//	u, err := url.Parse(origin)
		//	if err == nil {
		//		originHost := u.Scheme + "://" + u.Host
		//		if allowedOriginPattern.MatchString(originHost) {
		//			frameAncestors = "'self' " + originHost
		//		}
		//	}
		//}

		// no sniff MIME types
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		c.Writer.Header().Set("Referrer-Policy", "no-referrer")

		//csp := "default-src 'self'; " +
		//	"script-src 'self'; " +
		//	"connect-src 'self' wss://*.corp:45007; " +
		//	"frame-ancestors " + frameAncestors
		//c.Writer.Header().Set("Content-Security-Policy", csp)

		if sslEnabled {
			c.Writer.Header().Set("Strict-Transport-Security",
				"max-age=63072000; includeSubDomains; preload")
		}

		c.Next()
	}
}

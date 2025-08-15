package middleware

import (
	"den-den-mushi-Go/pkg/config"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

func Security(cfg *config.Security, sslEnabled bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.IsEnabled {
			c.Next()
			return
		}

		// clickjack
		c.Writer.Header().Set("X-Frame-Options", "DENY")

		// no sniff MIME types
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		// standard stuff
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		c.Writer.Header().Set("Cross-Origin-Resource-Policy", "same-site")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=(), usb=(), "+
			"bluetooth=(), gyroscope=(), magnetometer=()")

		if sslEnabled {
			c.Writer.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}

		c.Next() // set CSP/cache after we know what we’re returning

		ct := c.Writer.Header().Get("Content-Type")
		reqPath := c.Request.URL.Path
		isHTML := strings.HasPrefix(ct, "text/html") || wantsHTML(c.Request.Header.Get("Accept"))

		if !isHTML {
			if strings.HasPrefix(reqPath, "/assets/") {
				if c.Writer.Header().Get("Cache-Control") == "" {
					c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				}
			} else if strings.HasPrefix(reqPath, "/api/") {
				// no CSP needed
			} else if c.Writer.Header().Get("Cache-Control") == "" && looksLikeHashedFile(reqPath) {
				c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
		} else {
			if c.Writer.Header().Get("Cache-Control") == "" {
				c.Writer.Header().Set("Cache-Control", "no-store") // don’t cache the entry HTML
			}
			// csp
			connect := "'self'"
			for _, host := range cfg.ConnectSrc {
				connect += " " + host
			}
			cspParts := []string{
				"default-src 'self'",
				"base-uri 'none'",
				"frame-ancestors 'none'",
				"object-src 'none'",
				"form-action 'self'",
				"img-src 'self' data:",
				"font-src 'self' data:",
				"script-src 'self'", // no inline, no eval
				"connect-src " + connect,
			}

			if cfg.AllowStyleAttribute {
				cspParts = append(cspParts,
					"style-src 'self'",
					"style-src-elem 'self'",
					"style-src-attr 'unsafe-inline'", // allows style="..." attributes
				)
			} else {
				cspParts = append(cspParts, "style-src 'self'")
			}

			csp := strings.Join(cspParts, "; ")

			if cfg.EnforceCsp {
				c.Writer.Header().Set("Content-Security-Policy", csp)
			} else {
				c.Writer.Header().Set("Content-Security-Policy-Report-Only", csp)
			}

		}

		c.Next()
	}
}

func wantsHTML(accept string) bool {
	return strings.Contains(accept, "text/html")
}

func looksLikeHashedFile(p string) bool {
	// crude heuristic: /assets/app.2f7b1c9d.js
	b := path.Base(p)
	return strings.Count(b, ".") >= 2
}

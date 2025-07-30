package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func APIKeyAuth(expectedKey string, headerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(headerName)
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing API key",
			})
			return
		}
		if apiKey != expectedKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid API key",
			})
			return
		}

		c.Next()
	}
}

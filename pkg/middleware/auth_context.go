package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthContext struct {
	UserID  string
	OuGroup string // only the health check OU group
}

type contextKey string

const authCtxKey contextKey = "authContext"

func SetAuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user_id in auth context"})
			return
		}

		ouGroup, ok := c.Get("ou_group")
		if !ok {
			ouGroup = ""
		}

		auth := &AuthContext{
			UserID:  userID.(string),
			OuGroup: ouGroup.(string),
		}

		ctx := WithAuthContext(c.Request.Context(), auth)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func WithAuthContext(ctx context.Context, auth *AuthContext) context.Context {
	return context.WithValue(ctx, authCtxKey, auth)
}

func GetAuthContext(ctx context.Context) (*AuthContext, bool) {
	val, ok := ctx.Value(authCtxKey).(*AuthContext)
	return val, ok
}

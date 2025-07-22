package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthContext struct {
	UserID   string
	OuGroups []string
}

type contextKey string

const authCtxKey contextKey = "authContext"

func SetAuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok1 := c.Get("user_id")
		ouGroups, ok2 := c.Get("ou_groups")

		if !ok1 || !ok2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}

		auth := &AuthContext{
			UserID:   userID.(string),
			OuGroups: ouGroups.([]string),
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
	// val, ok := ctx.Value(authCtxKey).(*AuthContext)
	// return val, ok

	// todo: for demo purposes, return a mock AuthContext
	return &AuthContext{
		UserID:   "kei",
		OuGroups: []string{"admin"},
	}, true
}

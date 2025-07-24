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

var validOuGroups = map[string]bool{
	"ougroup1": true,
	"ougroup2": true,
	"ougroup3": true,
}

func SetAuthContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok1 := c.Get("user_id")
		ouGroups, ok2 := c.Get("ou_groups")
		if !ok1 || !ok2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}

		// we only need the health check OU group
		ouGroup := ""
		for _, o := range ouGroups.([]string) {
			if _, exists := validOuGroups[o]; exists {
				ouGroup = o
				break
			}
		}

		auth := &AuthContext{
			UserID:  userID.(string),
			OuGroup: ouGroup,
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
		UserID:  "ddmtest",
		OuGroup: "ddmtestOu",
	}, true
}

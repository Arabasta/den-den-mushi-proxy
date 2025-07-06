package middleware

import (
	"den-den-mushi-Go/internal/proxy/jwt_service"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func WsJwtMiddleware(v *jwt_service.Validator, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//h := c.GetHeader("Sec-WebSocket-Protocol")
		//rawToken, err := v.ExtractTokenFromHeader(h)
		//if err != nil {
		//	log.Error("Failed to extract JWT from WebSocket header", zap.String("header", h), zap.Error(err))
		//	c.AbortWithStatusJSON(400, gin.H{"error": "JWT validation failed"}) // don't return error details to client
		//	return
		//}
		//
		//token, claims, err := v.GetTokenAndClaims(rawToken)
		//if err != nil {
		//	log.Error("Failed to parse JWT", zap.String("rawToken", rawToken), zap.Error(err))
		//	c.AbortWithStatusJSON(401, gin.H{"error": "JWT validation failed"})
		//	return
		//}
		//
		//err = v.ValidateClaims(claims, token)
		//if err != nil {
		//	log.Error("Failed to validate claims", zap.Any("claims", claims), zap.Any("token", token), zap.Error(err))
		//	c.AbortWithStatusJSON(401, gin.H{"error": "JWT validation failed"})
		//	return
		//}
		//
		//log.Info("WebSocket JWT validated, setting claims in ctx", zap.String("user", claims.Subject), zap.String("jti", claims.ID))
		//
		//c.Header("Sec-WebSocket-Protocol", "jwt")

		log.Info("WebSocket JWT middleware called, setting mock claims in context")
		now := time.Now()
		claims := &token.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ID:        uuid.New().String(),
				ExpiresAt: jwt.NewNumericDate(now.Add(10000 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
			Connection: dto.Connection{
				Server: dto.ServerInfo{
					IP:     "52.221.194.96",
					Port:   "22",
					OSUser: "ec2-user",
				},
				Type:    "ssh_test_key",
				Purpose: "change_request",
				UserSession: dto.UserSession{
					Id:        "kei",
					StartRole: "implementor",
				},
				PtySession: dto.PtySession{
					IsNew:                     true,
					IsObserverEnabled:         true,
					MaxObservers:              10000,
					MaxHeadlessMinutes:        50000,
					MaxSessionDurationMinutes: 50000,
				},
				ChangeRequest: dto.ChangeRequest{
					Id:                       "123",
					ImplementorGroup:         "123",
					EndTime:                  time.Now().Add(100 * time.Hour).Format(time.RFC3339),
					ChangeGracePeriodMinutes: 50000,
				},
			},
		}

		log.Info("WebSocket JWT middleware mock claims", zap.Any("claims", claims))

		c.Set("claims", claims)
		c.Next()
	}
}

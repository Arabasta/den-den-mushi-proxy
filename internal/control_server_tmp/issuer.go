package control_server_tmp

import (
	"den-den-mushi-Go/pkg/connection"
	"den-den-mushi-Go/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// todo: move this to control service

type Issuer struct {
	secret []byte
	iss    string
	aud    string
	ttl    time.Duration
}

func NewIssuer(secret, issuer, audience string, ttl time.Duration) *Issuer {
	return &Issuer{
		secret: []byte(secret),
		iss:    issuer,
		aud:    audience,
		ttl:    ttl,
	}
}

func (i *Issuer) Mint(userID string, conn connection.Connection, jti string) (string, error) {
	now := time.Now()
	claims := token.Claims{
		Connection: conn,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.iss,
			Subject:   userID,
			Audience:  []string{i.aud},
			ExpiresAt: jwt.NewNumericDate(now.Add(i.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        jti,
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // todo: use asymm
	return tok.SignedString(i.secret)
}

func RegisterIssuerRoutes(r *gin.RouterGroup, issuer *Issuer, log *zap.Logger) {
	issr := r.Group("")
	issr.POST("/token", mintToken(issuer))
}

func mintToken(issuer *Issuer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			UserID     string                `json:"user_id"`
			Connection connection.Connection `json:"connection" binding:"required"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		jti := uuid.NewString()
		tok, err := issuer.Mint(body.UserID, body.Connection, jti)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"token": tok})
	}
}

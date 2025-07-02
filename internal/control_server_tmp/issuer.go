package control_server_tmp

import (
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// todo: move this to control server package

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

func (i *Issuer) Mint(userID string, conn dto.Connection, jti string) (string, error) {
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
	issr.POST("/token", mintToken(issuer, log))
}

// todo: follow https://www.rfc-editor.org/rfc/rfc8725.html#name-best-practices
func mintToken(issuer *Issuer, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.MintRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			log.Error("Failed to bind JSON", zap.Error(err))
			return
		}

		log.Info("Mint request received", zap.Any("request", body))

		// some simple example validation for now
		//todo: improve validation
		var startRole dto.StartRole

		if body.PtySessionId == "" {
			startRole = dto.Implementor // only implementor can start
		} else {
			fmt.Println(body.StartRole)
			if body.StartRole == "" {
				c.JSON(400, gin.H{"error": "start_role is required when joining existing session"})
				return
			} else if body.StartRole == dto.Implementor {
			} else if body.StartRole == dto.Observer {
				// ok, observer can join existing session
			} else {
				c.JSON(400, gin.H{"error": "invalid start_role"})
				return
			}
			startRole = body.StartRole
		}

		// Construct the Connection DTO to embed in JWT
		conn := dto.Connection{
			Server:  body.Server,
			Type:    body.Type, // todo: should be set based on server details
			Purpose: body.Purpose,
			UserSession: dto.UserSession{
				Id:        "kei", // todo: should be from ctx after auth
				StartRole: startRole,
			},
			PtySession: dto.PtySession{
				Id:    body.PtySessionId,
				IsNew: body.PtySessionId == "",

				// todo: all these should be set by config or from db
				IsObserverEnabled:         true,
				MaxObservers:              5,
				MaxHeadlessMinutes:        5,
				MaxSessionDurationMinutes: 360,
			},
			ChangeRequest: dto.ChangeRequest{
				Id:                       body.ChangeID,
				ImplementorGroup:         "myimplgroup",                                      // todo: should be from change service after validation
				EndTime:                  time.Now().Add(1 * time.Hour).Format(time.RFC3339), // todo: should be from change request
				ChangeGracePeriodMinutes: 30,                                                 // todo: should be from config or db
			},
		}

		jti := uuid.NewString()
		tokenStr, err := issuer.Mint(conn.UserSession.Id, conn, jti)
		if err != nil {
			log.Error("Failed to mint token", zap.Error(err))
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Info("Token minted successfully", zap.String("jti", jti), zap.String("userID", conn.UserSession.Id))
		c.JSON(200, gin.H{"token": tokenStr})
	}
}

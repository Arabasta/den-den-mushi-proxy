package control_server_tmp

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// todo: move this to control server package

type Issuer struct {
	cfg    *config.Config
	log    *zap.Logger
	secret []byte
	iss    string
	aud    string
	ttl    time.Duration
}

func New(cfg *config.Config, log *zap.Logger) *Issuer {
	log.Info("Initializing JWT Issuer",
		zap.String("issuer", cfg.Token.Issuer),
		zap.String("audience", cfg.Token.Audience),
		zap.Int("ttl_seconds", cfg.Token.Ttl),
	)

	return &Issuer{
		cfg:    cfg,
		log:    log,
		secret: []byte(cfg.Token.Secret),
		iss:    cfg.Token.Issuer,
		aud:    cfg.Token.Audience,
		ttl:    time.Duration(cfg.Token.Ttl) * time.Second,
	}
}

// Mint must follow https://www.rfc-editor.org/rfc/rfc8725.html
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
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // todo: use asymm, sign with priv

	// RFC8725 3.11 explicit type
	tok.Header["typ"] = "proxy/ws+jwt"

	return tok.SignedString(i.secret)
}

func RegisterIssuerRoutes(r *gin.RouterGroup, issuer *Issuer, log *zap.Logger) {
	issr := r.Group("/api/v1")
	issr.POST("/token", mintToken(issuer, log))
}

func mintToken(issuer *Issuer, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.MintRequestTmp
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			log.Error("Failed to bind JSON", zap.Error(err))
			return
		}

		log.Info("Mint request received", zap.Any("request", body))

		// todo: implement proper validation, this is just a placeholder demo simple junk thing
		var startRole types.StartRole

		if body.PtySessionId == "" {
			startRole = types.Implementor // only implementor can start
		} else {
			fmt.Println(body.StartRole)
			if body.StartRole == "" {
				c.JSON(400, gin.H{"error": "start_role is required when joining existing session"})
				return
			} else if body.StartRole == types.Implementor {
			} else if body.StartRole == types.Observer {
				// ok, observer can join existing session
			} else {
				c.JSON(400, gin.H{"error": "invalid start_role"})
				return
			}
			startRole = body.StartRole
		}

		if body.UserId == "" {
			body.UserId = "keiyam"
		}

		// Construct the Connection DTO to embed in JWT
		conn := dto.Connection{
			Server:  body.Server,
			Type:    body.Type, // todo: should be set based on server details
			Purpose: body.Purpose,
			UserSession: dto.UserSession{
				Id:        body.UserId + "/" + uuid.NewString(), // todo: should be set with keycloak user id
				StartRole: startRole,
			},
			PtySession: dto.PtySession{
				Id:    body.PtySessionId,
				IsNew: body.PtySessionId == "",
			},
			ChangeRequest: dto.ChangeRequest{
				Id:                body.ChangeID,
				ImplementorGroups: make([]string, 0),             // todo: should be from change service after validation
				EndTime:           time.Now().Add(1 * time.Hour), // todo: should be from change request
			},
			FilterType: body.FilterType,
		}

		jti := uuid.NewString() + strconv.FormatInt(time.Now().Unix(), 10)

		tokenStr, err := issuer.Mint(body.UserId, conn, jti)
		if err != nil {
			log.Error("Failed to mint token", zap.Error(err))
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Info("Token minted successfully", zap.String("jti", jti), zap.String("userID", conn.UserSession.Id))
		c.JSON(200, gin.H{"token": tokenStr})
	}
}

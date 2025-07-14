package jwt

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/pkg/dto"
	"den-den-mushi-Go/pkg/token"
	"den-den-mushi-Go/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Issuer struct {
	cfg    *config.Config
	log    *zap.Logger
	secret []byte
	iss    string
	ttl    time.Duration
}

func New(cfg *config.Config, log *zap.Logger) *Issuer {
	log.Info("Initializing JWT Issuer",
		zap.String("issuer", cfg.Token.Issuer),
		zap.Int("ttl_seconds", cfg.Token.Ttl),
	)

	return &Issuer{
		cfg:    cfg,
		log:    log,
		secret: []byte(cfg.Token.Secret),
		iss:    cfg.Token.Issuer,
		ttl:    time.Duration(cfg.Token.Ttl) * time.Second,
	}
}

// Mint must follow https://www.rfc-editor.org/rfc/rfc8725.html
func (i *Issuer) Mint(userID string, conn *dto.Connection, proxyType types.Proxy) (string, error) {
	now := time.Now()
	claims := token.Claims{
		Connection: *conn,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    i.iss,
			Subject:   userID,
			Audience:  proxyType.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString() + strconv.FormatInt(time.Now().Unix(), 10),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // todo: use asymm, sign with priv

	// RFC8725 3.11 explicit type
	tok.Header["typ"] = "proxy/ws+jwt"

	return tok.SignedString(i.secret)
}

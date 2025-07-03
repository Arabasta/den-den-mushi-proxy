package jwt_service

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Validator struct {
	parser *jwt.Parser
	secret []byte
	replay *jtiStore
	cfg    *config.Config
	log    *zap.Logger
}

func NewValidator(p *jwt.Parser, secret string, ttl time.Duration, cfg *config.Config, log *zap.Logger) *Validator {
	v := &Validator{
		parser: p,
		secret: []byte(secret),
		replay: &jtiStore{ttl: ttl},
		cfg:    cfg,
		log:    log}
	return v
}

// ValidateClaims Validate must follow https://www.rfc-editor.org/rfc/rfc8725.html
func (v *Validator) ValidateClaims(claims *token.Claims, tok *jwt.Token) error {
	v.log.Debug("Validating claims", zap.Any("claims", claims), zap.String("jti", claims.ID))

	// todo: RFC8725 3.8 validate subject against keycloak user

	// check typ RFC8725 3.11 3.12
	if !v.isExpectedTyp(tok.Header["typ"].(string), v.cfg.Token.ExpectedTyp) {
		v.log.Error("Token has unexpected type", zap.String("jti", claims.ID), zap.String("typ", tok.Header["typ"].(string)))
		return errors.New("token has unexpected type")
	}

	// check expiration
	if !v.isBeforeExp(claims.RegisteredClaims.ExpiresAt) {
		v.log.Error("Token is expired", zap.String("jti", claims.ID))
		return errors.New("token is expired")
	}

	// check replay
	if v.replay.isConsumed(claims.ID) {
		v.log.Error("Token already consumed", zap.String("jti", claims.ID))
		return errors.New("token already consumed")
	}

	if !v.replay.consume(claims.ID) {
		v.log.Error("Failed to consume token", zap.String("jti", claims.ID))
		return errors.New("failed to consume token")
	}

	return nil
}

func (v *Validator) isExpectedTyp(typ string, expectedTyp string) bool {
	if strings.TrimSpace(typ) == "" || strings.TrimSpace(expectedTyp) == "" {
		return false
	}
	return strings.TrimSpace(typ) == expectedTyp
}

func (v *Validator) isBeforeExp(exp *jwt.NumericDate) bool {
	if exp == nil {
		return false
	}
	return time.Now().Before(exp.Time)
}

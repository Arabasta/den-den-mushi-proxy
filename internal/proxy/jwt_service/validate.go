package jwt_service

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/jwt_service/jti"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strings"
)

type Validator struct {
	parser *jwt.Parser
	secret []byte
	jti    *jti.Service
	cfg    *config.Config
	log    *zap.Logger
}

func NewValidator(p *jwt.Parser, jti *jti.Service, secret string, cfg *config.Config, log *zap.Logger) *Validator {
	v := &Validator{
		parser: p,
		secret: []byte(secret),
		jti:    jti,
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

	// check replay
	if v.jti.IsConsumed(claims.ID) {
		v.log.Error("Token already consumed", zap.String("jti", claims.ID))
		return errors.New("token already consumed")
	}

	if !v.jti.Consume(claims.ID) {
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

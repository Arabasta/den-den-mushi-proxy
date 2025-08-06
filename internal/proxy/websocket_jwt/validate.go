package websocket_jwt

import (
	"den-den-mushi-Go/internal/proxy/websocket_jwt/jti"
	"den-den-mushi-Go/pkg/config"
	"den-den-mushi-Go/pkg/middleware"
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
	cfg    *config.JwtAudience
	log    *zap.Logger
}

func NewValidator(p *jwt.Parser, jti *jti.Service, cfg *config.JwtAudience, log *zap.Logger) *Validator {
	log.Info("Initializing JWT Validator...")
	return &Validator{
		parser: p,
		secret: []byte(cfg.Secret),
		jti:    jti,
		cfg:    cfg,
		log:    log,
	}
}

// ValidateClaims Validate must follow https://www.rfc-editor.org/rfc/rfc8725.html
func (v *Validator) ValidateClaims(claims *token.Claims, tok *jwt.Token, authCtx *middleware.AuthContext, cfg *config.Tmpauth) error {
	v.log.Debug("Validating claims", zap.Any("claims", claims), zap.String("jti", claims.ID))

	// RFC8725 3.8 validate subject against keycloak user
	if claims.Subject == "" {
		v.log.Error("Token has no subject", zap.String("jti", claims.ID))
	}
	if cfg.IsEnabled && claims.Subject != authCtx.UserID {
		v.log.Error("Token subject does not match auth context user ID", zap.String("jti", claims.ID),
			zap.String("subject", claims.Subject),
			zap.String("authUserID", authCtx.UserID))
		return errors.New("token subject does not match auth context user ID")
	}
	v.log.Debug("Validated subject", zap.String("jti", claims.ID), zap.String("subject", claims.Subject))

	// check typ RFC8725 3.11 3.12
	if !v.isExpectedTyp(tok.Header["typ"].(string), v.cfg.ExpectedTyp) {
		v.log.Error("Token has unexpected type", zap.String("jti", claims.ID), zap.String("typ", tok.Header["typ"].(string)))
		return errors.New("token has unexpected type")
	}

	if !v.jti.ConsumeIfNotExists(claims) {
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

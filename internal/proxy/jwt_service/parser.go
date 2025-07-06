package jwt_service

import (
	"den-den-mushi-Go/internal/proxy/config"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// NewParser must follow https://www.rfc-editor.org/rfc/rfc8725.html
func NewParser(cfg *config.Config, log *zap.Logger) *jwt.Parser {
	log.Info("Initialising JWT parser. Audience: " + cfg.Token.Audience + ", Issuer: " + cfg.Token.Issuer)

	return jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),

		// audience depends on the type of proxy
		// e.g., proxy-osdba, proxy-storage, proxy-appliances, etc
		jwt.WithAudience(cfg.Token.Audience),

		// RFC8725 3.8 validate issuer
		jwt.WithIssuer(cfg.Token.Issuer),

		// makes exp claim mandatory and checks it
		jwt.WithExpirationRequired(),

		// makes iat claim mandatory and checks it
		jwt.WithIssuedAt())
}

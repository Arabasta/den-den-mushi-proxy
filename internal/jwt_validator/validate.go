package jwt_validator

import (
	"context"
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Validator struct {
	parser *jwt.Parser
	secret []byte
	replay *jtiStore
}

func NewValidator(p *jwt.Parser, secret string, ttl time.Duration) *Validator {
	v := &Validator{
		parser: p,
		secret: []byte(secret),
		replay: &jtiStore{ttl: ttl}}
	return v
}

// todo: follow https://www.rfc-editor.org/rfc/rfc8725.html#name-best-practices
func (v *Validator) Validate(_ context.Context, tokenString string) (*token.Claims, error) {
	claims := new(token.Claims)
	_, err := v.parser.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return v.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if ok := v.replay.Check(claims.ID); !ok {
		return nil, errors.New("replay detected")
	}
	return claims, nil
}

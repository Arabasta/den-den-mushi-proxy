package jwt_service

import (
	"den-den-mushi-Go/pkg/token"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strings"
)

// ExtractTokenFromHeader just gets the second value
func (v *Validator) ExtractTokenFromHeader(h string) (string, error) {
	v.log.Debug("Extracting token from header", zap.String("header", h))
	parts := strings.SplitN(h, ",", 2)
	if len(parts) != 2 {
		v.log.Error("Failed to extract JWT from header", zap.String("header", h))
		return "", errors.New("failed to extract JWT from header")
	}

	return strings.TrimSpace(parts[1]), nil
}

// GetTokenAndClaims parses jwt, verifies signature. Returns token and custom claims
func (v *Validator) GetTokenAndClaims(tokenString string) (*jwt.Token, *token.Claims, error) {
	v.log.Debug("Extracting token and claims from raw token", zap.String("token", tokenString))
	claims := new(token.Claims)

	// ParseWithClaims handles RFC8725 3.1 algorithm verification, algo passed in NewParser
	tok, err := v.parser.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (interface{}, error) {
		return v.secret, nil
	})
	if err != nil {
		return nil, nil, err
	}

	return tok, claims, nil
}

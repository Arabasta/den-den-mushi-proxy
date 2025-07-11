package jwt_service

import (
	"den-den-mushi-Go/pkg/token"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"strings"
)

func (v *Validator) ExtractProxyTokenFromHeader(h string) (string, error) {
	v.log.Debug("Extracting token from header", zap.String("header", h))
	h = strings.TrimSpace(h)

	const expectedPrefix = "X-Proxy-Session-Token,"
	if !strings.HasPrefix(h, expectedPrefix) {
		v.log.Error("Invalid header format - missing expected prefix",
			zap.String("header", h),
			zap.String("expected", expectedPrefix))
		return "", fmt.Errorf("invalid header format: expected '%s' prefix", expectedPrefix)
	}

	// get second part of header
	tokenString := strings.TrimSpace(h[len(expectedPrefix):])
	if tokenString == "" {
		v.log.Error("Empty token after header prefix",
			zap.String("header", h))
		return "", errors.New("empty session token")
	}

	// simple JWT format validation
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		v.log.Error("Invalid JWT format - wrong number of segments",
			zap.String("token", tokenString))
		return "", errors.New("invalid JWT format")
	}

	return tokenString, nil
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

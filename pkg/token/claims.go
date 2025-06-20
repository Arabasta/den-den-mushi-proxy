package token

import (
	"den-den-mushi-Go/pkg/connection"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Connection connection.Connection `json:"connection"`
	jwt.RegisteredClaims
}

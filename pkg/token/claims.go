package token

import (
	"den-den-mushi-Go/pkg/dto/connection"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Connection connection.Connection `json:"connection"`
	jwt.RegisteredClaims
}

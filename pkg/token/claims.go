package token

import (
	"den-den-mushi-Go/pkg/dto"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Connection dto.Connection `json:"connection"`
	jwt.RegisteredClaims
	OuGroup string `json:"ou_group,omitempty"`
}

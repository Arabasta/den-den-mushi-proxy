package types

import "github.com/golang-jwt/jwt/v5"

type Proxy string

func (p Proxy) String() jwt.ClaimStrings {
	return jwt.ClaimStrings{string(p)}
}

const (
	OS      Proxy = "OS"
	DBA     Proxy = "DBA"
	STORAGE Proxy = "STORAGE"
	NETWORK Proxy = "NETWORK"
)

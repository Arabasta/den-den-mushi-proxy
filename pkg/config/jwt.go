package config

import "github.com/spf13/viper"

type JwtIssuer struct {
	Issuer        string
	Secret        string
	ExpirySeconds int
}

type JwtAudience struct {
	ExpectedIssuer   string
	ExpectedAudience string
	ExpectedTyp      string
	Secret           string
}

func BindJwtIssuerSecret(v *viper.Viper) {
	_ = v.BindEnv("jwtissuer.secret", "JWT_ISSUER_SECRET")
}

func BindJwtAudienceSecret(v *viper.Viper) {
	_ = v.BindEnv("jwtaudience.secret", "JWT_AUDIENCE_SECRET")
}

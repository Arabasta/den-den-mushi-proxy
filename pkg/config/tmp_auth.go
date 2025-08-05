package config

import "github.com/spf13/viper"

type Tmpauth struct {
	IsEnabled  bool
	UserIdKey  string `json:"UserIdKey"`
	OuGroupKey string `json:"OuGroupKey"`
	Secret     string `json:"-"`
}

func BindTmpAuthSecret(v *viper.Viper) {
	_ = v.BindEnv("tmpauth.secret", "TMPAUTH_SECRET")
}

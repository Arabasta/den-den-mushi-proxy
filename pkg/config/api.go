package config

import "github.com/spf13/viper"

type Api struct {
	IsKeyAuthEnabled bool
	Key              string `json:"-"` // API key for authentication
	KeyHeader        string
}

func BindApi(v *viper.Viper) {
	_ = v.BindEnv("api.key", "API_KEY")
}

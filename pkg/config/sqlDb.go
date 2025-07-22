package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type SqlDb struct {
	User     string
	Password string `json:"-"`
	Host     string `json:"-"`
	Port     int
	DBName   string
	Params   string
	// Pooling
	MaxIdleConns           int
	MaxOpenConns           int
	ConnMaxLifetimeMinutes int
}

func BindSqlDb(v *viper.Viper, envPrefix, viperPrefix string) *SqlDb {
	keys := []string{"password", "host"}

	for _, key := range keys {
		envKey := fmt.Sprintf("%s_%s", envPrefix, strings.ToUpper(key))
		viperKey := fmt.Sprintf("%s.%s", viperPrefix, key)
		_ = v.BindEnv(viperKey, envKey)
	}

	var db SqlDb
	_ = v.UnmarshalKey(viperPrefix, &db)
	return &db
}

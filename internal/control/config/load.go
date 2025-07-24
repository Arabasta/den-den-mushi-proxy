package config

import (
	"den-den-mushi-Go/pkg/config"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

func Load(path string) *Config {
	var cfg Config

	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType("json")
	v.AddConfigPath(dir)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// env bindings
	config.BindSsl(v)
	config.BindRedis(v)
	config.BindJwtIssuerSecret(v)
	cfg.DdmDB = config.BindSqlDb(v, "DDM_DB", "DdmDB")
	cfg.InvDB = config.BindSqlDb(v, "INV_DB", "InvDB")

	fmt.Println("Config host", cfg.DdmDB.Host)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return &cfg
}

func HotReload(path string, cfg *Config) {
	v := viper.New()
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
}

package config

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func Load(path string) *Config {
	dir := filepath.Dir(path)
	filename := filepath.Base(path)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	v := viper.New()
	v.SetConfigName(name)
	v.SetConfigType("json")
	v.AddConfigPath(dir)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg Config

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

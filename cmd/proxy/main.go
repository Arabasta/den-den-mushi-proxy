package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/server"
	"den-den-mushi-Go/pkg/logger"
	"flag"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	cfg := config.Load(configPath())

	log := logger.Init(cfg)
	defer func() {
		_ = log.Sync()
	}()

	s := server.New(static.Files, cfg, log)
	if err := server.Start(s, cfg); err != nil {
		log.Fatal("failed to start server: %v", zap.Error(err))
	}
}

func configPath() string {
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	var finalPath string
	if *configPath != "" {
		finalPath = *configPath
	} else {
		// default to config.json next to executable
		exe, _ := os.Executable()
		exeDir := filepath.Dir(exe)
		finalPath = filepath.Join(exeDir, "config.json")
	}

	return finalPath
}

package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/logger"
	"den-den-mushi-Go/internal/control/server"
	"flag"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	cfg := config.Load(configPath())

	log := logger.Init(cfg)
	if log == nil {
		panic("failed to initialize logger")
	}
	defer func() {
		_ = log.Sync()
	}()

	s := server.New(static.Files, cfg, log)
	if err := server.Start(s, cfg); err != nil {
		log.Fatal("failed to start server: %v", zap.Error(err))
	}
}

// configPath usage: go run main.go -config /path/to/config.json
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

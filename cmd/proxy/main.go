package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/server"
	"den-den-mushi-Go/pkg/logger"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)
	cfg := config.Load(filepath.Join(exeDir, "config.json"))

	log := logger.Init(cfg)
	defer func() {
		_ = log.Sync()
	}()

	s := server.New(static.Files, cfg, log)
	if err := server.Start(s, cfg); err != nil {
		log.Fatal("failed to start server: %v", zap.Error(err))
	}
}

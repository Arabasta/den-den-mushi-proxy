package main

import (
	"den-den-mushi-Go/internal/config"
	"den-den-mushi-Go/internal/server"
	"den-den-mushi-Go/pkg/logger"
	"embed"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

//go:embed static/* static/css/* static/js/* static/js/settings/*
var staticFiles embed.FS

func main() {
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)
	cfg := config.Load(filepath.Join(exeDir, "config.json"))

	log := logger.Init(cfg)
	defer func() {
		_ = log.Sync()
	}()

	s := server.New(staticFiles, cfg, log)
	if err := server.Start(s, cfg); err != nil {
		log.Fatal("failed to start server: %v", zap.Error(err))
	}
}

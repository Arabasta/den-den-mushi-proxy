package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/llm_external/config"
	"den-den-mushi-Go/internal/llm_external/server"
	"den-den-mushi-Go/pkg/logger"
	"den-den-mushi-Go/pkg/mysql"
	"flag"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load(".env")
	cfg := config.Load(configPath())

	log := logger.Init(cfg.Logger, cfg.App)
	if log == nil {
		panic("failed to initialize logger")
	}
	defer func() {
		_ = log.Sync()
	}()
	var db *gorm.DB
	var err error

	if !cfg.Development.IsUsingInvDb {
		db, err = mysql.Client(cfg.DdmDB, cfg.Ssl, log)
	} else {
		db, err = mysql.Client(cfg.InvDB, cfg.Ssl, log)
	}
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	s := server.New(db, root.Files, cfg, log)
	if err := server.Start(s, cfg, log); err != nil {
		log.Fatal("failed to start server", zap.Error(err))
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

package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/server"
	configpkg "den-den-mushi-Go/pkg/config"
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/logger"
	"den-den-mushi-Go/pkg/mysql"
	"flag"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	_ = godotenv.Load(".env")
	cfg := config.Load(configPath())

	log := logger.Init(configpkg.Logger{
		Level:       cfg.Logger.Level,
		Format:      cfg.Logger.Format,
		Output:      cfg.Logger.Output,
		FilePath:    cfg.Logger.FilePath,
		Environment: cfg.App.Environment,
	})
	if log == nil {
		panic("failed to initialize logger")
	}
	defer func() {
		_ = log.Sync()
	}()

	ddmDb, err := mysql.Client(configpkg.SqlDb{
		User:                   cfg.DdmDB.User,
		Password:               cfg.DdmDB.Password,
		Host:                   cfg.DdmDB.Host,
		Port:                   cfg.DdmDB.Port,
		DBName:                 cfg.DdmDB.DBName,
		Params:                 cfg.DdmDB.Params,
		MaxIdleConns:           cfg.DdmDB.MaxIdleConns,
		MaxOpenConns:           cfg.DdmDB.MaxOpenConns,
		ConnMaxLifetimeMinutes: cfg.DdmDB.ConnMaxLifetimeMinutes}, log)
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if cfg.App.Environment != "prod" {
		log.Info("Running AutoMigrate for non-production environment")
		if err := ddmDb.AutoMigrate(&pty_sessions.Model{}, &connections.Model{}, &proxy_host.Model{}); err != nil {
			log.Fatal("Failed to auto-migrate pty_sessions table", zap.Error(err))
		}
	}
	s := server.New(root.Files, ddmDb, cfg, log)
	if err := server.Start(s, cfg, log); err != nil {
		log.Fatal("failed to start server: %v", zap.Error(err))
	}
	// todo: add graceful shutdown
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

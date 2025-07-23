package main

import (
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/jwt_service/jti"
	"den-den-mushi-Go/internal/proxy/server"
	"den-den-mushi-Go/pkg/dto/connections"
	"den-den-mushi-Go/pkg/dto/proxy_host"
	"den-den-mushi-Go/pkg/dto/pty_sessions"
	"den-den-mushi-Go/pkg/logger"
	"den-den-mushi-Go/pkg/mysql"
	redispkg "den-den-mushi-Go/pkg/redis"
	"flag"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"path/filepath"
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

	if cfg.Development.IsSMX && cfg.App.Environment != "prod" && cfg.Development.IsAutoMigrateEnabled {
		log.Info("Running AutoMigrate for non-production environment")
		if err := db.AutoMigrate(&pty_sessions.Model{}, &connections.Model{}, &proxy_host.Model{},
			&jti.Model{}); err != nil {
			log.Fatal("Failed to auto-migrate", zap.Error(err))
		}
	}

	var redisClient *redis.Client

	if cfg.Development.UseRedis {
		redisClient, err = redispkg.Client(cfg.Redis, log) // todo pass redis to server
		if err != nil {
			log.Fatal("Failed to connect to Redis cluster", zap.Error(err))
		}
	}

	s := server.New(root.Files, db, redisClient, cfg, log)
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

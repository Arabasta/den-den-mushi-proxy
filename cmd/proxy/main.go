package main

import (
	"context"
	"den-den-mushi-Go"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/server"
	"den-den-mushi-Go/internal/proxy/websocket_jwt/jti"
	configpkg "den-den-mushi-Go/pkg/config"
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
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
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

	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		log.Error("Failed to get hostname during startup", zap.Error(err))
		os.Exit(1)
	}
	cfg.Host = &configpkg.Host{Name: hostname}
	log.Info("Starting Den Den Mushi Proxy", zap.String("version", cfg.App.Version), zap.String("hostname", cfg.Host.Name))

	var db *gorm.DB

	if !cfg.Development.IsUsingInvDb {
		db, err = mysql.Client(cfg.DdmDB, cfg.Ssl, log)
	} else {
		db, err = mysql.Client(cfg.InvDB, cfg.Ssl, log)
	}
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	if cfg.Development.IsSMX && cfg.Development.IsAutoMigrateEnabled {
		if err := db.AutoMigrate(&pty_sessions.Model{}, &connections.Model{}, &proxy_host.Model{},
			&jti.Model{}); err != nil {
			log.Fatal("Failed to auto-migrate", zap.Error(err))
		}
	}

	var redisClient *redis.Client

	if cfg.Development.UseRedis {
		redisClient, err = redispkg.Client(cfg.Redis, log)
		if err != nil {
			log.Fatal("Failed to connect to Redis cluster", zap.Error(err))
		}
	}

	s, sessionManager := server.New(root.Files, db, redisClient, cfg, log)

	err = sessionManager.CleanupActiveSessionsAndConnections()
	if err != nil {
		log.Error("Failed to cleanup active sessions and connections", zap.Error(err))
		// continue with server startup even if cleanup fails
	}

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-stop
		log.Info("Shutting down gracefully...")

		// todo: close http server listeners

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		sessionManager.Shutdown(ctx)

		log.Info("Shutdown requested, forcing exit in 10s...")
		time.AfterFunc(10*time.Second, func() {
			log.Warn("Forcing process exit")
			os.Exit(0)
		})

		// block until shutdown tasks finish or timeout
		<-ctx.Done()

		log.Info("Shutdown complete")
		os.Exit(0)
	}()

	if err := server.Start(s, cfg, log); err != nil {
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

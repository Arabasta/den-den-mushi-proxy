package server

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/pkg/middleware"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	log    *zap.Logger
}

func setupSecurityHeaders(r *gin.Engine) {

}

func New(staticFiles embed.FS, db *gorm.DB, redis *redis.Client, cfg *config.Config, log *zap.Logger) *Server {
	deps := initDependencies(db, redis, cfg, log)

	r := gin.New()
	r.Use(
		middleware.RequestLogger(log),
		middleware.Cors(cfg.Cors, log),
		gin.Recovery())

	registerUnprotectedRoutes(r, deps, cfg, log)
	registerWebsocketRoutes(r, deps, cfg, log)
	addStaticRoutes(r, staticFiles, cfg, log)

	return &Server{engine: r, cfg: cfg, log: log}
}

func Start(s *Server, cfg *config.Config, log *zap.Logger) error {
	addr := fmt.Sprintf(":%d", cfg.App.Port)

	if !cfg.Ssl.Enabled {
		log.Info("Starting server without TLS", zap.String("address", addr))
		return s.engine.Run(addr)
	} else {
		log.Info("Starting server with TLS")
		log.Debug("Cert details", zap.String("cert", cfg.Ssl.CertFile),
			zap.String("key", cfg.Ssl.KeyFile))
		return s.engine.RunTLS(addr, cfg.Ssl.CertFile, cfg.Ssl.KeyFile)
	}
}

package server

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/middleware"
	middleware2 "den-den-mushi-Go/pkg/middleware"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	log    *zap.Logger
}

func setupSecurityHeaders(r *gin.Engine) {

}

func New(staticFiles embed.FS, cfg *config.Config, log *zap.Logger) *Server {
	deps := initDependencies(cfg, log)

	r := gin.New()
	r.Use(
		middleware2.RequestLogger(log),
		middleware.Cors(cfg, log),
		gin.Recovery())

	registerUnprotectedRoutes(r, deps, cfg, log)
	registerWebsocketRoutes(r, deps, cfg, log)
	addStaticRoutes(r, staticFiles, cfg, log)

	return &Server{engine: r, cfg: cfg, log: log}
}

func Start(s *Server, cfg *config.Config) error {
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	return s.engine.Run(addr)
}

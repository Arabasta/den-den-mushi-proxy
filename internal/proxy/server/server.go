package server

import (
	"crypto/tls"
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/core/session_manager"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/profiler"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	log    *zap.Logger
}

func New(staticFiles embed.FS, db *gorm.DB, redis *redis.Client, cfg *config.Config, log *zap.Logger) (*Server, *session_manager.Service) {
	deps := initDependencies(db, redis, cfg, log)

	if cfg.App.Environment == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(
		gin.Recovery(),
		middleware.RequestLogger(log),
		middleware.Cors(cfg.Cors, log),
		middleware.CsrfGuardNoCookies(log),
		middleware.Security(cfg.Security, cfg.Ssl.Enabled),
		middleware.StripSetCookie(),
		middleware.MaxBody(10<<20), // 10 mb
	)

	if cfg.Development.IsDevRoutesEnabled {
		registerUnprotectedRoutes(r, deps, cfg, log)
		addStaticRoutes(r, staticFiles, cfg, log)
	}

	profiler.StartSidecar(cfg.Pprof, log)

	registerWebsocketRoutes(r, deps, cfg, log)

	r.NoRoute(func(c *gin.Context) { c.JSON(404, gin.H{"error": "not found"}) })
	r.NoMethod(func(c *gin.Context) { c.JSON(405, gin.H{"error": "method not allowed"}) })

	return &Server{engine: r, cfg: cfg, log: log}, deps.SessionManager
}

func Start(s *Server, cfg *config.Config, log *zap.Logger) error {
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	if cfg.App.IsLocalHost {
		addr = fmt.Sprintf("127.0.0.1:%d", cfg.App.Port)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MiB
	}

	if !cfg.Ssl.Enabled {
		log.Info("Starting server without TLS", zap.String("address", addr))
		return srv.ListenAndServe()
	}

	srv.TLSConfig = tls13OnlyConfig()

	log.Info("Starting server with TLS 1.3 only",
		zap.String("address", addr),
		zap.String("cert", cfg.Ssl.CertFile),
		zap.String("key", cfg.Ssl.KeyFile),
	)

	return srv.ListenAndServeTLS(cfg.Ssl.CertFile, cfg.Ssl.KeyFile)
}

func tls13OnlyConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
		NextProtos: []string{"h2", "http/1.1"},
	}
}

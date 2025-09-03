package server

import (
	"crypto/tls"
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/profiler"
	"embed"
	"fmt"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
	cfg    *config.Config
	log    *zap.Logger
}

func New(ddmDb *gorm.DB, staticFiles embed.FS, cfg *config.Config, log *zap.Logger) *Server {
	deps := initDependencies(ddmDb, cfg, log)

	if cfg.App.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(ginzap.Ginzap(log, "", true))
	r.Use(
		gin.Recovery(),
		middleware.RequestLogger(log),
		middleware.Cors(cfg.Cors, log),
		middleware.CsrfGuardNoCookies(log),
		middleware.Security(cfg.Security, cfg.Ssl.Enabled),
		middleware.MaxBody(10<<20), // 10 mb
	)

	profiler.Start(cfg.Pprof, log)

	serveSwagger(r, cfg.Swagger, log)
	registerProtectedRoutes(r, deps, cfg, log)
	addStaticRoutes(r, staticFiles, cfg, log)

	r.NoRoute(func(c *gin.Context) { c.JSON(404, gin.H{"error": "not found"}) })
	r.NoMethod(func(c *gin.Context) { c.JSON(405, gin.H{"error": "method not allowed"}) })

	return &Server{engine: r, cfg: cfg, log: log}
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
		MaxHeaderBytes:    1 << 20, // 1mb
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

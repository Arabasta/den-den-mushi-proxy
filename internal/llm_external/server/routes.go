package server

import (
	"den-den-mushi-Go/internal/llm_external/config"
	"den-den-mushi-Go/internal/llm_external/pty_sessions"
	oapi "den-den-mushi-Go/openapi/llm_external"
	"den-den-mushi-Go/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func registerProtectedRoutes(r *gin.Engine, deps *Deps, cfg *config.Config, log *zap.Logger) {
	protected := r.Group("")

	if cfg.Api.IsKeyAuthEnabled {
		log.Info("API Key authentication is enabled", zap.String("header", cfg.Api.KeyHeader))
		protected.Use(
			middleware.APIKeyAuth(cfg.Api.Key, cfg.Api.KeyHeader),
		)
	}

	m := &MasterHandler{
		PtySessions: &pty_sessions.Handler{
			Service: deps.PtySessionService,
			Log:     log,
		},
	}

	oapi.RegisterHandlers(protected, m)
}

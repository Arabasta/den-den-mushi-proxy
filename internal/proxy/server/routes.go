package server

import (
	"den-den-mushi-Go/internal/proxy/config"
	"den-den-mushi-Go/internal/proxy/tmp/admin_server_tmp"
	"den-den-mushi-Go/internal/proxy/tmp/control_server_tmp"
	"den-den-mushi-Go/internal/proxy/websocket"
	"den-den-mushi-Go/internal/proxy/websocket_jwt"
	middlewarepkg "den-den-mushi-Go/pkg/middleware"
	"embed"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

func registerUnprotectedRoutes(r *gin.Engine, deps *Deps, cfg *config.Config, log *zap.Logger) {
	unprotected := r.Group("")

	if cfg.App.Environment == "dev" {
		control_server_tmp.RegisterIssuerRoutes(unprotected, deps.Issuer, log)
		admin_server_tmp.RegisterAdminRoutes(unprotected, deps.SessionManager, log)
	}
}

func registerWebsocketRoutes(r *gin.Engine, deps *Deps, cfg *config.Config, log *zap.Logger) {
	protected := r.Group("")
	protected.Use(
		middlewarepkg.TmpAuth(log, cfg.TmpAuth),
		middlewarepkg.SetAuthContext(),
		websocket_jwt.Middleware(deps.Validator, log, cfg),
	)

	websocket.RegisterWebsocketRoutes(protected, log, deps.WebsocketService)
}

func addStaticRoutes(r *gin.Engine, staticFiles embed.FS, cfg *config.Config, log *zap.Logger) {
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal("Failed to load embedded static files", zap.Error(err))
	}

	r.StaticFS("/static", http.FS(subFS))

	r.GET("/", func(c *gin.Context) {
		data, err := fs.ReadFile(subFS, "index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to load index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
}

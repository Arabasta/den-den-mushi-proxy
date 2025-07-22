package server

import (
	"den-den-mushi-Go/internal/control/config"
	"den-den-mushi-Go/internal/control/pty_token"
	oapi "den-den-mushi-Go/openapi/control"
	"embed"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/fs"
	"net/http"
)

func registerProtectedRoutes(r *gin.Engine, deps *Deps, cfg *config.Config, log *zap.Logger) {
	protected := r.Group("")
	protected.Use(
	// todo use keycloak  / webseal middleware
	)

	m := &MasterHandler{
		PtyHandler: &pty_token.Handler{
			Service: deps.PtyTokenService,
			Log:     log,
		},
		//MakeChangeHandler: &make_change.Handler{
		//	Service: deps.MakeChangeService,
		//	Log:     log,
		//}
	}

	oapi.RegisterHandlers(protected, m)
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

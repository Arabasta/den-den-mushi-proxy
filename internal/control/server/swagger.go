package server

import (
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/config"
	"github.com/gin-gonic/gin"
	swgui "github.com/swaggest/swgui/v4"

	"go.uber.org/zap"
)

func serveSwagger(r *gin.Engine, cfg *config.App, log *zap.Logger) {
	if cfg.Environment == "dev" {
		log.Info("Serving Swagger UI in dev environment")
		r.StaticFile("/swagger-spec/control.yaml", "./swagger/control.yaml")

		r.GET("/swagger-spec/control.json", func(c *gin.Context) {
			swagger, err := oapi.GetSwagger()
			if err != nil {
				c.JSON(500, gin.H{"error": "cannot load control swagger"})
				return
			}
			c.JSON(200, swagger)
		})

		r.GET("/swagger/control/*any", gin.WrapH(swgui.NewHandler(
			"Den Den Mushi Control API Docs",
			"/swagger-spec/control.json",
			"/swagger/control/",
		)))

	} else {
		log.Info("Not serving Swagger UI in non-dev environment")
	}
}

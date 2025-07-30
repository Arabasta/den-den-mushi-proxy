package server

import (
	oapi "den-den-mushi-Go/openapi/llm_external"
	"den-den-mushi-Go/pkg/config"
	"github.com/gin-gonic/gin"
	swgui "github.com/swaggest/swgui/v4"

	"go.uber.org/zap"
)

func serveSwagger(r *gin.Engine, cfg *config.App, log *zap.Logger) {
	if cfg.Environment == "dev" {
		log.Info("Serving Swagger UI in dev environment")
		r.StaticFile("/swagger-spec/llm_external.yaml", "./swagger/llm_external.yaml")

		r.GET("/swagger-spec/llm_external.json", func(c *gin.Context) {
			swagger, err := oapi.GetSwagger()
			if err != nil {
				c.JSON(500, gin.H{"error": "cannot load llm_external swagger"})
				return
			}
			c.JSON(200, swagger)
		})

		r.GET("/swagger/llm_external/*any", gin.WrapH(swgui.NewHandler(
			"Den Den Mushi LLM External API Docs",
			"/swagger-spec/llm_external.json",
			"/swagger/llm_external/",
		)))

	} else {
		log.Info("Not serving Swagger UI in non-dev environment")
	}
}

package server

import (
	"den-den-mushi-Go/internal/admin/config"
	// oapi "den-den-mushi-Go/openapi/admin"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func registerProtectedRoutes(r *gin.Engine, deps *Deps, cfg *config.Config, log *zap.Logger) {
	// todo : add handlers
	//protected := r.Group("")
	//
	//m := &MasterHandler{}
	//
	//oapi.RegisterHandlers(protected, m)
}

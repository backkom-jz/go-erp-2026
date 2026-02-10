package bootstrap

import (
	"go-erp/internal/api"
	"go-erp/internal/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(cfg ServerConfig, logger *zap.Logger) *gin.Engine {
	if cfg.Mode != "" {
		gin.SetMode(cfg.Mode)
	}

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	api.RegisterRoutes(r)
	return r
}

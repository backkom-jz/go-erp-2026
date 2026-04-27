package bootstrap

import (
	"go-erp/internal/api"
	"go-erp/internal/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRouter(cfg ServerConfig, logger *zap.Logger, app *App) *gin.Engine {
	if cfg.Mode != "" {
		gin.SetMode(cfg.Mode)
	}

	r := gin.New()
	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())
	r.Use(middleware.JWT(app.JWTManager))
	r.Use(middleware.Tenant())
	r.Use(middleware.RBAC())
	r.Use(middleware.Idempotency(app.IdempotencyStore, 5*time.Minute))

	api.RegisterRoutes(r, api.RouteModules{
		Auth:      app.AuthHandler,
		User:      app.UserHandler,
		Product:   app.ProductHandler,
		Inventory: app.InventoryHandler,
		Order:     app.OrderHandler,
		Payment:   app.PaymentHandler,
	})
	return r
}

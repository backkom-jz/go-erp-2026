package api

import (
	"github.com/gin-gonic/gin"
)

type Module interface {
	// Register 注册模块路由。
	Register(rg *gin.RouterGroup)
}

type RouteModules struct {
	AI        Module
	Auth      Module
	User      Module
	Product   Module
	Inventory Module
	Order     Module
	Payment   Module
}

// RegisterRoutes 注册所有 v1 业务路由。
// 备注：统一前缀为 /api/v1。
func RegisterRoutes(r *gin.Engine, modules RouteModules) {
	v1 := r.Group("/api/v1")
	{
		modules.AI.Register(v1)
		modules.Auth.Register(v1)
		modules.User.Register(v1)
		modules.Product.Register(v1)
		modules.Inventory.Register(v1)
		modules.Order.Register(v1)
		modules.Payment.Register(v1)
	}
}

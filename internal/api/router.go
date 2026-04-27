package api

import (
	"github.com/gin-gonic/gin"
)

type Module interface {
	Register(rg *gin.RouterGroup)
}

type RouteModules struct {
	Auth      Module
	User      Module
	Product   Module
	Inventory Module
	Order     Module
	Payment   Module
}

func RegisterRoutes(r *gin.Engine, modules RouteModules) {
	v1 := r.Group("/api/v1")
	{
		modules.Auth.Register(v1)
		modules.User.Register(v1)
		modules.Product.Register(v1)
		modules.Inventory.Register(v1)
		modules.Order.Register(v1)
		modules.Payment.Register(v1)
	}
}

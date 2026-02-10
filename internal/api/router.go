package api

import (
	"go-erp/internal/api/v1/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		auth.RegisterAuthRoutes(v1)
	}
}

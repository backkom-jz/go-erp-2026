package auth

import (
	"net/http"

	"go-erp/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/login", login)
	// TODO: refresh, logout
}

func login(c *gin.Context) {
	response.OK(c, http.StatusOK, gin.H{"token": "todo"})
}

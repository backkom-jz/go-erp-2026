package middleware

import (
	"net/http"

	"go-erp/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: parse/validate JWT, set user/tenant in context
		c.Next()
		if c.IsAborted() {
			return
		}
		_ = http.StatusOK
	}
}

func AbortUnauthorized(c *gin.Context, message string) {
	response.Fail(c, http.StatusUnauthorized, 401, message)
	c.Abort()
}

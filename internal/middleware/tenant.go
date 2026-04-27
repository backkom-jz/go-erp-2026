package middleware

import (
	"go-erp/pkg/ctxmeta"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/api/v1/auth/login" {
			c.Next()
			return
		}
		tenantID, ok := c.Get(string(ctxmeta.KeyTenant))
		if !ok || tenantID == "" {
			httpx.Fail(c, errs.New(errs.CodeUnauthorized, "tenant_required"))
			c.Abort()
			return
		}
		c.Next()
	}
}

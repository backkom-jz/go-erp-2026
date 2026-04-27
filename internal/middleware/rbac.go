package middleware

import (
	"go-erp/pkg/ctxmeta"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"
	"strings"

	"github.com/gin-gonic/gin"
)

func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/api/v1/auth/login" {
			c.Next()
			return
		}

		roleValue, _ := c.Get(string(ctxmeta.KeyRole))
		role, _ := roleValue.(string)
		if role == "" {
			role = "viewer"
		}

		if strings.HasPrefix(path, "/api/v1/orders/create") && role == "viewer" {
			httpx.Fail(c, errs.New(errs.CodeForbidden, "permission_denied"))
			c.Abort()
			return
		}
		c.Next()
	}
}

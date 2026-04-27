package middleware

import (
	"go-erp/pkg/auth/jwt"
	"go-erp/pkg/ctxmeta"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWT(manager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/api/v1/auth/login" {
			c.Next()
			return
		}

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			httpx.Fail(c, errs.New(errs.CodeUnauthorized, "missing_token"))
			c.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		claims, err := manager.Parse(token)
		if err != nil {
			httpx.Fail(c, errs.Wrap(errs.CodeUnauthorized, "invalid_token", err))
			c.Abort()
			return
		}

		c.Set(string(ctxmeta.KeyUserID), claims.UserID)
		c.Set(string(ctxmeta.KeyTenant), claims.TenantID)
		c.Set(string(ctxmeta.KeyRole), claims.Role)
		ctx := c.Request.Context()
		ctx = ctxmeta.WithUserID(ctx, claims.UserID)
		ctx = ctxmeta.WithTenantID(ctx, claims.TenantID)
		ctx = ctxmeta.WithRole(ctx, claims.Role)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AbortUnauthorized(c *gin.Context, message string) {
	httpx.Fail(c, errs.New(errs.CodeUnauthorized, message))
	c.Abort()
}

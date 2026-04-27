package middleware

import (
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"
	"go-erp/pkg/idempotency"
	"time"

	"github.com/gin-gonic/gin"
)

func Idempotency(store *idempotency.Store, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if store == nil {
			c.Next()
			return
		}
		key := c.GetHeader("Idempotency-Key")
		if key == "" {
			c.Next()
			return
		}
		if err := store.Reserve(c.Request.Context(), "idem:"+key, ttl); err != nil {
			httpx.Fail(c, errs.New(errs.CodeDuplicate, "duplicate_request"))
			c.Abort()
			return
		}
		c.Next()
	}
}

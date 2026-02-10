package middleware

import "github.com/gin-gonic/gin"

func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: extract tenant_id from context/JWT and set for DB scope
		c.Next()
	}
}

package middleware

import "github.com/gin-gonic/gin"

func RBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: implement RBAC check
		c.Next()
	}
}

package middleware

import (
	"go-erp/pkg/ctxmeta"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func Logger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceID := uuid.NewString()
		c.Set(string(ctxmeta.KeyTraceID), traceID)
		c.Request = c.Request.WithContext(ctxmeta.WithTraceID(c.Request.Context(), traceID))

		c.Next()

		latency := time.Since(start)
		log.Info("http_request",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

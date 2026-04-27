package ctxmeta

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	KeyTraceID ctxKey = "trace_id"
	KeyUserID  ctxKey = "user_id"
	KeyTenant  ctxKey = "tenant_id"
	KeyRole    ctxKey = "role"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, KeyTraceID, traceID)
}

func TraceID(ctx context.Context) string {
	if v, ok := ctx.Value(KeyTraceID).(string); ok {
		return v
	}
	return ""
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, KeyUserID, userID)
}

func UserID(ctx context.Context) string {
	if v, ok := ctx.Value(KeyUserID).(string); ok {
		return v
	}
	return ""
}

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, KeyTenant, tenantID)
}

func TenantID(ctx context.Context) string {
	if v, ok := ctx.Value(KeyTenant).(string); ok {
		return v
	}
	return ""
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, KeyRole, role)
}

func Role(ctx context.Context) string {
	if v, ok := ctx.Value(KeyRole).(string); ok {
		return v
	}
	return ""
}

func GetTraceIDFromGin(c *gin.Context) string {
	if v, ok := c.Get(string(KeyTraceID)); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

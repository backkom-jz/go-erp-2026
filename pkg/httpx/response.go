package httpx

import (
	"go-erp/pkg/ctxmeta"
	"go-erp/pkg/errs"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    int(errs.CodeOK),
		Message: "ok",
		Data:    data,
		TraceID: ctxmeta.GetTraceIDFromGin(c),
	})
}

func Fail(c *gin.Context, err error) {
	status, code, message := errs.ToHTTP(err)
	c.JSON(status, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: ctxmeta.GetTraceIDFromGin(c),
	})
}

type PageData[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}

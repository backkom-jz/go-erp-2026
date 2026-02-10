package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceID string      `json:"trace_id"`
}

func OK(c *gin.Context, httpStatus int, data interface{}) {
	c.JSON(httpStatus, Response{
		Code:    0,
		Message: "ok",
		Data:    data,
		TraceID: getTraceID(c),
	})
}

func Fail(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: getTraceID(c),
	})
}

func getTraceID(c *gin.Context) string {
	if v, ok := c.Get("trace_id"); ok {
		if s, ok2 := v.(string); ok2 {
			return s
		}
	}
	return ""
}

func HTTPStatusFromCode(code int) int {
	if code == 0 {
		return http.StatusOK
	}
	return http.StatusBadRequest
}

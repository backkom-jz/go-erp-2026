package response

import (
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

func OK(c *gin.Context, httpStatus int, data interface{}) {
	_ = httpStatus
	httpx.OK(c, data)
}

func Fail(c *gin.Context, httpStatus int, code int, message string) {
	_ = httpStatus
	httpx.Fail(c, errs.New(errs.Code(code), message))
}

func HTTPStatusFromCode(code int) int {
	status, _, _ := errs.ToHTTP(errs.New(errs.Code(code), ""))
	return status
}

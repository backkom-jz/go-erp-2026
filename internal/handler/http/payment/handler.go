package payment

import (
	dtopayment "go-erp/internal/dto/payment"
	paymentsvc "go-erp/internal/service/payment"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *paymentsvc.Service
}

// NewHandler 创建支付模块处理器。
func NewHandler(service *paymentsvc.Service) *Handler {
	return &Handler{service: service}
}

// Register 注册支付接口。
// 接口备注：
// - POST /api/v1/payments/callback 支付回调
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/payments/callback", h.Callback)
}

// Callback 支付回调接口。
func (h *Handler) Callback(c *gin.Context) {
	var req dtopayment.CallbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_payment_callback", err))
		return
	}
	if err := h.service.Callback(c.Request.Context(), req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"processed": true})
}

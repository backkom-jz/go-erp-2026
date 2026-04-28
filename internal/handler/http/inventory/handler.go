package inventory

import (
	dtoinventory "go-erp/internal/dto/inventory"
	inventorysvc "go-erp/internal/service/inventory"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *inventorysvc.Service
}

// NewHandler 创建库存模块处理器。
func NewHandler(service *inventorysvc.Service) *Handler {
	return &Handler{service: service}
}

// Register 注册库存接口。
// 接口备注：
// - POST /api/v1/inventory/deduct 扣减库存
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/inventory/deduct", h.Deduct)
}

// Deduct 扣减库存接口。
func (h *Handler) Deduct(c *gin.Context) {
	var req dtoinventory.DeductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_inventory_request", err))
		return
	}
	if err := h.service.Deduct(c.Request.Context(), req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"deducted": true})
}

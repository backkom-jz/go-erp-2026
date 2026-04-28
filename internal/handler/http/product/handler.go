package product

import (
	"go-erp/internal/dto/product"
	productsvc "go-erp/internal/service/product"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *productsvc.Service
}

// NewHandler 创建商品模块处理器。
func NewHandler(service *productsvc.Service) *Handler {
	return &Handler{service: service}
}

// Register 注册商品接口。
// 接口备注：
// - POST /api/v1/products/spu 创建 SPU
// - GET  /api/v1/products/spu 查询 SPU 列表
// - POST /api/v1/products/sku 创建 SKU
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/products/spu", h.CreateSPU)
	rg.GET("/products/spu", h.ListSPU)
	rg.POST("/products/sku", h.CreateSKU)
}

// CreateSPU 创建 SPU。
func (h *Handler) CreateSPU(c *gin.Context) {
	var req product.CreateSPURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_spu_request", err))
		return
	}
	row, err := h.service.CreateSPU(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, row)
}

// ListSPU 查询 SPU 列表。
func (h *Handler) ListSPU(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.service.ListSPU(c.Request.Context(), limit)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, rows)
}

// CreateSKU 创建 SKU。
func (h *Handler) CreateSKU(c *gin.Context) {
	var req product.CreateSKURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_sku_request", err))
		return
	}
	row, err := h.service.CreateSKU(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, row)
}

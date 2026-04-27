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

func NewHandler(service *productsvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/products/spu", h.CreateSPU)
	rg.GET("/products/spu", h.ListSPU)
	rg.POST("/products/sku", h.CreateSKU)
}

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

func (h *Handler) ListSPU(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.service.ListSPU(c.Request.Context(), limit)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, rows)
}

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

package order

import (
	dtoorder "go-erp/internal/dto/order"
	ordersvc "go-erp/internal/service/order"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *ordersvc.Service
}

func NewHandler(service *ordersvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/order/create", h.Create)
	rg.GET("/order/:id", h.Get)
}

func (h *Handler) Create(c *gin.Context) {
	var req dtoorder.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_order_request", err))
		return
	}
	row, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, row)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_order_id", err))
		return
	}
	header, items, getErr := h.service.GetByID(c.Request.Context(), uint(id))
	if getErr != nil {
		httpx.Fail(c, getErr)
		return
	}
	httpx.OK(c, gin.H{
		"order": header,
		"items": items,
	})
}

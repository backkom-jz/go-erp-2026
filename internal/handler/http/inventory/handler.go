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

func NewHandler(service *inventorysvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/inventory/deduct", h.Deduct)
}

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

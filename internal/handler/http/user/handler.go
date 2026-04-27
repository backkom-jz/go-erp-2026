package user

import (
	dtouser "go-erp/internal/dto/user"
	usersvc "go-erp/internal/service/user"
	"go-erp/pkg/ctxmeta"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *usersvc.Service
}

func NewHandler(service *usersvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/users", h.Create)
	rg.GET("/users/me", h.Me)
}

func (h *Handler) Create(c *gin.Context) {
	var req dtouser.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_user_request", err))
		return
	}
	if err := h.service.Create(c.Request.Context(), req); err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, gin.H{"created": true})
}

func (h *Handler) Me(c *gin.Context) {
	userNo := ctxmeta.UserID(c.Request.Context())
	if userNo == "" {
		httpx.Fail(c, errs.New(errs.CodeUnauthorized, "missing_user"))
		return
	}
	entity, err := h.service.GetByUserNo(c.Request.Context(), userNo)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, entity)
}

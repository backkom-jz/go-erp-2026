package auth

import (
	dtoauth "go-erp/internal/dto/auth"
	authsvc "go-erp/internal/service/auth"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *authsvc.Service
}

func NewHandler(service *authsvc.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/auth/login", h.Login)
}

func (h *Handler) Login(c *gin.Context) {
	var req dtoauth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_login_request", err))
		return
	}
	if req.Role == "" {
		req.Role = "admin"
	}
	access, refresh, err := h.service.Login(c.Request.Context(), req.UserNo, req.TenantID, req.Role)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, dtoauth.LoginResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

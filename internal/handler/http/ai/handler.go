package ai

import (
	dtoai "go-erp/internal/dto/ai"
	aisvc "go-erp/internal/service/ai"
	"go-erp/pkg/errs"
	"go-erp/pkg/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *aisvc.Service
}

// NewHandler 创建 AI 模块处理器。
func NewHandler(service *aisvc.Service) *Handler {
	return &Handler{service: service}
}

// Register 注册 AI 接口。
// 接口备注：
// - POST /api/v1/ai/chat AI 对话（DeepSeekV4-pro）
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.POST("/ai/chat", h.Chat)
}

// Chat AI 对话接口。
func (h *Handler) Chat(c *gin.Context) {
	var req dtoai.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.Fail(c, errs.Wrap(errs.CodeBadRequest, "invalid_ai_chat_request", err))
		return
	}
	reply, err := h.service.Chat(c.Request.Context(), req)
	if err != nil {
		httpx.Fail(c, err)
		return
	}
	httpx.OK(c, dtoai.ChatResponse{Reply: reply})
}

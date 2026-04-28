package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"net/http"
	"strings"
	"time"

	dtoai "go-erp/internal/dto/ai"
	"go-erp/pkg/errs"
)

type Config struct {
	Enabled        bool
	BaseURL        string
	APIKey         string
	Model          string
	TimeoutSeconds int
	Temperature    float64
	MaxTokens      int
}

type Service struct {
	client *http.Client
	cfg    Config
}

type deepSeekRequest struct {
	Model       string          `json:"model"`
	Messages    []dtoai.Message `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Stream      bool            `json:"stream"`
}

type deepSeekResponse struct {
	Choices []struct {
		Message dtoai.Message `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewService 创建 AI 对话服务。
func NewService(cfg Config) *Service {
	timeout := cfg.TimeoutSeconds
	if timeout <= 0 {
		timeout = 60
	}
	return &Service{
		client: &http.Client{Timeout: time.Duration(timeout) * time.Second},
		cfg:    cfg,
	}
}

// Chat 通过后端代理调用 DeepSeekV4-pro 对话接口。
func (s *Service) Chat(ctx context.Context, req dtoai.ChatRequest) (string, error) {
	if !s.cfg.Enabled {
		return "", errs.New(errs.CodeBadRequest, "ai_service_disabled")
	}
	apiKey := strings.TrimSpace(s.cfg.APIKey)
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", errs.New(errs.CodeInternal, "ai_api_key_missing")
	}

	model := s.cfg.Model
	if strings.TrimSpace(model) == "" {
		model = "DeepSeekV4-pro"
	}
	// DeepSeek 官方兼容模型名映射，避免因别名导致请求失败。
	switch strings.ToLower(strings.TrimSpace(model)) {
	case "deepseekv4-pro", "deepseek-v4-pro":
		model = "deepseek-chat"
	}
	baseURL := strings.TrimSpace(s.cfg.BaseURL)
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/chat/completions"
	}
	baseURL = normalizeDeepSeekURL(baseURL)

	messages := make([]dtoai.Message, 0, len(req.History)+1)
	for _, item := range req.History {
		role := strings.TrimSpace(item.Role)
		content := strings.TrimSpace(item.Content)
		if role == "" || content == "" {
			continue
		}
		messages = append(messages, dtoai.Message{Role: role, Content: content})
	}
	messages = append(messages, dtoai.Message{Role: "user", Content: strings.TrimSpace(req.Message)})

	payload, err := json.Marshal(deepSeekRequest{
		Model:       model,
		Messages:    messages,
		Temperature: s.cfg.Temperature,
		MaxTokens:   s.cfg.MaxTokens,
		Stream:      false,
	})
	if err != nil {
		return "", errs.Wrap(errs.CodeInternal, "marshal_ai_request_failed", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL, bytes.NewReader(payload))
	if err != nil {
		return "", errs.Wrap(errs.CodeInternal, "create_ai_request_failed", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", errs.Wrap(errs.CodeInternal, "call_ai_failed", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errs.Wrap(errs.CodeInternal, "read_ai_response_failed", err)
	}

	var parsed deepSeekResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", errs.Wrap(errs.CodeInternal, "decode_ai_response_failed", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		if resp.StatusCode == http.StatusUnauthorized {
			return "", errs.New(errs.CodeUnauthorized, "ai_invalid_token")
		}
		if parsed.Error != nil && parsed.Error.Message != "" {
			msg := strings.ToLower(parsed.Error.Message)
			if strings.Contains(msg, "invalid") && strings.Contains(msg, "token") {
				return "", errs.New(errs.CodeUnauthorized, "ai_invalid_token")
			}
			if strings.Contains(msg, "invalid") && strings.Contains(msg, "key") {
				return "", errs.New(errs.CodeUnauthorized, "ai_invalid_token")
			}
			return "", errs.New(errs.CodeInternal, "ai_error:"+parsed.Error.Message)
		}
		return "", errs.New(errs.CodeInternal, fmt.Sprintf("ai_http_%d", resp.StatusCode))
	}
	if len(parsed.Choices) == 0 {
		return "", errs.Wrap(errs.CodeInternal, "ai_empty_choices", errors.New("no choices"))
	}

	reply := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if reply == "" {
		return "", errs.Wrap(errs.CodeInternal, "ai_empty_reply", errors.New("empty content"))
	}
	return reply, nil
}

func normalizeDeepSeekURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/chat/completions"
	}
	return u.String()
}

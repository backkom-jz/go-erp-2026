package ai

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message string    `json:"message" binding:"required"`
	History []Message `json:"history"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

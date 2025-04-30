package message_model

import "github.com/cloudwego/eino/schema"

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message,omitempty"`
	Data    *schema.Message `json:"data,omitempty"`
}

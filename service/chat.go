package service

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
)

// ChatService 聊天服务接口
type ChatService interface {
	SendMessage(ctx context.Context, message string) (string, error)
}

// chatService 聊天服务实现
type chatService struct {
	chatModel *ollama.ChatModel
}

// NewChatService 创建聊天服务实例
func NewChatService() ChatService {
	return &chatService{
		chatModel: global.ChatModel,
	}
}

// SendMessage 发送聊天消息
func (s *chatService) SendMessage(ctx context.Context, message string) (string, error) {
	// TODO: 实现消息发送逻辑
	// 1. 调用模型进行对话
	// 2. 处理响应
	// 3. 返回结果
	return "", nil
}

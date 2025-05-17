package pipeline

import (
	"MoonAgent/internal/global"
	"context"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

func NewChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {

	config := &ark.ChatModelConfig{
		APIKey: global.ServerConfig.LLMConfig.API_KEY,
		Model:  global.ServerConfig.LLMConfig.MODEL,
	}

	chatModel, err := ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}

	return chatModel, nil
}

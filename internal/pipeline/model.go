package pipeline

import (
	"context"

	"MoonAgent/internal/global"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context) (cm model.ToolCallingChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ark.ChatModelConfig{
		APIKey: global.ServerConfig.LLMConfig.API_KEY,
		Model:  global.ServerConfig.LLMConfig.MODEL,
	}
	cm, err = ark.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

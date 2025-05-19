package tests

import (
	"context"

	"MoonAgent/pkg/config"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context, config *config.ServerConfig) (cm model.ToolCallingChatModel, err error) {
	// TODO Modify component configuration here.
	modelConfig := &ark.ChatModelConfig{
		APIKey: config.LLMConfig.API_KEY,
		Model:  config.LLMConfig.MODEL,
	}
	cm, err = ark.NewChatModel(ctx, modelConfig)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

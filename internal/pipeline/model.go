package pipeline

import (
	"context"

	"MoonAgent/cmd/di"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/model"
)

func newChatModel(ctx context.Context, app *di.Application) (cm model.ToolCallingChatModel, err error) {
	// TODO Modify component configuration here.
	modelConfig := &ark.ChatModelConfig{
		APIKey: app.ServerConfig.LLMConfig.API_KEY,
		Model:  app.ServerConfig.LLMConfig.MODEL,
	}
	cm, err = ark.NewChatModel(ctx, modelConfig)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

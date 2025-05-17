package pipeline

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type ToolImpl struct {
	config *ToolConfig
}

type ToolConfig struct {
}

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	// config := &ToolConfig{}
	// bt = &ToolImpl{config: config}
	// return bt, nil
	googleAPIKey := "AIzaSyDrKLkc290NdtNfC8fkOVQTVCPq_yuXZpA"
	googleSearchEngineID := "60af7f26ff64c4d55"
	tool, err := googlesearch.NewTool(ctx, &googlesearch.Config{
		APIKey:         googleAPIKey,
		SearchEngineID: googleSearchEngineID,
	})
	if err != nil {
		return nil, err
	}
	return tool, nil
}

func (impl *ToolImpl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	panic("implement me")
}

func (impl *ToolImpl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	panic("implement me")
}

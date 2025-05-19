package tests

import (
	"context"
	"encoding/json"

	"MoonAgent/pkg/utils"

	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func newGoogleSearchTool(ctx context.Context) (bt tool.BaseTool, err error) {
	config := &googlesearch.Config{
		APIKey:         "AIzaSyDrKLkc290NdtNfC8fkOVQTVCPq_yuXZpA",
		SearchEngineID: "60af7f26ff64c4d55",
	}
	bt, err = googlesearch.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}

type JumpWebPageImpl struct {
	config *JumpWebPageConfig
}

type JumpWebPageConfig struct {
}

func newJumpWebPage(ctx context.Context) (bt tool.BaseTool, err error) {
	config := &JumpWebPageConfig{}
	bt = &JumpWebPageImpl{config: config}
	return bt, nil
}

func (impl *JumpWebPageImpl) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "网页跳转",
		Desc: "跳转到指定网页",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"url": {
				Type: "string",
				Desc: "网页地址",
			},
		}),
	}, nil
}

func (impl *JumpWebPageImpl) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	p := &GoToWebPageParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}
	if p.URL != "" {
		return utils.GoToWebPage(ctx, p.URL)
	}
	return "", nil
}

type GoToWebPageParam struct {
	URL string `json:"url"`
}

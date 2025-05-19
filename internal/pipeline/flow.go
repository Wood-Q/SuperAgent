package pipeline

import (
	"context"

	"MoonAgent/pkg/config"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// newLambda component initialization function of node 'Lambda3' in graph 'Assitant'
func newLambda(ctx context.Context, config *config.ServerConfig) (lba *compose.Lambda, err error) {

	messageModifier := func(ctx context.Context, input []*schema.Message) []*schema.Message {
		res := make([]*schema.Message, 0, len(input)+1)
		res = append(res, schema.SystemMessage("你是一个热情的小助手，会活泼地回答问题"))
		res = append(res, input...)
		return res
	}

	agentConfig := &react.AgentConfig{
		MessageModifier: messageModifier,
	}

	chatModelIns11, err := newChatModel(ctx, config)
	if err != nil {
		return nil, err
	}

	agentConfig.ToolCallingModel = chatModelIns11

	toolIns21, err := newGoogleSearchTool(ctx)

	toolIns22, err := newJumpWebPage(ctx)

	if err != nil {
		return nil, err
	}

	agentConfig.ToolsConfig.Tools = []tool.BaseTool{toolIns21, toolIns22}

	ins, err := react.NewAgent(ctx, agentConfig)
	if err != nil {
		return nil, err
	}

	lba, err = compose.AnyLambda(ins.Generate, ins.Stream, nil, nil)

	if err != nil {
		return nil, err
	}

	return lba, nil
}

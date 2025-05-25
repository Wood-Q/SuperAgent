package toolcallagent

import (
	reactagent "MoonAgent/internal/agents/ReActAgent"
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

type ToolCallAgent struct {
	ReActAgent *reactagent.ReActAgent
	Tools      []schema.ToolInfo
}

func NewToolCallAgent(name string, systemPrompt string, nextPrompt string, chatModel *model.BaseChatModel, tools []schema.ToolInfo) *ToolCallAgent {
	ta := &ToolCallAgent{
		ReActAgent: reactagent.NewReActAgent(name, systemPrompt, nextPrompt, chatModel),
		Tools:      tools,
	}
	ta.ReActAgent.ThinkFunc = ta.Think
	ta.ReActAgent.ActFunc = ta.Act
	return ta
}

func (ta *ToolCallAgent) Think(ctx context.Context, input string) (bool, error) {
	if len(ta.Tools) == 0 {
		return false, errors.New("no tools available")
	}
	toolNames := make([]string, 0)
	toolDescriptions := make([]string, 0)
	for _, tool := range ta.Tools {
		toolNames = append(toolNames, tool.Name)
		toolDescriptions = append(toolDescriptions, tool.Desc)
	}
	input += "\n\n" + "You have the following tools available: " + strings.Join(toolNames, ", ") + "\n\n" + strings.Join(toolDescriptions, "\n")
	chatResp, err := ta.ReActAgent.BaseAgent.Run(ctx, input)
	if err != nil {
		return false, err
	}
	if chatResp.ToolCalls != nil {
		zap.S().Info("思考内容为：", chatResp.Content)
		zap.S().Info("工具调用为：", chatResp.ToolCalls)
		return true, nil
	} else {
		zap.S().Info("思考内容为：", chatResp.Content)
		return false, nil
	}
}

func (ta *ToolCallAgent) Act(ctx context.Context, input string) (string, error) {
	if len(ta.Tools) == 0 {
		return "", errors.New("no tools available")
	}
	toolNames := make([]string, 0)
	for _, tool := range ta.Tools {
		toolNames = append(toolNames, tool.Name)
	}
	input += "\n\n" + "You have the following tools available: " + strings.Join(toolNames, ", ")
	chatResp, err := ta.ReActAgent.BaseAgent.Run(ctx, input)
	if err != nil {
		return "", err
	}
	return chatResp.Content, nil
}

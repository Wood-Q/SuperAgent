package reactagent

import (
	baseagent "MoonAgent/internal/agents/baseAgent"
	"context"

	"github.com/cloudwego/eino/components/model"
)

type ReActAgentInterface interface {
	Think(ctx context.Context, input string) (bool, error)
	Act(ctx context.Context, input string) (string, error)
}

type ReActAgent struct {
	BaseAgent *baseagent.BaseAgent
	ThinkFunc func(ctx context.Context, input string) (bool, error)
	ActFunc   func(ctx context.Context, input string) (string, error)
}

func NewReActAgent(name string, systemPrompt string, nextPrompt string, chatModel *model.BaseChatModel) *ReActAgent {
	ra := &ReActAgent{
		BaseAgent: baseagent.NewBaseAgent(name, systemPrompt, nextPrompt, chatModel),
	}
	ra.BaseAgent.StepFunc = ra.Step
	return ra
}

func (ra *ReActAgent) Step(ctx context.Context) (string, error) {
	shoudAct, err := ra.ThinkFunc(ctx, ctx.Value("userPrompt").(string))
	if err != nil {
		return "", err
	}
	if !shoudAct {
		return "No action needed", nil
	}
	return ra.ActFunc(ctx, ctx.Value("userPrompt").(string))
}

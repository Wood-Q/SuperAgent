package baseagent

import (
	"MoonAgent/internal/constants"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

type BaseAgentInterface interface {
	Run(ctx context.Context, input string) (string, error)
	Step(ctx context.Context) (string, error)
}

type BaseAgent struct {
	name         string
	systemPrompt string
	nextPrompt   string

	state       constants.AgentState
	maxSteps    int
	currentStep int

	StepFunc func(ctx context.Context) (string, error)

	chatModel *model.BaseChatModel
}

func NewBaseAgent(name string, systemPrompt string, nextPrompt string, chatModel *model.BaseChatModel) *BaseAgent {
	return &BaseAgent{
		name:         name,
		systemPrompt: systemPrompt,
		nextPrompt:   nextPrompt,
		state:        constants.AgentStateIdle,
		maxSteps:     10,
		currentStep:  0,
		chatModel:    chatModel,
	}
}

func (a *BaseAgent) Run(ctx context.Context, userInput string) (schema.Message, error) {
	if a.state != constants.AgentStateIdle {
		return schema.Message{}, errors.New("agent is not idle")
	}
	//更改状态
	a.state = constants.AgentStateRunning
	//记录上下文
	ctx = context.WithValue(ctx, "userPrompt", userInput)
	//保存结果
	resultsList := make([]string, 0)
	for i := 0; i < a.maxSteps && a.state == constants.AgentStateRunning; i++ {
		stepNumber := i + 1
		a.currentStep = stepNumber
		zap.L().Info("agent is running",
			zap.Int("step", stepNumber),
			zap.Int("maxSteps", a.maxSteps))
		stepResult, err := a.StepFunc(ctx)
		if err != nil {
			return schema.Message{}, err
		}
		result := "Step " + strconv.Itoa(stepNumber) + " result: " + stepResult
		resultsList = append(resultsList, result)
	}
	if a.currentStep >= a.maxSteps {
		a.state = constants.AgentStateSuccess
		resultsList = append(resultsList, "Agent completed all steps")
	}
	return schema.Message{
		Role:    "assistant",
		Content: strings.Join(resultsList, "\n"),
	}, nil
}

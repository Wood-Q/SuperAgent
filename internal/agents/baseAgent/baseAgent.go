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
	Run(ctx context.Context, input string) (*schema.Message, error)
	RunStream(ctx context.Context, input string) (<-chan *schema.Message, error)
	Step(ctx context.Context) (*schema.Message, error)
	Reset()
	GetState() constants.AgentState
	SetMaxSteps(maxSteps int)
}

type BaseAgent struct {
	name         string
	systemPrompt string
	nextPrompt   string

	state       constants.AgentState
	maxSteps    int
	currentStep int
	stepHistory []string

	StepFunc func(ctx context.Context) (*schema.Message, error)

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
		stepHistory:  make([]string, 0),
		chatModel:    chatModel,
	}
}

func (a *BaseAgent) Run(ctx context.Context, userInput string) (*schema.Message, error) {
	if a.state != constants.AgentStateIdle {
		return nil, errors.New("agent is not idle")
	}

	// 重置状态
	a.Reset()
	a.state = constants.AgentStateRunning

	// 记录上下文
	ctx = context.WithValue(ctx, "userPrompt", userInput)

	// 保存结果
	resultsList := make([]string, 0)

	defer func() {
		if a.state == constants.AgentStateRunning {
			a.state = constants.AgentStateSuccess
		}
	}()

	for i := 0; i < a.maxSteps && a.state == constants.AgentStateRunning; i++ {
		stepNumber := i + 1
		a.currentStep = stepNumber

		zap.L().Info("agent is running",
			zap.String("agent", a.name),
			zap.Int("step", stepNumber),
			zap.Int("maxSteps", a.maxSteps))

		stepResult, err := a.StepFunc(ctx)
		if err != nil {
			a.state = constants.AgentStateError
			return nil, err
		}

		if stepResult == nil {
			continue
		}

		result := "Step " + strconv.Itoa(stepNumber) + " result: " + stepResult.Content
		resultsList = append(resultsList, result)
		a.stepHistory = append(a.stepHistory, stepResult.Content)

		// 检查是否应该结束
		if a.shouldStop(stepResult) {
			break
		}
	}

	if a.currentStep >= a.maxSteps {
		a.state = constants.AgentStateSuccess
		resultsList = append(resultsList, "Agent completed all steps")
	}

	return &schema.Message{
		Role:    "assistant",
		Content: strings.Join(resultsList, "\n"),
	}, nil
}

func (a *BaseAgent) RunStream(ctx context.Context, input string) (<-chan *schema.Message, error) {
	if a.state != constants.AgentStateIdle {
		return nil, errors.New("agent is not idle")
	}

	resultChan := make(chan *schema.Message, 10)

	go func() {
		defer close(resultChan)

		// 重置状态
		a.Reset()
		a.state = constants.AgentStateRunning

		// 记录上下文
		ctx = context.WithValue(ctx, "userPrompt", input)

		defer func() {
			if a.state == constants.AgentStateRunning {
				a.state = constants.AgentStateSuccess
			}
		}()

		for i := 0; i < a.maxSteps && a.state == constants.AgentStateRunning; i++ {
			stepNumber := i + 1
			a.currentStep = stepNumber

			stepResult, err := a.StepFunc(ctx)
			if err != nil {
				a.state = constants.AgentStateError
				resultChan <- &schema.Message{
					Role:    "assistant",
					Content: "Error: " + err.Error(),
				}
				return
			}

			if stepResult != nil {
				a.stepHistory = append(a.stepHistory, stepResult.Content)
				resultChan <- stepResult

				// 检查是否应该结束
				if a.shouldStop(stepResult) {
					break
				}
			}
		}

		// 发送完成消息
		resultChan <- &schema.Message{
			Role:    "assistant",
			Content: "Task completed",
		}
	}()

	return resultChan, nil
}

func (a *BaseAgent) Step(ctx context.Context) (*schema.Message, error) {
	if a.StepFunc == nil {
		return nil, errors.New("step function not implemented")
	}
	return a.StepFunc(ctx)
}

func (a *BaseAgent) Reset() {
	a.state = constants.AgentStateIdle
	a.currentStep = 0
	a.stepHistory = make([]string, 0)
}

func (a *BaseAgent) GetState() constants.AgentState {
	return a.state
}

func (a *BaseAgent) SetMaxSteps(maxSteps int) {
	a.maxSteps = maxSteps
}

func (a *BaseAgent) GetName() string {
	return a.name
}

func (a *BaseAgent) GetStepHistory() []string {
	return a.stepHistory
}

// shouldStop 判断是否应该停止执行
func (a *BaseAgent) shouldStop(result *schema.Message) bool {
	// 子类可以重写这个方法来实现自定义的停止逻辑
	return false
}

// GetChatModel 获取聊天模型
func (a *BaseAgent) GetChatModel() *model.BaseChatModel {
	return a.chatModel
}

// GetSystemPrompt 获取系统提示
func (a *BaseAgent) GetSystemPrompt() string {
	return a.systemPrompt
}

// GetNextPrompt 获取下一步提示
func (a *BaseAgent) GetNextPrompt() string {
	return a.nextPrompt
}

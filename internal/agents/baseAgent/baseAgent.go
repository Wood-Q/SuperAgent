package baseagent

import (
	"MoonAgent/internal/constants"
	"MoonAgent/internal/agents/orchestration"
	"errors"
	"strconv"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

type BaseAgentInterface interface {
	//整体运行
	Run(octx *orchestration.OrchestrationContext, input string) (*schema.Message, error)
	//流式运行
	RunStream(octx *orchestration.OrchestrationContext, input string) (<-chan *schema.Message, error)
	//每一步运行
	Step(octx *orchestration.OrchestrationContext) (*schema.Message, error)
	//重置agent状态
	Reset()
	//获取当前agent状态
	GetState() constants.AgentState
	//定义最大步数，防止死循环
	SetMaxSteps(maxSteps int)
}

type BaseAgent struct {
	name         string
	//系统prompt
	systemPrompt string
	//下一步prompt指示
	nextPrompt   string
	//agent状态
	state       constants.AgentState
	//限定最大步数
	maxSteps    int
	//当前步数
	currentStep int
	//每一步记录历史
	stepHistory []string

	StepFunc func(octx *orchestration.OrchestrationContext) (*schema.Message, error)

	chatModel model.ToolCallingChatModel
}

//返回新的BaseAgent结构体
func NewBaseAgent(name string, systemPrompt string, nextPrompt string, chatModel model.ToolCallingChatModel) *BaseAgent {
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

func (a *BaseAgent) Run(octx *orchestration.OrchestrationContext, userInput string) (*schema.Message, error) {
	//如果状态不为空闲，报错
	if a.state != constants.AgentStateIdle {
		return nil, errors.New("agent is not idle")
	}

	// 重置状态
	a.Reset()
	a.state = constants.AgentStateRunning

	// 设置用户输入到编排上下文
	octx.SetInput("userPrompt", userInput)
	octx.AddUserMessage(userInput)

	// 保存结果
	resultsList := make([]string, 0)

	//切换模型状态
	defer func() {
		if a.state == constants.AgentStateRunning {
			a.state = constants.AgentStateSuccess
		}
	}()
	
	//重复步骤运行
	for i := 0; i < a.maxSteps && a.state == constants.AgentStateRunning; i++ {
		stepNumber := i + 1
		a.currentStep = stepNumber

		zap.L().Info("agent is running",
			zap.String("agent", a.name),
			zap.Int("step", stepNumber),
			zap.Int("maxSteps", a.maxSteps))

		stepResult, err := a.StepFunc(octx)
		if err != nil {
			a.state = constants.AgentStateError
			return nil, err
		}

		if stepResult == nil {
			continue
		}

		//存入resultsList和stepHistory
		result := "Step " + strconv.Itoa(stepNumber) + " result: " + stepResult.Content
		resultsList = append(resultsList, result)
		a.stepHistory = append(a.stepHistory, stepResult.Content)

		// 将步骤结果添加到内存
		octx.AddAssistantMessage(stepResult.Content)

		// 检查是否应该结束
		if a.shouldStop(stepResult) {
			break
		}
	}

	if a.currentStep >= a.maxSteps {
		a.state = constants.AgentStateSuccess
		resultsList = append(resultsList, "Agent completed all steps")
	}

	//最终输出的相应
	finalMessage := &schema.Message{
		Role:    "assistant",
		Content: strings.Join(resultsList, "\n"),
	}

	// 添加最终响应到内存
	octx.AddAssistantMessage(finalMessage.Content)

	return finalMessage, nil
}

//流式实现
func (a *BaseAgent) RunStream(octx *orchestration.OrchestrationContext, input string) (<-chan *schema.Message, error) {
	if a.state != constants.AgentStateIdle {
		return nil, errors.New("agent is not idle")
	}

	resultChan := make(chan *schema.Message, 10)

	go func() {
		defer close(resultChan)

		// 重置状态
		a.Reset()
		a.state = constants.AgentStateRunning

		// 设置用户输入到编排上下文
		octx.SetInput("userPrompt", input)
		octx.AddUserMessage(input)

		defer func() {
			if a.state == constants.AgentStateRunning {
				a.state = constants.AgentStateSuccess
			}
		}()

		for i := 0; i < a.maxSteps && a.state == constants.AgentStateRunning; i++ {
			stepNumber := i + 1
			a.currentStep = stepNumber

			stepResult, err := a.StepFunc(octx)
			if err != nil {
				a.state = constants.AgentStateError
				errorMessage := &schema.Message{
					Role:    "assistant",
					Content: "Error: " + err.Error(),
				}
				octx.AddAssistantMessage(errorMessage.Content)
				resultChan <- errorMessage
				return
			}

			if stepResult != nil {
				a.stepHistory = append(a.stepHistory, stepResult.Content)
				octx.AddAssistantMessage(stepResult.Content)
				resultChan <- stepResult

				// 检查是否应该结束
				if a.shouldStop(stepResult) {
					break
				}
			}
		}

		// 发送完成消息
		completionMessage := &schema.Message{
			Role:    "assistant",
			Content: "Task completed",
		}
		octx.AddAssistantMessage(completionMessage.Content)
		resultChan <- completionMessage
	}()

	return resultChan, nil
}

func (a *BaseAgent) Step(octx *orchestration.OrchestrationContext) (*schema.Message, error) {
	if a.StepFunc == nil {
		return nil, errors.New("step function not implemented")
	}
	return a.StepFunc(octx)
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
func (a *BaseAgent) GetChatModel() model.ToolCallingChatModel {
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

package reactagent

import (
	baseagent "MoonAgent/internal/agents/baseAgent"
	"MoonAgent/internal/agents/orchestration"
	"errors"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

type ReActAgentInterface interface {
	//ReAct所需的思考行为
	Think(octx *orchestration.OrchestrationContext) (*schema.Message, error)
	//ReAct所需的行动行为
	Act(octx *orchestration.OrchestrationContext, action string) (*schema.Message, error)
	//ReAct所需的观察行为
	Observe(octx *orchestration.OrchestrationContext, observation string) (*schema.Message, error)
}

type ReActStep struct {
	StepType    string // "think", "act", "observe"
	Content     string
	Action      string
	Observation string
}

type ReActAgent struct {
	BaseAgent *baseagent.BaseAgent

	// ReAct specific fields
	thoughts     []string
	actions      []string
	observations []string
	currentLoop  int
	maxLoops     int

	// Custom functions for ReAct cycle
	ThinkFunc   func(octx *orchestration.OrchestrationContext, history []ReActStep) (*schema.Message, error)
	ActFunc     func(octx *orchestration.OrchestrationContext, thought string) (*schema.Message, error)
	ObserveFunc func(octx *orchestration.OrchestrationContext, action string) (*schema.Message, error)
}

func NewReActAgent(name string, systemPrompt string, nextPrompt string, chatModel model.ToolCallingChatModel) *ReActAgent {
	ra := &ReActAgent{
		BaseAgent:    baseagent.NewBaseAgent(name, systemPrompt, nextPrompt, chatModel),
		thoughts:     make([]string, 0),
		actions:      make([]string, 0),
		observations: make([]string, 0),
		currentLoop:  0,
		maxLoops:     5, // 默认最多5个循环
	}
	ra.BaseAgent.StepFunc = ra.Step
	return ra
}

func (ra *ReActAgent) Step(octx *orchestration.OrchestrationContext) (*schema.Message, error) {
	if ra.currentLoop >= ra.maxLoops {
		return &schema.Message{
			Role:    "assistant",
			Content: "ReAct循环已达到最大次数，结束执行",
		}, nil
	}

	// 1. Think - 推理阶段
	thinkResult, err := ra.Think(octx)
	if err != nil {
		return nil, err
	}

	//思考内容合并
	if thinkResult != nil && thinkResult.Content != "" {
		ra.thoughts = append(ra.thoughts, thinkResult.Content)
		zap.L().Info("ReAct Think",
			zap.Int("loop", ra.currentLoop+1),
			zap.String("thought", thinkResult.Content))

		// 检查是否需要采取行动
		if ra.needsAction(thinkResult.Content) {
			
			// 2. Act - 行动阶段
			actResult, err := ra.Act(octx, thinkResult.Content)
			if err != nil {
				return nil, err
			}

			if actResult != nil && actResult.Content != "" {
				ra.actions = append(ra.actions, actResult.Content)
				zap.L().Info("ReAct Act",
					zap.Int("loop", ra.currentLoop+1),
					zap.String("action", actResult.Content))

				// 3. Observe - 观察阶段
				observeResult, err := ra.Observe(octx, actResult.Content)
				if err != nil {
					return nil, err
				}

				if observeResult != nil && observeResult.Content != "" {
					ra.observations = append(ra.observations, observeResult.Content)
					zap.L().Info("ReAct Observe",
						zap.Int("loop", ra.currentLoop+1),
						zap.String("observation", observeResult.Content))
				}
			}
		}
	}

	ra.currentLoop++

	// 构建完整的响应
	return ra.buildCompleteResponse(), nil
}

func (ra *ReActAgent) Think(octx *orchestration.OrchestrationContext) (*schema.Message, error) {
	if ra.ThinkFunc != nil {
		history := ra.buildHistory()
		return ra.ThinkFunc(octx, history)
	}

	// 默认的思考实现
	userInput, exists := octx.GetInputString("userPrompt")
	if !exists {
		return nil, errors.New("userPrompt not found in orchestration context")
	}

	prompt := ra.buildThinkPrompt(userInput)

	resp, err := ra.BaseAgent.GetChatModel().Generate(octx.Context(), []*schema.Message{
		{Role: "system", Content: ra.BaseAgent.GetSystemPrompt()},
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ra *ReActAgent) Act(octx *orchestration.OrchestrationContext, thought string) (*schema.Message, error) {
	if ra.ActFunc != nil {
		return ra.ActFunc(octx, thought)
	}

	// 默认的行动实现
	prompt := ra.buildActPrompt(thought)

	resp, err := ra.BaseAgent.GetChatModel().Generate(octx.Context(), []*schema.Message{
		{Role: "system", Content: ra.BaseAgent.GetSystemPrompt()},
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ra *ReActAgent) Observe(octx *orchestration.OrchestrationContext, action string) (*schema.Message, error) {
	if ra.ObserveFunc != nil {
		return ra.ObserveFunc(octx, action)
	}

	// 默认的观察实现
	prompt := ra.buildObservePrompt(action)

	resp, err := ra.BaseAgent.GetChatModel().Generate(octx.Context(), []*schema.Message{
		{Role: "system", Content: ra.BaseAgent.GetSystemPrompt()},
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 构建历史记录
func (ra *ReActAgent) buildHistory() []ReActStep {
	history := make([]ReActStep, 0)
	maxLen := len(ra.thoughts)

	for i := 0; i < maxLen; i++ {
		step := ReActStep{StepType: "think"}
		if i < len(ra.thoughts) {
			step.Content = ra.thoughts[i]
		}
		if i < len(ra.actions) {
			step.Action = ra.actions[i]
		}
		if i < len(ra.observations) {
			step.Observation = ra.observations[i]
		}
		history = append(history, step)
	}

	return history
}

// 构建思考提示
func (ra *ReActAgent) buildThinkPrompt(userInput string) string {
	var prompt strings.Builder
	prompt.WriteString("用户问题: " + userInput + "\n\n")

	if len(ra.thoughts) > 0 {
		prompt.WriteString("之前的思考过程:\n")
		for i, thought := range ra.thoughts {
			prompt.WriteString("思考" + string(rune(i+1)) + ": " + thought + "\n")
			if i < len(ra.actions) {
				prompt.WriteString("行动" + string(rune(i+1)) + ": " + ra.actions[i] + "\n")
			}
			if i < len(ra.observations) {
				prompt.WriteString("观察" + string(rune(i+1)) + ": " + ra.observations[i] + "\n")
			}
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("现在请进行下一步思考，分析当前情况并决定下一步行动。")
	return prompt.String()
}

// 构建行动提示
func (ra *ReActAgent) buildActPrompt(thought string) string {
	return "基于以下思考内容，请制定具体的行动计划:\n" + thought + "\n\n请明确说明要采取的具体行动。"
}

// 构建观察提示
func (ra *ReActAgent) buildObservePrompt(action string) string {
	return "已执行以下行动:\n" + action + "\n\n请观察和分析这个行动的结果和影响。"
}

// 判断是否需要采取行动
func (ra *ReActAgent) needsAction(thought string) bool {
	// 简单的启发式判断，可以根据实际需求优化
	actionKeywords := []string{"需要", "应该", "计划", "执行", "行动", "查询", "搜索", "调用"}
	thoughtLower := strings.ToLower(thought)

	for _, keyword := range actionKeywords {
		if strings.Contains(thoughtLower, keyword) {
			return true
		}
	}

	return false
}

// 构建完整响应
func (ra *ReActAgent) buildCompleteResponse() *schema.Message {
	var response strings.Builder

	for i := 0; i < len(ra.thoughts); i++ {
		response.WriteString("🤔 思考: " + ra.thoughts[i] + "\n")

		if i < len(ra.actions) {
			response.WriteString("🎯 行动: " + ra.actions[i] + "\n")
		}

		if i < len(ra.observations) {
			response.WriteString("👁️ 观察: " + ra.observations[i] + "\n")
		}

		response.WriteString("\n")
	}

	return &schema.Message{
		Role:    "assistant",
		Content: response.String(),
	}
}

// 重置ReAct状态
func (ra *ReActAgent) Reset() {
	ra.BaseAgent.Reset()
	ra.thoughts = make([]string, 0)
	ra.actions = make([]string, 0)
	ra.observations = make([]string, 0)
	ra.currentLoop = 0
}

// SetMaxLoops 设置最大循环次数
func (ra *ReActAgent) SetMaxLoops(maxLoops int) {
	ra.maxLoops = maxLoops
}

// GetCurrentLoop 获取当前循环次数
func (ra *ReActAgent) GetCurrentLoop() int {
	return ra.currentLoop
}

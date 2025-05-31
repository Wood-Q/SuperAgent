package toolcallagent

import (
	reactagent "MoonAgent/internal/agents/ReActAgent"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

type ToolCallAgent struct {
	ReActAgent *reactagent.ReActAgent
	Tools      []schema.ToolInfo
	toolMap    map[string]schema.ToolInfo
}

func NewToolCallAgent(name string, systemPrompt string, nextPrompt string, chatModel *model.BaseChatModel, tools []schema.ToolInfo) *ToolCallAgent {
	ta := &ToolCallAgent{
		ReActAgent: reactagent.NewReActAgent(name, systemPrompt, nextPrompt, chatModel),
		Tools:      tools,
		toolMap:    make(map[string]schema.ToolInfo),
	}

	// 构建工具映射
	for _, tool := range tools {
		ta.toolMap[tool.Name] = tool
	}

	// 设置自定义的Think和Act函数
	ta.ReActAgent.ThinkFunc = ta.Think
	ta.ReActAgent.ActFunc = ta.Act
	ta.ReActAgent.ObserveFunc = ta.Observe

	return ta
}

func (ta *ToolCallAgent) Think(ctx context.Context, history []reactagent.ReActStep) (*schema.Message, error) {
	if len(ta.Tools) == 0 {
		return nil, errors.New("no tools available")
	}

	userInput := ctx.Value("userPrompt").(string)
	prompt := ta.buildThinkPrompt(userInput, history)

	// 构建包含工具信息的系统提示
	systemPrompt := ta.buildSystemPromptWithTools()

	resp, err := (*ta.ReActAgent.BaseAgent.GetChatModel()).Generate(ctx, []*schema.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: prompt},
	})

	if err != nil {
		return nil, err
	}

	zap.L().Info("ToolCallAgent Think", zap.String("thought", resp.Content))
	return resp, nil
}

func (ta *ToolCallAgent) Act(ctx context.Context, thought string) (*schema.Message, error) {
	// 分析思考内容，确定需要调用的工具
	toolCall := ta.parseToolCall(thought)

	if toolCall == nil {
		// 如果没有工具调用，返回普通响应
		return &schema.Message{
			Role:    "assistant",
			Content: "基于当前思考，我认为不需要使用工具，可以直接回答。",
		}, nil
	}

	// 执行工具调用
	_, err := ta.executeToolCall(ctx, toolCall)
	if err != nil {
		return &schema.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("工具调用失败: %s", err.Error()),
		}, nil
	}

	return &schema.Message{
		Role:      "assistant",
		Content:   fmt.Sprintf("调用工具 %s，参数: %s", toolCall.Function.Name, toolCall.Function.Arguments),
		ToolCalls: []schema.ToolCall{*toolCall},
	}, nil
}

func (ta *ToolCallAgent) Observe(ctx context.Context, action string) (*schema.Message, error) {
	// 观察工具调用的结果
	observation := fmt.Sprintf("工具执行完成。行动内容: %s", action)

	return &schema.Message{
		Role:    "assistant",
		Content: observation,
	}, nil
}

// 构建包含工具信息的系统提示
func (ta *ToolCallAgent) buildSystemPromptWithTools() string {
	basePrompt := ta.ReActAgent.BaseAgent.GetSystemPrompt()

	var toolsDesc strings.Builder
	toolsDesc.WriteString("\n\n你可以使用以下工具:\n")

	for _, tool := range ta.Tools {
		toolsDesc.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Desc))
	}

	toolsDesc.WriteString("\n当你需要使用工具时，请在思考中明确说明要调用的工具名称和参数。")

	return basePrompt + toolsDesc.String()
}

// 构建思考提示
func (ta *ToolCallAgent) buildThinkPrompt(userInput string, history []reactagent.ReActStep) string {
	var prompt strings.Builder
	prompt.WriteString("用户问题: " + userInput + "\n\n")

	if len(history) > 0 {
		prompt.WriteString("历史思考过程:\n")
		for i, step := range history {
			if step.Content != "" {
				prompt.WriteString(fmt.Sprintf("思考%d: %s\n", i+1, step.Content))
			}
			if step.Action != "" {
				prompt.WriteString(fmt.Sprintf("行动%d: %s\n", i+1, step.Action))
			}
			if step.Observation != "" {
				prompt.WriteString(fmt.Sprintf("观察%d: %s\n", i+1, step.Observation))
			}
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("请分析当前情况，决定是否需要使用工具来解决问题。如果需要使用工具，请明确说明工具名称和参数。")

	return prompt.String()
}

// 解析工具调用
func (ta *ToolCallAgent) parseToolCall(thought string) *schema.ToolCall {
	// 简单的工具调用解析逻辑
	// 在实际应用中，可以使用更复杂的NLP技术或让LLM生成结构化的工具调用

	thoughtLower := strings.ToLower(thought)

	for _, tool := range ta.Tools {
		toolNameLower := strings.ToLower(tool.Name)
		if strings.Contains(thoughtLower, toolNameLower) ||
			strings.Contains(thoughtLower, "调用"+toolNameLower) ||
			strings.Contains(thoughtLower, "使用"+toolNameLower) {

			// 构建工具调用
			return &schema.ToolCall{
				ID:   fmt.Sprintf("call_%s_%d", tool.Name, len(ta.ReActAgent.BaseAgent.GetStepHistory())),
				Type: "function",
				Function: schema.FunctionCall{
					Name:      tool.Name,
					Arguments: ta.extractArguments(thought, tool),
				},
			}
		}
	}

	return nil
}

// 提取参数
func (ta *ToolCallAgent) extractArguments(thought string, tool schema.ToolInfo) string {
	// 简化的参数提取逻辑
	// 在实际应用中，应该使用更智能的方法来提取参数

	args := make(map[string]interface{})

	// 根据工具类型设置默认参数
	switch tool.Name {
	case "search", "google_search":
		// 提取搜索关键词
		if strings.Contains(thought, "搜索") {
			// 简单提取搜索词
			args["query"] = "用户查询内容"
		}
	case "browse", "web_browse":
		// 提取URL
		if strings.Contains(thought, "访问") || strings.Contains(thought, "浏览") {
			args["url"] = "https://example.com"
		}
	}

	if argsBytes, err := json.Marshal(args); err == nil {
		return string(argsBytes)
	}

	return "{}"
}

// 执行工具调用
func (ta *ToolCallAgent) executeToolCall(ctx context.Context, toolCall *schema.ToolCall) (string, error) {
	tool, exists := ta.toolMap[toolCall.Function.Name]
	if !exists {
		return "", fmt.Errorf("tool %s not found", toolCall.Function.Name)
	}

	// 这里应该调用实际的工具执行逻辑
	// 目前返回模拟结果
	result := fmt.Sprintf("工具 %s 执行成功，参数: %s", tool.Name, toolCall.Function.Arguments)

	zap.L().Info("Tool executed",
		zap.String("tool", tool.Name),
		zap.String("arguments", toolCall.Function.Arguments),
		zap.String("result", result))

	return result, nil
}

// Run 重写Run方法以支持工具调用
func (ta *ToolCallAgent) Run(ctx context.Context, input string) (*schema.Message, error) {
	return ta.ReActAgent.BaseAgent.Run(ctx, input)
}

// RunStream 重写RunStream方法以支持流式工具调用
func (ta *ToolCallAgent) RunStream(ctx context.Context, input string) (<-chan *schema.Message, error) {
	return ta.ReActAgent.BaseAgent.RunStream(ctx, input)
}

// Reset 重置状态
func (ta *ToolCallAgent) Reset() {
	ta.ReActAgent.Reset()
}

// GetTools 获取可用工具列表
func (ta *ToolCallAgent) GetTools() []schema.ToolInfo {
	return ta.Tools
}

// AddTool 添加工具
func (ta *ToolCallAgent) AddTool(tool schema.ToolInfo) {
	ta.Tools = append(ta.Tools, tool)
	ta.toolMap[tool.Name] = tool
}

// RemoveTool 移除工具
func (ta *ToolCallAgent) RemoveTool(toolName string) {
	delete(ta.toolMap, toolName)

	// 从切片中移除
	for i, tool := range ta.Tools {
		if tool.Name == toolName {
			ta.Tools = append(ta.Tools[:i], ta.Tools[i+1:]...)
			break
		}
	}
}

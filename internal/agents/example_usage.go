package main

import (
	manus "MoonAgent/internal/agents/Manus"
	reactagent "MoonAgent/internal/agents/ReActAgent"
	toolcallagent "MoonAgent/internal/agents/ToolCallAgent"
	baseagent "MoonAgent/internal/agents/baseAgent"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// ExampleUsage 展示如何使用各种Agent
func main() {
	var chatModel model.ToolCallingChatModel
	// 1. 使用BaseAgent
	fmt.Println("=== BaseAgent 示例 ===")
	baseAgentExample(chatModel)

	// 2. 使用ReActAgent
	fmt.Println("\n=== ReActAgent 示例 ===")
	reactAgentExample(chatModel)

	// 3. 使用ToolCallAgent
	fmt.Println("\n=== ToolCallAgent 示例 ===")
	toolCallAgentExample(chatModel)

	// 4. 使用Manus
	fmt.Println("\n=== Manus 示例 ===")
	manusExample(chatModel)
}

// baseAgentExample BaseAgent使用示例
func baseAgentExample(chatModel model.ToolCallingChatModel) {
	// 创建BaseAgent
	agent := baseagent.NewBaseAgent(
		"BaseAgent示例",
		"你是一个基础AI助手",
		"请继续处理",
		chatModel,
	)

	// 设置自定义步骤函数
	agent.StepFunc = func(ctx context.Context) (*schema.Message, error) {
		userInput := ctx.Value("userPrompt").(string)
		return &schema.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("处理用户输入: %s", userInput),
		}, nil
	}

	// 运行
	ctx := context.Background()
	result, err := agent.Run(ctx, "你好")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("结果: %s\n", result.Content)
	fmt.Printf("状态: %s\n", agent.GetState())
}

// reactAgentExample ReActAgent使用示例
func reactAgentExample(chatModel model.ToolCallingChatModel) {
	// 创建ReActAgent
	agent := reactagent.NewReActAgent(
		"ReActAgent示例",
		"你是一个使用ReAct框架的AI助手",
		"请继续思考和行动",
		chatModel,
	)

	// 设置最大循环次数
	agent.SetMaxLoops(3)

	// 可以自定义Think、Act、Observe函数
	agent.ThinkFunc = func(ctx context.Context, history []reactagent.ReActStep) (*schema.Message, error) {
		userInput := ctx.Value("userPrompt").(string)
		return &schema.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("思考: 用户问了 '%s'，我需要分析这个问题", userInput),
		}, nil
	}

	agent.ActFunc = func(ctx context.Context, thought string) (*schema.Message, error) {
		return &schema.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("行动: 基于思考 '%s'，我决定提供直接回答", thought),
		}, nil
	}

	agent.ObserveFunc = func(ctx context.Context, action string) (*schema.Message, error) {
		return &schema.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("观察: 行动 '%s' 已完成，用户应该得到了满意的回答", action),
		}, nil
	}

	// 运行
	ctx := context.Background()
	result, err := agent.BaseAgent.Run(ctx, "什么是人工智能？")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("结果: %s\n", result.Content)
	fmt.Printf("当前循环: %d\n", agent.GetCurrentLoop())
}

// toolCallAgentExample ToolCallAgent使用示例
func toolCallAgentExample(chatModel model.ToolCallingChatModel) {
	// 定义工具
	tools := []schema.ToolInfo{
		{
			Name: "search",
			Desc: "搜索互联网信息",
		},
		{
			Name: "calculator",
			Desc: "执行数学计算",
		},
	}

	// 创建ToolCallAgent
	agent := toolcallagent.NewToolCallAgent(
		"ToolCallAgent示例",
		"你是一个可以使用工具的AI助手",
		"请继续使用工具处理",
		chatModel,
		tools,
	)

	// 运行
	ctx := context.Background()
	result, err := agent.Run(ctx, "请搜索最新的AI技术发展")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}

	fmt.Printf("结果: %s\n", result.Content)
	fmt.Printf("可用工具: %v\n", agent.GetTools())
}

// manusExample Manus使用示例
func manusExample(chatModel model.ToolCallingChatModel) {
	// 定义工具
	tools := []schema.ToolInfo{
		{
			Name: "web_search",
			Desc: "搜索网页信息",
		},
		{
			Name: "document_query",
			Desc: "查询文档内容",
		},
		{
			Name: "code_generator",
			Desc: "生成代码",
		},
	}

	// 方式1: 使用默认配置
	manus1 := manus.NewManusWithDefaults("Manus助手", chatModel, tools)

	// 方式2: 使用自定义配置
	config := &manus.ManusConfig{
		Name: "高级Manus",
		SystemPrompt: `你是一个高级AI助手，专门帮助用户解决复杂问题。
你具备强大的推理能力和丰富的工具集。`,
		NextPrompt:  "请继续深入分析",
		MaxSteps:    15,
		MaxLoops:    8,
		EnableDebug: true,
	}
	manus2 := manus.NewManus(config, chatModel, tools)

	// 使用Manus1
	fmt.Println("--- 使用默认配置的Manus ---")
	ctx := context.Background()
	result1, err := manus1.Run(ctx, "帮我分析一下Go语言的优势")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %s\n", result1.Content)
	}

	// 使用Manus2
	fmt.Println("\n--- 使用自定义配置的Manus ---")
	result2, err := manus2.Run(ctx, "设计一个微服务架构")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
	} else {
		fmt.Printf("结果: %s\n", result2.Content)
	}

	// 获取调试信息
	if debugInfo := manus2.GetDebugInfo(); debugInfo != nil {
		fmt.Printf("调试信息: %+v\n", debugInfo)
	}

	// 动态添加工具
	newTool := schema.ToolInfo{
		Name: "image_generator",
		Desc: "生成图像",
	}
	manus2.AddTool(newTool)
	fmt.Printf("添加工具后的工具列表: %v\n", len(manus2.GetTools()))

	// 流式处理示例
	fmt.Println("\n--- 流式处理示例 ---")
	streamChan, err := manus2.RunStream(ctx, "解释量子计算的原理")
	if err != nil {
		fmt.Printf("流式处理错误: %v\n", err)
		return
	}

	// 处理流式响应
	for message := range streamChan {
		fmt.Printf("流式消息: %s\n", message.Content)
	}
}

// StreamingExample 流式处理的详细示例
func StreamingExample(chatModel model.ToolCallingChatModel) {
	fmt.Println("=== 流式处理详细示例 ===")

	tools := []schema.ToolInfo{
		{Name: "search", Desc: "搜索信息"},
		{Name: "analyze", Desc: "分析数据"},
	}

	manus := manus.NewManusWithDefaults("流式Manus", chatModel, tools)

	ctx := context.Background()
	streamChan, err := manus.RunStream(ctx, "分析当前AI技术的发展趋势")
	if err != nil {
		fmt.Printf("启动流式处理失败: %v\n", err)
		return
	}

	messageCount := 0
	for message := range streamChan {
		messageCount++
		fmt.Printf("[消息 %d] %s: %s\n", messageCount, message.Role, message.Content)

		// 如果有工具调用
		if len(message.ToolCalls) > 0 {
			for _, toolCall := range message.ToolCalls {
				fmt.Printf("  工具调用: %s(%s)\n", toolCall.Function.Name, toolCall.Function.Arguments)
			}
		}
	}

	fmt.Printf("流式处理完成，共收到 %d 条消息\n", messageCount)
	fmt.Printf("最终状态: %s\n", manus.GetState())
}

package manus

import (
	toolcallagent "MoonAgent/internal/agents/ToolCallAgent"
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

// ManusConfig Manus配置
type ManusConfig struct {
	Name         string
	SystemPrompt string
	NextPrompt   string
	MaxSteps     int
	MaxLoops     int
	EnableDebug  bool
}

// Manus 智能助手，基于ToolCallAgent构建
type Manus struct {
	ToolCallAgent *toolcallagent.ToolCallAgent
	config        *ManusConfig
	logger        *zap.Logger
}

// NewManus 创建新的Manus实例
func NewManus(config *ManusConfig, chatModel model.ToolCallingChatModel, tools []schema.ToolInfo) *Manus {
	if config == nil {
		config = DefaultManusConfig()
	}

	manus := &Manus{
		ToolCallAgent: toolcallagent.NewToolCallAgent(
			config.Name,
			config.SystemPrompt,
			config.NextPrompt,
			chatModel,
			tools,
		),
		config: config,
		logger: zap.L().Named("manus"),
	}

	// 配置参数
	manus.ToolCallAgent.ReActAgent.BaseAgent.SetMaxSteps(config.MaxSteps)
	manus.ToolCallAgent.ReActAgent.SetMaxLoops(config.MaxLoops)

	return manus
}

// NewManusWithDefaults 使用默认配置创建Manus
func NewManusWithDefaults(name string, chatModel model.ToolCallingChatModel, tools []schema.ToolInfo) *Manus {
	config := DefaultManusConfig()
	config.Name = name
	return NewManus(config, chatModel, tools)
}

// DefaultManusConfig 默认配置
func DefaultManusConfig() *ManusConfig {
	return &ManusConfig{
		Name: "Manus",
		SystemPrompt: `你是Manus，一个智能AI助手。你具备以下能力：
1. 理解和分析用户的问题
2. 使用可用的工具来获取信息或执行任务
3. 基于ReAct框架进行推理：思考-行动-观察
4. 提供准确、有用的回答

请始终保持友好、专业的态度，并尽力帮助用户解决问题。`,
		NextPrompt:  "请继续分析并采取下一步行动。",
		MaxSteps:    10,
		MaxLoops:    5,
		EnableDebug: false,
	}
}

// Run 运行Manus处理用户输入
func (m *Manus) Run(ctx context.Context, input string) (*schema.Message, error) {
	m.logger.Info("Manus开始处理用户请求",
		zap.String("input", input),
		zap.String("name", m.config.Name))

	result, err := m.ToolCallAgent.Run(ctx, input)
	if err != nil {
		m.logger.Error("Manus处理失败", zap.Error(err))
		return nil, fmt.Errorf("Manus处理失败: %w", err)
	}

	m.logger.Info("Manus处理完成",
		zap.String("result", result.Content))

	return result, nil
}

// RunStream 流式运行Manus
func (m *Manus) RunStream(ctx context.Context, input string) (<-chan *schema.Message, error) {
	m.logger.Info("Manus开始流式处理用户请求",
		zap.String("input", input),
		zap.String("name", m.config.Name))

	return m.ToolCallAgent.RunStream(ctx, input)
}

// AddTool 添加工具
func (m *Manus) AddTool(tool schema.ToolInfo) {
	m.ToolCallAgent.AddTool(tool)
	m.logger.Info("添加工具", zap.String("tool", tool.Name))
}

// RemoveTool 移除工具
func (m *Manus) RemoveTool(toolName string) {
	m.ToolCallAgent.RemoveTool(toolName)
	m.logger.Info("移除工具", zap.String("tool", toolName))
}

// GetTools 获取所有工具
func (m *Manus) GetTools() []schema.ToolInfo {
	return m.ToolCallAgent.GetTools()
}

// Reset 重置状态
func (m *Manus) Reset() {
	m.ToolCallAgent.Reset()
	m.logger.Info("Manus状态已重置")
}

// GetConfig 获取配置
func (m *Manus) GetConfig() *ManusConfig {
	return m.config
}

// UpdateConfig 更新配置
func (m *Manus) UpdateConfig(config *ManusConfig) {
	if config != nil {
		m.config = config
		m.ToolCallAgent.ReActAgent.BaseAgent.SetMaxSteps(config.MaxSteps)
		m.ToolCallAgent.ReActAgent.SetMaxLoops(config.MaxLoops)
		m.logger.Info("Manus配置已更新")
	}
}

// GetState 获取当前状态
func (m *Manus) GetState() string {
	state := m.ToolCallAgent.ReActAgent.BaseAgent.GetState()
	return string(state)
}

// GetStepHistory 获取步骤历史
func (m *Manus) GetStepHistory() []string {
	return m.ToolCallAgent.ReActAgent.BaseAgent.GetStepHistory()
}

// SetDebugMode 设置调试模式
func (m *Manus) SetDebugMode(enable bool) {
	m.config.EnableDebug = enable
	if enable {
		m.logger.Info("Manus调试模式已启用")
	} else {
		m.logger.Info("Manus调试模式已禁用")
	}
}

// GetDebugInfo 获取调试信息
func (m *Manus) GetDebugInfo() map[string]interface{} {
	if !m.config.EnableDebug {
		return nil
	}

	return map[string]interface{}{
		"name":         m.config.Name,
		"state":        m.GetState(),
		"tools_count":  len(m.GetTools()),
		"step_history": m.GetStepHistory(),
		"current_loop": m.ToolCallAgent.ReActAgent.GetCurrentLoop(),
		"max_steps":    m.config.MaxSteps,
		"max_loops":    m.config.MaxLoops,
	}
}

// Think 执行思考步骤（暴露给外部使用）
func (m *Manus) Think(ctx context.Context) (*schema.Message, error) {
	return m.ToolCallAgent.ReActAgent.Think(ctx)
}

// Act 执行行动步骤（暴露给外部使用）
func (m *Manus) Act(ctx context.Context, thought string) (*schema.Message, error) {
	return m.ToolCallAgent.ReActAgent.Act(ctx, thought)
}

// Observe 执行观察步骤（暴露给外部使用）
func (m *Manus) Observe(ctx context.Context, action string) (*schema.Message, error) {
	return m.ToolCallAgent.ReActAgent.Observe(ctx, action)
}

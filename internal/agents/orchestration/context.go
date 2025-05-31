package orchestration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
)

// MemoryState 接口定义对话的短期记忆状态管理
type MemoryState interface {
	AddMessage(role, content string)
	GetSummary() string
	GetRecentMessages(n int) []schema.Message
	Clear()
	GetMessageCount() int
}

// OrchestrationContext 统一的编排上下文
type OrchestrationContext struct {
	ctx      context.Context        // 用于信号管理
	Input    map[string]interface{} // 用户输入的原始数据
	Memory   MemoryState            // 对话的短期记忆状态
	Metadata map[string]string      // 元数据，如模型版本等
	mu       sync.RWMutex           // 读写锁保护并发访问
}

// SimpleMemoryState 简单的内存状态实现
type SimpleMemoryState struct {
	messages    []schema.Message
	maxMessages int
	mu          sync.RWMutex
	createdAt   time.Time
}

// NewSimpleMemoryState 创建新的简单内存状态
func NewSimpleMemoryState(maxMessages int) *SimpleMemoryState {
	if maxMessages <= 0 {
		maxMessages = 100 // 默认保留100条消息
	}
	return &SimpleMemoryState{
		messages:    make([]schema.Message, 0),
		maxMessages: maxMessages,
		createdAt:   time.Now(),
	}
}

// AddMessage 添加消息到内存
func (s *SimpleMemoryState) AddMessage(role, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message := schema.Message{
		Role:    schema.RoleType(role),
		Content: content,
	}

	s.messages = append(s.messages, message)

	// 如果超过最大消息数，移除最早的消息
	if len(s.messages) > s.maxMessages {
		s.messages = s.messages[1:]
	}
}

// GetSummary 获取对话摘要
func (s *SimpleMemoryState) GetSummary() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.messages) == 0 {
		return "No conversation history"
	}

	// 简单实现：返回消息数量和时间跨度
	duration := time.Since(s.createdAt)
	return fmt.Sprintf("Conversation with %d messages over %v", len(s.messages), duration.Round(time.Second))
}

// GetRecentMessages 获取最近的n条消息
func (s *SimpleMemoryState) GetRecentMessages(n int) []schema.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if n <= 0 || len(s.messages) == 0 {
		return []schema.Message{}
	}

	if n >= len(s.messages) {
		// 返回所有消息的副本
		result := make([]schema.Message, len(s.messages))
		copy(result, s.messages)
		return result
	}

	// 返回最后n条消息的副本
	start := len(s.messages) - n
	result := make([]schema.Message, n)
	copy(result, s.messages[start:])
	return result
}

// Clear 清空内存状态
func (s *SimpleMemoryState) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = make([]schema.Message, 0)
	s.createdAt = time.Now()
}

// GetMessageCount 获取消息数量
func (s *SimpleMemoryState) GetMessageCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.messages)
}

// NewOrchestrationContext 创建新的编排上下文
func NewOrchestrationContext(ctx context.Context) *OrchestrationContext {
	return &OrchestrationContext{
		ctx:      ctx,
		Input:    make(map[string]interface{}),
		Memory:   NewSimpleMemoryState(100),
		Metadata: make(map[string]string),
	}
}

// NewOrchestrationContextWithMemory 使用指定内存状态创建编排上下文
func NewOrchestrationContextWithMemory(ctx context.Context, memory MemoryState) *OrchestrationContext {
	return &OrchestrationContext{
		ctx:      ctx,
		Input:    make(map[string]interface{}),
		Memory:   memory,
		Metadata: make(map[string]string),
	}
}

// Context 获取底层的context.Context
func (oc *OrchestrationContext) Context() context.Context {
	return oc.ctx
}

// WithContext 返回使用新context的副本
func (oc *OrchestrationContext) WithContext(ctx context.Context) *OrchestrationContext {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	return &OrchestrationContext{
		ctx:      ctx,
		Input:    oc.copyInput(),
		Memory:   oc.Memory,
		Metadata: oc.copyMetadata(),
	}
}

// SetInput 设置输入数据
func (oc *OrchestrationContext) SetInput(key string, value interface{}) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	oc.Input[key] = value
}

// GetInput 获取输入数据
func (oc *OrchestrationContext) GetInput(key string) (interface{}, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()
	value, exists := oc.Input[key]
	return value, exists
}

// GetInputString 获取字符串类型的输入数据
func (oc *OrchestrationContext) GetInputString(key string) (string, bool) {
	if value, exists := oc.GetInput(key); exists {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// SetMetadata 设置元数据
func (oc *OrchestrationContext) SetMetadata(key, value string) {
	oc.mu.Lock()
	defer oc.mu.Unlock()
	oc.Metadata[key] = value
}

// GetMetadata 获取元数据
func (oc *OrchestrationContext) GetMetadata(key string) (string, bool) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()
	value, exists := oc.Metadata[key]
	return value, exists
}

// copyInput 复制输入数据（内部使用，需要持有锁）
func (oc *OrchestrationContext) copyInput() map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range oc.Input {
		copy[k] = v
	}
	return copy
}

// copyMetadata 复制元数据（内部使用，需要持有锁）
func (oc *OrchestrationContext) copyMetadata() map[string]string {
	copy := make(map[string]string)
	for k, v := range oc.Metadata {
		copy[k] = v
	}
	return copy
}

// AddUserMessage 添加用户消息到内存
func (oc *OrchestrationContext) AddUserMessage(content string) {
	oc.Memory.AddMessage("user", content)
}

// AddAssistantMessage 添加助手消息到内存
func (oc *OrchestrationContext) AddAssistantMessage(content string) {
	oc.Memory.AddMessage("assistant", content)
}

// AddSystemMessage 添加系统消息到内存
func (oc *OrchestrationContext) AddSystemMessage(content string) {
	oc.Memory.AddMessage("system", content)
}

// GetConversationHistory 获取对话历史
func (oc *OrchestrationContext) GetConversationHistory(n int) []schema.Message {
	return oc.Memory.GetRecentMessages(n)
}

// ClearMemory 清空内存
func (oc *OrchestrationContext) ClearMemory() {
	oc.Memory.Clear()
}

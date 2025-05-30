package pipeline

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type ChatTemplateConfig struct {
	FormatType schema.FormatType
	Templates  []schema.MessagesTemplate
}

// newChatTemplate component initialization function of node 'ChatTemplate2' in graph 'Assitant'
func newChatTemplate(ctx context.Context) (ctp prompt.ChatTemplate, err error) {
	config := &ChatTemplateConfig{
		FormatType: schema.FString,
		Templates: []schema.MessagesTemplate{
			schema.SystemMessage("你是一个活泼的小助手，会用活泼的方式回答问题"),
			schema.SystemMessage(`你是一个专业的规划代理，负责通过结构化计划高效解决问题。
									你的职责是：
									1. 分析请求以理解任务范围。
									2. 创建清晰、可操作的计划。
									3. 根据需要使用可用工具执行步骤。
									4. 跟踪进度并在必要时调整计划。
									5. 回答的时候会详细介绍每一步及使用的工具`),
			schema.SystemMessage("根据用户回答检索到的内容为{retrieve_result}"),
			schema.UserMessage(ctx.Value("user_input").(string)),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}

package tests

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
			schema.SystemMessage("你是一个热情的小助手，会活泼地回答问题"),
			schema.SystemMessage("根据用户回答检索到的内容为{retrieve_result}"),
			schema.UserMessage(ctx.Value("user_input").(string)),
		},
	}
	ctp = prompt.FromMessages(config.FormatType, config.Templates...)
	return ctp, nil
}

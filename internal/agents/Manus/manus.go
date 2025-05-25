package manus

import (
	toolcallagent "MoonAgent/internal/agents/ToolCallAgent"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type Manus struct {
	ReActAgent *toolcallagent.ToolCallAgent
}

func NewManus(name string, systemPrompt string, nextPrompt string, chatModel *model.BaseChatModel, tools []schema.ToolInfo) *Manus {
	manus := &Manus{
		ReActAgent: toolcallagent.NewToolCallAgent(name, systemPrompt, nextPrompt, chatModel, tools),
	}
	return manus
}

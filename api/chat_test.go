package api

import (
	"context"
	"testing"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
)

func TestSendChat(t *testing.T) {
	chatModel, _ := ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "deepseek-r1:7b",         // 模型名称
	})
	messages := []*schema.Message{
		schema.SystemMessage("你是一个助手"),
		schema.UserMessage("你好"),
	}
	resp, err := chatModel.Generate(context.Background(), messages)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)
	t.Logf("resp: %v", resp)
}

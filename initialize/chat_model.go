package initialize

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"go.uber.org/zap"
)

func InitChatModel() {
	chatModel, err := ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: global.ServerConfig.LLMConfig.BASE_URL, // Ollama 服务地址
		Model:   global.ServerConfig.LLMConfig.MODEL,    // 模型名称
	})
	if err != nil {
		zap.S().Errorw("[InitChatModel]", "err", err.Error())
		return
	}
	global.ChatModel = chatModel
}

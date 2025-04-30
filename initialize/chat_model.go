package initialize

import (
	"SuperAgent/global"
	"SuperAgent/tools"
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

	// 绑定工具
	chatModel.BindTools(tools.BodyReportTool())

	global.ChatModel = chatModel
}

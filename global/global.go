package global

import (
	"SuperAgent/config"
	"github.com/cloudwego/eino-ext/components/model/ollama"
)

var (
	ServerConfig *config.ServerConfig
	ChatModel *ollama.ChatModel
)

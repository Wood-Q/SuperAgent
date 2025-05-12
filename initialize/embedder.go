package initialize

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"go.uber.org/zap"
)

func InitEmbedder() {
	embedder, err := ark.NewEmbedder(context.Background(), &ark.EmbeddingConfig{
		APIKey: global.ServerConfig.DocumentConfig.ArkApiKey,
		Model:  global.ServerConfig.DocumentConfig.ArkModel,
	})
	if err != nil {
		zap.S().Error("Failed to create embedder: %v", zap.String("error", err.Error()))
	}
	zap.S().Info("Embedder created: %v", zap.String("model", global.ServerConfig.DocumentConfig.ArkModel))
	global.Embedder = embedder
}

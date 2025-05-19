package embedder

import (
	"MoonAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"go.uber.org/zap"
)

func ProvideEmbedder(cfg *config.ServerConfig) (*ark.Embedder, error) {
	embedder, err := ark.NewEmbedder(context.Background(), &ark.EmbeddingConfig{
		APIKey: cfg.DocumentConfig.API_KEY,
		Model:  cfg.DocumentConfig.Model,
	})
	if err != nil {
		zap.S().Error("Failed to create embedder: %v", zap.String("error", err.Error()))
		return nil, err
	}
	zap.S().Info("Embedder created: %v", zap.String("model", cfg.DocumentConfig.Model))
	return embedder, nil
}

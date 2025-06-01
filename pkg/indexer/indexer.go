package indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"go.uber.org/zap"
)

// InitIndexer 初始化Indexer

func NewIndexer(ctx context.Context, es8Config *es8.IndexerConfig) (*es8.Indexer, error) {
	indexer, err := es8.NewIndexer(ctx, es8Config)
	if err != nil {
		zap.S().Panic("Failed to create indexer: %v", zap.String("error", err.Error()))
	}
	zap.S().Info("Indexer created successfully")
	return indexer, nil
}

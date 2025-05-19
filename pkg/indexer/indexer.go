package indexer

import (
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"go.uber.org/zap"
)

// InitIndexer 初始化Indexer

func NewIndexer(ctx context.Context, milvusConfig *milvus.IndexerConfig) (*milvus.Indexer, error) {
	indexer, err := milvus.NewIndexer(ctx, milvusConfig)
	if err != nil {
		zap.S().Panic("Failed to create indexer: %v", zap.String("error", err.Error()))
	}
	zap.S().Info("Indexer created successfully")
	return indexer, nil
}

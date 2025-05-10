package initialize

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"go.uber.org/zap"
)

func InitIndexer() {
	indexer, err := milvus.NewIndexer(context.Background(), &milvus.IndexerConfig{
		Client:     *global.MilvusClient,
		Embedding:  global.Embedder,
	})
	if err != nil {
		zap.S().Error("Failed to create indexer: %v", zap.String("error", err.Error()))
	}
	zap.S().Info("Indexer created successfully")
	global.Indexer = indexer
}

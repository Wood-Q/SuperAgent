package indexer

import (
	"MoonAgent/pkg/config"
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
)

func NewIndexer(ctx context.Context, cli *client.Client, emb *ark.Embedder, docs []*schema.Document, cfg *config.ServerConfig) (*milvus.Indexer, error) {
	// Create an embedding model
	emb, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey: cfg.DocumentConfig.API_KEY,
		Model:  cfg.DocumentConfig.Model,
	})
	if err != nil {
		zap.S().Error("Failed to create embedding: %v", zap.String("error", err.Error()))
		return nil, err
	}

	// Create an indexer
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:    *cli,
		Embedding: emb,
	})
	if err != nil {
		zap.S().Error("Failed to create indexer: %v", zap.String("error", err.Error()))
		return nil, err
	}
	zap.S().Info("Indexer created success")
	// Store documents
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		zap.S().Error("Failed to store: %v", zap.String("error", err.Error()))
		return nil, err
	}
	zap.S().Info("Store success, ids: %v", ids)
	return indexer, nil
}

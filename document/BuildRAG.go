package document

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func BuildRAG(id, content string, metadata map[string]any, ctx context.Context) error {
	// Store documents
	docs := []*schema.Document{
		{
			ID:       id,
			Content:  content,
			MetaData: metadata,
		},
	}
	ids, err := global.Indexer.Store(ctx, docs)
	if err != nil {
		zap.S().Error("Failed to store: %v", zap.String("error", err.Error()))
		return err
	}
	zap.S().Info("Store success, ids: ", ids)
	return nil
}

package splitter

import (
	"context"
	"strconv"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func SplitDocs(ctx context.Context, embedder *ark.Embedder, docs []*schema.Document) ([]*schema.Document, error) {
	splitter, err := semantic.NewSplitter(ctx, &semantic.Config{
		Embedding:    embedder,
		BufferSize:   2,
		MinChunkSize: 100,
		Percentile:   0.9,
	})
	if err != nil {
		zap.S().Error("Failed to create splitter: %v", zap.String("error", err.Error()))
		return nil, err
	}
	// 执行分割
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		zap.S().Error("Failed to transform: %v", zap.String("error", err.Error()))
		return nil, err
	}
	for i, result := range results {
		result.ID = docs[0].ID + "_" + strconv.Itoa(i)
	}
	return results, nil
}

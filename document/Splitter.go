package document

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func Split(ctx context.Context, docs []*schema.Document) ([]*schema.Document, error) {
	// 初始化分割器
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
		ChunkSize:   500,
		OverlapSize: 100,
		Separators:  []string{"\n", "。", "！", "？"},
		KeepType:    recursive.KeepTypeEnd,
	})
	if err != nil {
		zap.S().Error("new splitter failed", zap.Error(err))
	}

	// 执行分割
	results, err := splitter.Transform(ctx, docs)
	if err != nil {
		zap.S().Error("split failed", zap.Error(err))
	}

	// 处理分割结果
	for i, doc := range results {
		zap.S().Info("片段", zap.Int("片段", i+1), zap.String("内容", doc.Content))
	}

	return results, nil
}

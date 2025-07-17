package splitter

import (
	"context"
	"strconv"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"
)

func SplitDocs(ctx context.Context, embedder *ark.Embedder, docs []*schema.Document) ([]*schema.Document, error) {
	// splitter, err := semantic.NewSplitter(ctx, &semantic.Config{
	// 	Embedding:    embedder,
	// 	BufferSize:   2,
	// 	MinChunkSize: 100,
	// 	Percentile:   0.9,
	// })
	splitter, err := recursive.NewSplitter(ctx, &recursive.Config{
    ChunkSize:    400,           // 必需：目标片段大小
    OverlapSize:  10,            // 可选：片段重叠大小
    Separators:   []string{"\n", ".", "?", "!"}, // 可选：分隔符列表
    LenFunc:      nil,            // 可选：自定义长度计算函数
    KeepType:     recursive.KeepTypeNone, // 可选：分隔符保留策略
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

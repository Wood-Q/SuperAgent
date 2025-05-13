package initialize

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"go.uber.org/zap"
)

// InitIndexer 初始化Indexer
func InitIndexer() {
	// 检查collection是否存在
	client := *global.MilvusClient
	exists, err := client.HasCollection(context.Background(), "agent")
	if err != nil {
		zap.S().Error("检查collection失败: %v", zap.String("error", err.Error()))
	}

	if !exists {
		zap.S().Warn("Collection 'agent' 不存在，请先创建")
	}

	fields := []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "255",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector", // 确保字段名匹配
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "81920",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}

	indexer, err := milvus.NewIndexer(context.Background(), &milvus.IndexerConfig{
		Client:     client,
		Embedding:  global.Embedder,
		Collection: "try",
		Fields:     fields,
	})
	if err != nil {
		zap.S().Panic("Failed to create indexer: %v", zap.String("error", err.Error()))
	}
	zap.S().Info("Indexer created successfully")
	global.Indexer = indexer
}

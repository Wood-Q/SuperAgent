package initialize

import (
	"SuperAgent/global"
	"context"

	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"go.uber.org/zap"
)

func InitRetriever() {
	// Create a retriever
	retriever, err := milvus.NewRetriever(context.Background(), &milvus.RetrieverConfig{
		Client:      *global.MilvusClient,
		Collection:  "",
		Partition:   nil,
		VectorField: "",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      1,
		Embedding: global.Embedder,
	})
	if err != nil {
		zap.S().Error("Failed to create retriever: %v", zap.String("error", err.Error()))
	}
	global.Retriever = retriever
}

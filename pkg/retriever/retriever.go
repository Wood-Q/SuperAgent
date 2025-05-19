package retriever

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

// ProvideRetriever 提供检索器
func ProvideRetriever(milvusCli *client.Client, embedder *ark.Embedder) (*milvus.Retriever, error) {
	retriever, err := milvus.NewRetriever(context.Background(), &milvus.RetrieverConfig{
		Client:      *milvusCli,
		Collection:  "try",
		Partition:   nil,
		VectorField: "",
		OutputFields: []string{
			"id",
			"content",
			"metadata",
		},
		TopK:      1,
		Embedding: embedder,
	})
	if err != nil {
		zap.S().Error("Failed to create retriever: %v", zap.String("error", err.Error()))
		return nil, err
	}
	return retriever, nil
}

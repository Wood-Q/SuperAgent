package di

import (
	"context"

	"MoonAgent/pkg/config"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/google/wire"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

// ProvideMilvusClient 提供 Milvus 客户端
func ProvideMilvusClient() (*client.Client, error) {
	Client, err := client.NewClient(context.Background(), client.Config{
		Address: config.GlobalConfig.DocumentConfig.Addr,
	})
	if err != nil {
		zap.S().Error("Failed to create client: %v", zap.String("error", err.Error()))
		return nil, err
	}
	zap.S().Info("Client created successfully")
	return &Client, nil
}

// ProvideEmbedder 提供嵌入器
func ProvideEmbedder() (*ark.Embedder, error) {
	embedder, err := ark.NewEmbedder(context.Background(), &ark.EmbeddingConfig{
		APIKey: config.GlobalConfig.DocumentConfig.API_KEY,
		Model:  config.GlobalConfig.DocumentConfig.Model,
	})
	if err != nil {
		zap.S().Error("Failed to create embedder: %v", zap.String("error", err.Error()))
		return nil, err
	}
	zap.S().Info("Embedder created: %v", zap.String("model", config.GlobalConfig.DocumentConfig.Model))
	return embedder, nil
}

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

// Application 应用程序结构体
type Application struct {
	MilvusClient *client.Client
	Embedder     *ark.Embedder
	Retriever    *milvus.Retriever
}

// ProvideApplication 提供应用程序实例
func ProvideApplication(
	milvusClient *client.Client,
	embedder *ark.Embedder,
	retriever *milvus.Retriever,
) *Application {
	return &Application{
		MilvusClient: milvusClient,
		Embedder:     embedder,
		Retriever:    retriever,
	}
}

// ProviderSet 依赖注入提供者集合
var ProviderSet = wire.NewSet(
	ProvideMilvusClient,
	ProvideEmbedder,
	ProvideRetriever,
	ProvideApplication,
)

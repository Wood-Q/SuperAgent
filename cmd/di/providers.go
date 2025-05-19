package di

import (
	"MoonAgent/pkg/embedder"
	userClient "MoonAgent/pkg/milvus"

	retriever "MoonAgent/pkg/retriever"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/google/wire"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

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
	userClient.ProvideMilvusClient,
	embedder.ProvideEmbedder,
	retriever.ProvideRetriever,
	ProvideApplication,
)

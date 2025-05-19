package di

import (
	"MoonAgent/pkg/config"
	"MoonAgent/pkg/embedder"
	userClient "MoonAgent/pkg/milvus"
	"context"

	retriever "MoonAgent/pkg/retriever"

	indexer "github.com/cloudwego/eino-ext/components/indexer/milvus"

	fields "MoonAgent/pkg/indexer"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/google/wire"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

// Application 应用程序结构体
type Application struct {
	ServerConfig  *config.ServerConfig
	MilvusClient  *client.Client
	Embedder      *ark.Embedder
	IndexerConfig *indexer.IndexerConfig
	Indexer       *indexer.Indexer
	Retriever     *milvus.Retriever
}

// ProvideContext 提供上下文
func ProvideContext() context.Context {
	return context.Background()
}

// ProvideApplication 提供应用程序实例
func ProvideApplication(
	serverConfig *config.ServerConfig,
	milvusClient *client.Client,
	embedder *ark.Embedder,
	retriever *milvus.Retriever,
	indexerConfig *indexer.IndexerConfig,
	indexer *indexer.Indexer,
) *Application {
	return &Application{
		ServerConfig:  serverConfig,
		MilvusClient:  milvusClient,
		Embedder:      embedder,
		Retriever:     retriever,
		IndexerConfig: indexerConfig,
		Indexer:       indexer,
	}
}

// ProviderSet 依赖注入提供者集合
var ProviderSet = wire.NewSet(
	// 1. 首先提供基础依赖
	ProvideContext,
	config.NewConfig,

	// 2. 提供中间依赖
	userClient.ProvideMilvusClient,
	embedder.ProvideEmbedder,

	// 3. 提供配置
	fields.NewIndexerConfig,

	// 4. 提供主要组件
	indexer.NewIndexer,
	retriever.ProvideRetriever,

	// 5. 最后提供应用实例
	ProvideApplication,
)

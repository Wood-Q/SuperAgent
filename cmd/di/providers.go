package di

import (
	"MoonAgent/pkg/config"
	"MoonAgent/pkg/embedder"
	userClient "MoonAgent/pkg/es8"
	"context"

	retriever "MoonAgent/pkg/retriever"

	indexer "github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/elastic/go-elasticsearch/v8"

	fields "MoonAgent/pkg/indexer"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/google/wire"
)

// Application 应用程序结构体
type Application struct {
	ServerConfig  *config.ServerConfig
	Es8Client     *elasticsearch.Client
	Embedder      *ark.Embedder
	IndexerConfig *indexer.IndexerConfig
	Indexer       *indexer.Indexer
	Retriever     *es8.Retriever
}

// ProvideContext 提供上下文
func ProvideContext() context.Context {
	return context.Background()
}

// ProvideApplication 提供应用程序实例
func ProvideApplication(
	serverConfig *config.ServerConfig,
	es8Client *elasticsearch.Client,
	embedder *ark.Embedder,
	retriever *es8.Retriever,
	indexerConfig *indexer.IndexerConfig,
	indexer *indexer.Indexer,
) *Application {
	return &Application{
		ServerConfig:  serverConfig,
		Es8Client:     es8Client,
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
	userClient.ProvideEs8Client,
	embedder.ProvideEmbedder,

	// 3. 提供配置
	fields.NewIndexerConfig,

	// 4. 提供主要组件
	indexer.NewIndexer,
	retriever.ProvideRetriever,

	// 5. 最后提供应用实例
	ProvideApplication,
)

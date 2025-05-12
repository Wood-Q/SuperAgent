package global

import (
	"SuperAgent/config"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	retriever "github.com/cloudwego/eino-ext/components/retriever/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
)

var (
	ServerConfig *config.ServerConfig
	ChatModel    *ollama.ChatModel
	MilvusClient *client.Client
	Embedder     *ark.Embedder
	Indexer      *milvus.Indexer
	Retriever    *retriever.Retriever
)

package initialize

import (
	"SuperAgent/global"
	"context"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"go.uber.org/zap"
)

func InitClient() {
	Client, err := client.NewClient(context.Background(), client.Config{
		Address: global.ServerConfig.DocumentConfig.Addr,
	})
	if err != nil {
		zap.S().Error("Failed to create client: %v", zap.String("error", err.Error()))
	}
	zap.S().Info("Client created successfully")
	global.MilvusClient = &Client
}

package milvus

import (
	"MoonAgent/pkg/config"
	"context"

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

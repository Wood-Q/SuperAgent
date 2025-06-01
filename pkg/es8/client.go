package es8

import (
	"MoonAgent/pkg/config"

	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

// ProvideEs8Client 提供 Es8 客户端
func ProvideEs8Client(cfg *config.ServerConfig) (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.DocumentConfig.Addr},
		Username:  cfg.DocumentConfig.UserName,
		Password:  cfg.DocumentConfig.Password,
	})
	if err != nil {
		zap.S().Errorw("Failed to create es8 client", "error", err.Error())
		return nil, err
	}
	zap.S().Info("Es8 client created successfully")
	return client, nil
}

package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 全局变量，外部只读
var GlobalConfig *ServerConfig

// NewConfig 加载完毕后直接把结构体返回
func NewConfig() (*ServerConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/Users/woodq/FullStack/SuperAgent/configs")
	if err := v.ReadInConfig(); err != nil {
		zap.S().Errorw("Failed to read config file", "error", err.Error())
		return nil, err
	}

	// 直接分配一个实例
	tmp := new(ServerConfig)
	if err := v.Unmarshal(tmp); err != nil {
		zap.S().Errorw("Failed to unmarshal config", "error", err.Error())
		return nil, err
	}
	return tmp, nil
}

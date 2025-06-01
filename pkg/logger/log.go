package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger 初始化logger
func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Logger = logger
	zap.ReplaceGlobals(logger)
}

// InitDevelopmentLogger 初始化开发模式的logger，输出更友好的日志格式
func InitDevelopmentLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
	zap.ReplaceGlobals(logger)
}

// InitTestLogger 初始化测试环境的logger，显示详细信息
func InitTestLogger() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.DisableStacktrace = true

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	Logger = logger
	zap.ReplaceGlobals(logger)
}

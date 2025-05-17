package logger

import (
	"MoonAgent/internal/global"

	"go.uber.org/zap"
)

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	global.Logger = logger
}

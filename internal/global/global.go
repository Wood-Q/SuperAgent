package global

import (
	"MoonAgent/internal/models"

	"go.uber.org/zap"
)

var (
	Logger       *zap.Logger
	ServerConfig *models.ServerConfig
)

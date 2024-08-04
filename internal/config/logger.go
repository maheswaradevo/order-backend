package config

import (
	"order-service-backend/internal/common/constants"

	"go.uber.org/zap"
)

func NewLogger(config *Config) *zap.Logger {
	var logger *zap.Logger
	if config.AppEnvironment == constants.LocalEnv {
		logger, _ = zap.NewDevelopment()
	} else if config.AppEnvironment == constants.ProductionEnv {
		logger, _ = zap.NewProduction()
	}
	return logger
}

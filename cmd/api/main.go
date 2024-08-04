package main

import (
	"fmt"
	"order-service-backend/internal/common/constants"
	"order-service-backend/internal/config"

	"go.uber.org/zap"
)

func main() {
	config.Init()

	cfg := config.GetConfig()

	logger := config.NewLogger(cfg)
	db := config.NewDatabase(cfg, logger)
	app := config.NewEcho(cfg)
	rabbitMqClient, err := config.NewRabbitMQ(cfg.RabbitMqConfig)
	if err != nil {
		logger.Fatal("failed to start rabbitmq client: ", zap.Error(err))
	}

	config.Bootstrap(&config.BootstrapConfig{
		DB:           db,
		Log:          logger,
		App:          app,
		Config:       cfg,
		RabbitMQConn: rabbitMqClient.Conn,
		RabbitMQChan: rabbitMqClient.Chann,
		RabbitMQQuit: rabbitMqClient.QuitChann,
		Events:       config.NewEvent(cfg),
	})

	var address string
	if cfg.AppEnvironment == constants.LocalEnv {
		address = fmt.Sprintf("%s:%s", "localhost", cfg.Port)
	} else if cfg.AppEnvironment == constants.ProductionEnv {
		address = fmt.Sprintf("%s:%s", cfg.ProductionEnvironment, cfg.Port)
	}

	if err := app.Start(address); err != nil {
		logger.Fatal("failed to start server: ", zap.Error(err))
	}
}
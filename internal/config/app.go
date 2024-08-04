package config

import (
	"order-service-backend/internal/delivery/http/route"
	"order-service-backend/internal/models"

	"github.com/labstack/echo/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB           *gorm.DB
	App          *echo.Echo
	Log          *zap.Logger
	Config       *Config
	Events       models.Events
	RabbitMQConn *amqp.Connection
	RabbitMQChan *amqp.Channel
	RabbitMQQuit chan bool
}

func Bootstrap(config *BootstrapConfig) {
	// Setup Repositories

	// Setup PubSub

	// Setup usecases

	// Setup Controller

	routeConfig := route.RouteConfig{
		App: config.App,
	}

	routeConfig.Setup()
}

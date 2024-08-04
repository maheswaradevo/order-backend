package route

import (
	"time"

	"github.com/maheswaradevo/order-backend/internal/config"
	"github.com/maheswaradevo/order-backend/internal/delivery/http"
	events "github.com/maheswaradevo/order-backend/internal/delivery/messaging"
	"github.com/maheswaradevo/order-backend/internal/gateway/messaging"
	"github.com/maheswaradevo/order-backend/internal/models"
	"github.com/maheswaradevo/order-backend/internal/models/consumer"
	"github.com/maheswaradevo/order-backend/internal/repository"
	"github.com/maheswaradevo/order-backend/internal/usecase"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RouteConfig struct {
	App             *echo.Echo
	OrderController *http.OrderController
}
type BootstrapConfig struct {
	DB             *gorm.DB
	App            *echo.Echo
	Log            *zap.Logger
	Config         *config.Config
	Events         models.Events
	RabbitMQConn   *amqp.Connection
	RabbitMQChan   *amqp.Channel
	RabbitMQQuit   chan bool
	AuthMiddleware *echo.MiddlewareFunc
	Redis          *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	// repository
	orderRepository := repository.NewOrderRepository(config.Log)

	creditLimitCh := make(chan []consumer.CreditLimitEvent)

	// pub
	orderMessaging := messaging.NewOrderPublisher(&config.Events, config.Log)

	// use case
	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, orderRepository, orderMessaging, creditLimitCh)

	// sub
	orderConsumer := events.NewOrderConsumer(orderUseCase, creditLimitCh)

	//controller
	orderController := http.NewOrderController(config.Log, orderUseCase)

	// config.App.Use(config.AuthMiddleware)

	routeConfig := RouteConfig{
		App:             config.App,
		OrderController: orderController,
	}

	time.AfterFunc(2*time.Second, func() {
		orderConsumer.ConsumeCreditLimitData(config.Log, &config.Events)
	})

	routeConfig.Setup()
}

func (c *RouteConfig) Setup() {
	//Setup Route
	c.SetupOrderRoute()

}

func (c *RouteConfig) SetupOrderRoute() {
	// Setup endpoint
	c.App.POST("/api/orders", c.OrderController.CreateOrder)
}

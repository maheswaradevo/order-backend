package route

import (
	"order-service-backend/internal/config"
	"order-service-backend/internal/delivery/http"

	events "order-service-backend/internal/delivery/messaging"
	"order-service-backend/internal/models"
	"order-service-backend/internal/repository"
	"order-service-backend/internal/usecase"
	"time"

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
}

func Bootstrap(config *BootstrapConfig) {
	orderRepository := repository.NewOrderRepository(config.Log)

	orderUseCase := usecase.NewOrderUseCase(config.DB, config.Log, orderRepository)

	orderController := http.NewOrderController(config.Log, orderUseCase)

	config.App.Use(*config.AuthMiddleware)

	routeConfig := RouteConfig{
		App:             config.App,
		OrderController: orderController,
	}

	time.AfterFunc(5*time.Second, func() {
		events.ConsumeUserData(config.Log, &config.Events)
	})

	routeConfig.Setup()
}

func (c *RouteConfig) Setup() {
	//Setup Route

}

func (c *RouteConfig) SetupOrderRoute() {
	// Setup endpoint
}

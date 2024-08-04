package usecase

import (
	"order-service-backend/internal/entity"
	"order-service-backend/internal/models"
	"order-service-backend/internal/repository"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB              *gorm.DB
	Log             *zap.Logger
	OrderRepository *repository.OrderRepository
}

func NewOrderUseCase(db *gorm.DB, logger *zap.Logger, orderRepository *repository.OrderRepository) *OrderUseCase {
	return &OrderUseCase{
		DB:              db,
		Log:             logger,
		OrderRepository: orderRepository,
	}
}

func (u *OrderUseCase) CreateOrder(ctx echo.Context, request *models.OrderCreateRequest) (*entity.Order, error) {
	return nil, nil
}

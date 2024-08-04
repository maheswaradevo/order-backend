package http

import (
	"order-service-backend/internal/usecase"

	"go.uber.org/zap"
)

type OrderController struct {
	Log     *zap.Logger
	UseCase *usecase.OrderUseCase
}

func NewOrderController(logger *zap.Logger, orderUseCase *usecase.OrderUseCase) *OrderController {
	return &OrderController{
		Log:     logger,
		UseCase: orderUseCase,
	}
}

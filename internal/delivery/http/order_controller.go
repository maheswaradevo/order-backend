package http

import (
	"net/http"
	commons "order-service-backend/internal/common"
	"order-service-backend/internal/common/jwt"
	"order-service-backend/internal/models"
	"order-service-backend/internal/usecase"
	"strings"

	"github.com/labstack/echo/v4"
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

func (c *OrderController) CreateOrder(ctx echo.Context) error {

	headerValue := ctx.Request().Header.Get("Authorization")
	token := strings.Replace(headerValue, "Bearer ", "", -1)

	res, _ := jwt.Decode(token)
	customerId, ok := res["id"].(float64)
	if !ok {
		return ctx.JSON(http.StatusBadRequest, commons.ResponseFailed("failed type assertion"))
	}
	customerIdUint := uint64(customerId)
	pyld := models.OrderCreateRequest{}

	if err := ctx.Bind(&pyld); err != nil {
		return ctx.JSON(http.StatusBadRequest, commons.ResponseFailed(err.Error()))
	}

	pyld.CustomerID = customerIdUint

	order, err := c.UseCase.CreateOrder(ctx, &pyld)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, commons.ResponseFailed(err.Error()))
	}

	return ctx.JSON(http.StatusOK, commons.ResponseSuccess("success", order))
}

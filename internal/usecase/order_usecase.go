package usecase

import (
	"errors"
	"fmt"

	"strings"

	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/order-backend/internal/common/helpers"
	"github.com/maheswaradevo/order-backend/internal/entity"
	"github.com/maheswaradevo/order-backend/internal/gateway/messaging"
	"github.com/maheswaradevo/order-backend/internal/models"
	"github.com/maheswaradevo/order-backend/internal/models/consumer"
	"github.com/maheswaradevo/order-backend/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderUseCase struct {
	DB              *gorm.DB
	Log             *zap.Logger
	OrderMessaging  *messaging.OrderPublisher
	OrderRepository *repository.OrderRepository
	CreditLimitCh   chan []consumer.CreditLimitEvent
}

func NewOrderUseCase(db *gorm.DB, logger *zap.Logger, orderRepository *repository.OrderRepository, orderMessaging *messaging.OrderPublisher, ch chan []consumer.CreditLimitEvent) *OrderUseCase {
	return &OrderUseCase{
		DB:              db,
		Log:             logger,
		OrderRepository: orderRepository,
		OrderMessaging:  orderMessaging,
		CreditLimitCh:   ch,
	}
}

func (u *OrderUseCase) CreateOrder(ctx echo.Context, request *models.OrderCreateRequest) (*entity.Order, error) {
	tx := u.DB.WithContext(helpers.Context(ctx)).Begin()
	defer tx.Rollback()

	var result *entity.Order
	noContract, err := u.GenerateContractNumber(ctx)
	if err != nil {
		u.Log.Error("failed to create new order: ", zap.Error(err))
		return nil, err
	}

	var limitLeft float64
	go u.PushCreditLimitRequest(request.CustomerID)

	creditLimitData := <-u.CreditLimitCh
	for _, v := range creditLimitData {
		duration := v.EndDate.Sub(v.StartDate)
		month := int(duration.Hours() / (24 * 30))

		if month == request.Tenor {
			perMonthPayment := request.OTR / float64(request.Tenor)
			perMonthPayment = perMonthPayment + request.AdminFee + request.InterestValue

			if v.CreditLimit >= perMonthPayment {
				// Not Eligible for credit
				limitLeft = v.CreditLimit - perMonthPayment
				// Update limit
				go u.PushUpdateCreditLimit(models.CreditLimitUpdate{
					ID:          uint64(v.ID),
					CreditLimit: limitLeft,
				})
				// create order
				order := models.OrderCreateRequest{
					ContractNumber:   noContract,
					OTR:              request.OTR,
					InterestValue:    request.InterestValue,
					InstallmentValue: perMonthPayment,
					AssetName:        request.AssetName,
					Tenor:            request.Tenor,
					CustomerID:       request.CustomerID,
				}
				result, err = u.OrderRepository.Create(tx, &order)
				if err != nil {
					u.Log.Error("failed to create order: ", zap.Error(err))
					return nil, err
				}
			} else {
				return nil, errors.New("credit limit reached")
			}
		} else {
			continue
		}
	}
	tx.Commit()
	return result, nil
}

func (u *OrderUseCase) GenerateContractNumber(ctx echo.Context) (string, error) {
	tx := u.DB.WithContext(helpers.Context(ctx)).Begin()
	defer tx.Rollback()

	_, count, err := u.OrderRepository.GetAll(tx, &models.OrderFilter{})
	if err != nil {
		u.Log.Error("failed to generate contract number: ", zap.Error(err))
		return "", err
	}

	noContract := fmt.Sprintf("%05d", count+1)
	noContract = strings.Replace(noContract, " ", "0", -1)

	return noContract, nil
}

func (u *OrderUseCase) PushCreditLimitRequest(customerID uint64) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			u.Log.Error("recovered from panic ", zap.Any("error", err))
		}
	}()

	return u.OrderMessaging.PushCreditLimitRequest(&models.CreditLimitRequest{
		CustomerID: customerID,
	})
}

func (u *OrderUseCase) PushUpdateCreditLimit(data models.CreditLimitUpdate) (bool, error) {
	defer func() {
		if err := recover(); err != nil {
			u.Log.Error("recovered from panic ", zap.Any("error", err))
		}
	}()
	return u.OrderMessaging.PushUpdateCreditLimit(&models.CreditLimitUpdate{
		ID:          data.ID,
		CreditLimit: data.CreditLimit,
	})
}

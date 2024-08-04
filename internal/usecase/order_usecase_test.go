package usecase_test

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/maheswaradevo/order-backend/internal/config"
	"github.com/maheswaradevo/order-backend/internal/entity"
	"github.com/maheswaradevo/order-backend/internal/gateway/messaging"
	"github.com/maheswaradevo/order-backend/internal/models"
	"github.com/maheswaradevo/order-backend/internal/repository"
	"github.com/maheswaradevo/order-backend/internal/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type GormDBInterface interface {
	WithContext(ctx context.Context) *gorm.DB
	Begin() *gorm.DB
	Commit() error
	Rollback() error
}

type MockGormDB struct {
	mock.Mock
}

func (m *MockGormDB) WithContext(ctx context.Context) *gorm.DB {
	args := m.Called(ctx)
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Begin() *gorm.DB {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}

func (m *MockGormDB) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGormDB) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(db *gorm.DB, order interface{}) (*entity.Order, error) {
	args := m.Called(db, order)
	return args.Get(0).(*entity.Order), args.Error(1)
}

type MockOrderPublisher struct {
	mock.Mock
}

func (m *MockOrderPublisher) PushCreditLimitRequest(request interface{}) (bool, error) {
	args := m.Called(request)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderPublisher) PushUpdateCreditLimit(update interface{}) (bool, error) {
	args := m.Called(update)
	return args.Bool(0), args.Error(1)
}

type OrderUseCaseTestSuite struct {
	suite.Suite
	DB             *MockGormDB
	OrderRepo      *MockOrderRepository
	OrderMessaging *MockOrderPublisher
	Service        *usecase.OrderUseCase
}

func (s *OrderUseCaseTestSuite) SetupTest() {
	config.Init()
	s.DB = new(MockGormDB)
	s.OrderRepo = new(MockOrderRepository)
	s.OrderMessaging = new(MockOrderPublisher)
	logger := config.NewLogger(config.GetConfig())
	events := config.NewEvent(config.GetConfig())

	s.Service = &usecase.OrderUseCase{
		DB:              config.NewDatabase(config.GetConfig(), logger),
		OrderMessaging:  messaging.NewOrderPublisher(&events, logger),
		OrderRepository: repository.NewOrderRepository(logger),
	}
}

func (s *OrderUseCaseTestSuite) TestCreateOrder_Success() {
	mockCtx := echo.New().NewContext(nil, nil)

	expectedOrder := &entity.Order{
		CustomerID:     1,
		ContractNumber: "00001",
		OTR:            500.0,
		AdminFee:       10.0,
		InterestValue:  50.0,
		AssetName:      "Asset 1",
		Tenor:          1,
	}

	s.OrderRepo.On("Create", s.DB, mock.Anything).Return(expectedOrder, nil)

	s.OrderMessaging.On("PushCreditLimitRequest", mock.Anything).Return(true, nil)
	s.OrderMessaging.On("PushUpdateCreditLimit", mock.Anything).Return(true, nil)

	s.DB.On("WithContext", mock.Anything).Return(s.DB)
	s.DB.On("Begin").Return(s.DB)
	s.DB.On("Commit").Return(nil)
	s.DB.On("Rollback").Return(nil)

	request := &models.OrderCreateRequest{
		CustomerID:    1,
		OTR:           500.0,
		AdminFee:      10.0,
		InterestValue: 50.0,
		AssetName:     "Asset 1",
		Tenor:         1,
	}

	result, err := s.Service.CreateOrder(mockCtx, request)

	s.NoError(err)
	s.Equal(expectedOrder, result)
}

func (s *OrderUseCaseTestSuite) TestCreateOrder_CreditLimitReached() {
	mockCtx := echo.New().NewContext(nil, nil)

	mockOrder := &entity.Order{ID: 1}
	s.OrderRepo.On("Create", mock.Anything, mock.Anything).Return(mockOrder, nil)

	s.OrderMessaging.On("PushCreditLimitRequest", mock.Anything).Return(true, nil)

	request := &models.OrderCreateRequest{
		OTR:           1000.0,
		Tenor:         6,
		AdminFee:      10.0,
		InterestValue: 50.0,
		CustomerID:    1,
		AssetName:     "Test Asset",
	}

	result, err := s.Service.CreateOrder(mockCtx, request)

	s.Nil(result)
	s.EqualError(err, "credit limit reached")
}

func TestOrderUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderUseCaseTestSuite))
}

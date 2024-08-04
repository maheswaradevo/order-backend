package repository

import (
	"github.com/maheswaradevo/order-backend/internal/entity"
	"github.com/maheswaradevo/order-backend/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderRepository struct {
	logger *zap.Logger
}

func NewOrderRepository(logger *zap.Logger) *OrderRepository {
	return &OrderRepository{logger: logger}
}

func (r *OrderRepository) Create(db *gorm.DB, data *models.OrderCreateRequest) (*entity.Order, error) {
	result := entity.Order{
		CustomerID:       data.CustomerID,
		ContractNumber:   data.ContractNumber,
		OTR:              data.OTR,
		AdminFee:         data.AdminFee,
		InstallmentValue: data.InstallmentValue,
		InterestValue:    data.InterestValue,
		AssetName:        data.AssetName,
		Tenor:            data.Tenor,
		UpdatedAt:        nil,
	}

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Order{}).Create(&result).Error; err != nil {
			r.logger.Error("failed to create order: ", zap.Error(err))
			return err
		}
		return nil
	})

	return &result, nil
}

func (r *OrderRepository) Get(db *gorm.DB, data *models.OrderGetRequest) (*entity.Order, error) {
	var order entity.Order
	tx := db.Model(&order)

	if data.ID != 0 {
		tx = tx.Where("id = ?", data.ID)
	}

	if data.ContractNumber != "" {
		tx = tx.Where("no_contract = ?", data.ContractNumber)
	}

	err := tx.First(&order).Error
	if err != nil {
		r.logger.Error("failed to get order: ", zap.Error(err))
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) Update(db *gorm.DB, data *models.OrderUpdateRequest) error {
	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.Order{}).Updates(&models.OrderUpdateRequest{
			OTR:              data.OTR,
			AdminFee:         data.AdminFee,
			InstallmentValue: data.InstallmentValue,
			InterestValue:    data.InterestValue,
			AssetName:        data.AssetName,
		}).Where("id = ?", data.ID).Error; err != nil {
			r.logger.Error("failed to update order: ", zap.Error(err))
			return err
		}
		return nil
	})

	return nil
}

func (r *OrderRepository) Delete(db *gorm.DB, id uint64) error {
	if err := db.Delete(&entity.Order{}, db.Where("id = ?", id)).Error; err != nil {
		r.logger.Error("failed to delete order: ", zap.Error(err))
		return err
	}
	return nil
}

func (r *OrderRepository) GetAll(db *gorm.DB, data *models.OrderFilter) (*[]entity.Order, int64, error) {
	var result []entity.Order
	var count int64

	query := db.Model(&entity.Order{})

	if data.ContractNumber != "" {
		query = query.Where("no_contract = ?", data.ContractNumber)
	}

	if data.ID != 0 {
		query = query.Where("id = ?", data.ID)
	}

	err := query.Count(&count).Error

	if err != nil {
		r.logger.Error("failed to get all order: ", zap.Error(err))
		return nil, 0, err
	}

	err = query.Find(&result).Error

	if err != nil {
		r.logger.Error("failed to get all order: ", zap.Error(err))
		return nil, 0, err
	}

	return &result, count, nil
}

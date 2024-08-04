package entity

import "time"

type Order struct {
	ID               uint64     `gorm:"column:id;primaryKey"`
	CustomerID       uint64     `gorm:"column:customer_id"`
	ContractNumber   string     `gorm:"column:no_contract"`
	OTR              float64    `gorm:"column:otr"`
	AdminFee         float64    `gorm:"column:admin_fee"`
	InstallmentValue float64    `gorm:"column:installment_value"`
	InterestValue    float64    `gorm:"column:interest_value"`
	AssetName        string     `gorm:"column:asset_name"`
	Tenor            int        `gorm:"column:tenor"`
	CreatedAt        *time.Time `gorm:"column:created_at"`
	UpdatedAt        *time.Time `gorm:"column:updated_at"`
}

func (o *Order) TableName() string {
	return "orders"
}

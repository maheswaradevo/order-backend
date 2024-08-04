package models

type CreditLimitRequest struct {
	CustomerID uint64 `json:"customer_id"`
}

type CreditLimitUpdate struct {
	ID          uint64  `json:"id"`
	CreditLimit float64 `json:"credit_limit"`
}

package consumer

import "time"

type CreditLimitEvent struct {
	ID          uint      `json:"id"`
	CustomerID  uint64    `json:"customer_id"`
	CreditLimit float64   `json:"credit_limit"`
	TenorID     uint      `json:"tenor_id"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

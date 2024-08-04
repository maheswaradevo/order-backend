package models

type (
	OrderCreateRequest struct {
		OTR              float64 `json:"otr" form:"otr"`
		AdminFee         float64 `json:"admin_fee" form:"admin_fee"`
		InstallmentValue float64 `json:"installment_value" form:"installment_value"`
		InterestValue    float64 `json:"interest_value" form:"interest_value"`
		AssetName        string  `json:"asset_name" form:"asset_name"`
	}
	OrderGetRequest struct {
		ID             uint64 `query:"id"`
		ContractNumber string `query:"contract_number"`
	}

	OrderUpdateRequest struct {
		ID               uint64  `param:"id"`
		OTR              float64 `json:"otr" form:"otr"`
		AdminFee         float64 `json:"admin_fee" form:"admin_fee"`
		InstallmentValue float64 `json:"installment_value" form:"installment_value"`
		InterestValue    float64 `json:"interest_value" form:"interest_value"`
		AssetName        string  `json:"asset_name" form:"asset_name"`
	}
)

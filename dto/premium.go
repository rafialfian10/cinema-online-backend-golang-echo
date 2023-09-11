package dto

type CreatePremiumRequest struct {
	BuyerID int  `json:"buyer_id" form:"buyer_id" validate:"required"`
	Price   int  `json:"price" form:"price" validate:"required"`
	Status  bool `json:"status" form:"status" validate:"required"`
}

type UpdatePremiumRequest struct {
	Status bool `json:"status" form:"status"`
}

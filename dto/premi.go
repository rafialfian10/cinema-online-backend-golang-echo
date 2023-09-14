package dto

type UpdatePremiRequest struct {
	UserID int  `json:"user_id" form:"user_id" validate:"required"`
	Price  int  `json:"price" form:"price"`
	Status bool `json:"status" form:"status"`
}

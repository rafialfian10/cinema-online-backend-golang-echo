package dto

type CreateTransactionRequest struct {
	MovieID  int    `json:"movie_id" form:"movie_id" validate:"required"`
	BuyerID  int    `json:"buyer_id" form:"buyer_id" validate:"required"`
	SellerID int    `json:"seller_id" form:"seller_id" validate:"required"`
	Price    int    `json:"price" form:"price" validate:"required"`
	Status   string `json:"status" form:"status" validate:"required"`
}

type UpdateTransactionRequest struct {
	Status string `json:"status" form:"status"`
}

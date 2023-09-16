package dto

type UpdatePremiRequest struct {
	Price  int  `json:"price" form:"price"`
	Status bool `json:"status" form:"status"`
}

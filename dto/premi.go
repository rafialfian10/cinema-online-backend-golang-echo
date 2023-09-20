package dto

import "time"

type UpdatePremiRequest struct {
	Price  int  `json:"price" form:"price"`
	Status bool `json:"status" form:"status"`
}

type UpdatePremiExpiredRequest struct {
	OrderID     int       `json:"order_id"`
	Price       int       `json:"price" form:"price"`
	Status      bool      `json:"status" form:"status"`
	Token       string    `json:"token" form:"token"`
	ActivatedAt time.Time `json:"activated_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

package models

import "time"

type Transaction struct {
	ID        int           `json:"id" gorm:"primary_key:auto_increment"`
	MovieID   int           `json:"movie_id"`
	Movie     MovieResponse `json:"movie"`
	BuyerID   int           `json:"buyer_id"`
	Buyer     UserResponse  `json:"buyer"`
	SellerID  int           `json:"seller_id"`
	Seller    UserResponse  `json:"seller"`
	Price     int           `json:"price"`
	Status    string        `json:"status"  gorm:"type:varchar(25)"`
	Token     string        `json:"token" gorm:"type: varchar(255)"`
	CreatedAt time.Time     `json:"-"`
	UpdatedAt time.Time     `json:"-"`
}

type TransactionResponse struct {
	ID       int           `json:"id" gorm:"primary_key:auto_increment"`
	MovieID  int           `json:"movie_id"`
	Movie    MovieResponse `json:"movie"`
	BuyerID  int           `json:"buyer_id"`
	Buyer    UserResponse  `json:"buyer"`
	SellerID int           `json:"seller_id"`
	Seller   UserResponse  `json:"seller"`
	Price    int           `json:"price"`
	Status   string        `json:"status" gorm:"type:varchar(25)"`
	Token    string        `json:"token" gorm:"type: varchar(255)"`
}

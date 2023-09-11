package models

import "time"

type Premium struct {
	ID        int          `json:"id" gorm:"primary_key:auto_increment"`
	Status    bool         `json:"status"`
	Price     int          `json:"price"`
	Token     string       `json:"token" gorm:"type: varchar(255)"`
	BuyerID   int          `json:"buyer_id"`
	Buyer     UserResponse `json:"buyer"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
}

type PremiumResponse struct {
	ID      int          `json:"id" gorm:"primary_key:auto_increment"`
	Status  bool         `json:"status"`
	Price   int          `json:"price"`
	Token   string       `json:"token" gorm:"type: varchar(255)"`
	BuyerID int          `json:"buyer_id"`
	Buyer   UserResponse `json:"buyer"`
}

func (PremiumResponse) TableName() string {
	return "premiums"
}

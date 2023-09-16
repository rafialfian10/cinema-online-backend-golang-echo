package models

import "time"

type Premi struct {
	ID          int    `json:"id" gorm:"primary_key:auto_increment"`
	OrderID     int    `json:"order_id"`
	Status      bool   `json:"status"`
	Price       int    `json:"price"`
	Token       string `json:"token" gorm:"type: varchar(255)"`
	UserID      int    `json:"user_id"`
	User        UserResponse
	ActivatedAt time.Time `json:"activated_at"`
	ExpiredAt   time.Time `json:"expired_at"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

type PremiResponse struct {
	ID          int       `json:"id" gorm:"primary_key:auto_increment"`
	OrderID     int       `json:"order_id"`
	Status      bool      `json:"status"`
	Price       int       `json:"price"`
	Token       string    `json:"token" gorm:"type: varchar(255)"`
	UserID      int       `json:"user_id"`
	ActivatedAt time.Time `json:"activated_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func (PremiResponse) TableName() string {
	return "premis"
}

package models

import "time"

type Premi struct {
	ID        int    `json:"id" gorm:"primary_key:auto_increment"`
	Status    bool   `json:"status"`
	Price     int    `json:"price"`
	Token     string `json:"token" gorm:"type: varchar(255)"`
	UserID    int    `json:"user_id"`
	User      UserResponse
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type PremiResponse struct {
	ID     int    `json:"id" gorm:"primary_key:auto_increment"`
	Status bool   `json:"status"`
	Price  int    `json:"price"`
	Token  string `json:"token" gorm:"type: varchar(255)"`
	UserID int    `json:"user_id"`
}

func (PremiResponse) TableName() string {
	return "premis"
}

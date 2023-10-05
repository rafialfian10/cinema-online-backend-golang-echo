package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" gorm:"type: varchar(255)"`
	Email     string    `json:"email" gorm:"type: varchar(255)"`
	Password  string    `json:"password" gorm:"type: varchar(255)"`
	Gender    string    `json:"gender" gorm:"type: varchar(255)"`
	Phone     string    `json:"phone" gorm:"type: varchar(255)"`
	Address   string    `json:"address" gorm:"type: text"`
	Photo     string    `json:"photo" gorm:"type: varchar(255)"`
	Role      string    `json:"role" gorm:"type: varchar(255)"`
	Premi     Premi     `json:"premi" gorm:"foreignKey:UserID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID       int           `json:"id"`
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Gender   string        `json:"gender"`
	Phone    string        `json:"phone"`
	Address  string        `json:"address"`
	Photo    string        `json:"photo"`
	Role     string        `json:"role"`
	Premi    PremiResponse `json:"premi" gorm:"foreignKey:UserID"`
}

func (UserResponse) TableName() string {
	return "users"
}

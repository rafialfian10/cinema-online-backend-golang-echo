package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" gorm:"type: varchar(255)"`
	Email     string    `json:"email" gorm:"type: varchar(255)"`
	Password  string    `json:"password" gorm:"type: varchar(255)"`
	Role      string    `json:"role" gorm:"type: varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

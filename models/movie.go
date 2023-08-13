package models

import "time"

type Movie struct {
	ID          int       `json:"id" gorm:"primary_key:auto_increment"`
	Title       string    `json:"title" gorm:"type: varchar(255)"`
	Category    string    `json:"category" gorm:"type: varchar(255)"`
	Price       int       `json:"price" gorm:"type: int"`
	Link        string    `json:"link" gorm:"type: varchar(255)"`
	Description string    `json:"description" gorm:"type: text"`
	Thumbnail   string    `json:"thumbnail" gorm:"type: varchar(255)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MovieResponse struct {
	ID          int    `json:"id" gorm:"primary_key:auto_increment"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}

func (MovieResponse) TableName() string {
	return "movies"
}

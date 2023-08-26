package models

import "time"

type Movie struct {
	ID          int        `json:"id" gorm:"primary_key:auto_increment"`
	Title       string     `json:"title" gorm:"type: varchar(255)"`
	ReleaseDate time.Time  `json:"release_date"`
	Category    []Category `json:"category" gorm:"many2many:movie_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CategoryID  []int      `json:"category_id" form:"category_id" gorm:"-"`
	Price       int        `json:"price" gorm:"type: int"`
	Link        string     `json:"link" gorm:"type: varchar(255)"`
	Description string     `json:"description" gorm:"type: text"`
	Thumbnail   string     `json:"thumbnail" gorm:"type: varchar(255)"`
	Trailer     string     `json:"trailer" gorm:"type: varchar(255)"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type MovieResponse struct {
	ID          int        `json:"id" gorm:"primary_key:auto_increment"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"release_date"`
	Category    []Category `json:"category" gorm:"many2many:movie_categories"`
	CategoryID  []int      `json:"category_id" form:"category_id" gorm:"-"`
	Price       int        `json:"price"`
	Link        string     `json:"link"`
	Description string     `json:"description"`
	Thumbnail   string     `json:"thumbnail"`
	Trailer     string     `json:"trailer" gorm:"type: varchar(255)"`
}

func (MovieResponse) TableName() string {
	return "movies"
}

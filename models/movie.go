package models

import (
	"time"
)

type Movie struct {
	ID          int          `json:"id" gorm:"primary_key:auto_increment"`
	Title       string       `json:"title" gorm:"type: varchar(255)"`
	ReleaseDate time.Time    `json:"release_date"`
	CategoryID  []int        `json:"category_id" form:"category_id" gorm:"-"`
	Category    []Category   `json:"category" gorm:"many2many:movie_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Price       int          `json:"price" gorm:"type: int"`
	Link        string       `json:"link" gorm:"type: varchar(255)"`
	Description string       `json:"description" gorm:"type: text"`
	Thumbnail   string       `json:"thumbnail" gorm:"type: varchar(255)"`
	Trailer     string       `json:"trailer" gorm:"type: varchar(255)"`
	FullMovie   string       `json:"full_movie" gorm:"type: varchar(255)"`
	UserID      int          `json:"user_id" form:"user_id"`
	User        UserResponse `json:"user"`
	// RatingID    []int        `json:"rating_id" form:"rating_id" gorm:"-"`
	// Rating      []Rating     `json:"rating" gorm:"many2many:movie_ratings;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MovieResponse struct {
	ID          int                `json:"id" gorm:"primary_key:auto_increment"`
	Title       string             `json:"title"`
	ReleaseDate time.Time          `json:"release_date"`
	CategoryID  []int              `json:"category_id" form:"category_id" gorm:"-"`
	Category    []CategoryResponse `json:"category" gorm:"many2many:movie_categories"`
	Price       int                `json:"price"`
	Link        string             `json:"link"`
	Description string             `json:"description"`
	Thumbnail   string             `json:"thumbnail"`
	FullMovie   string             `json:"full_movie"`
	Trailer     string             `json:"trailer"`
	UserID      int                `json:"user_id" form:"user_id"`
	User        UserResponse       `json:"user"`
	// RatingID    []int              `json:"rating_id" form:"rating_id" gorm:"-"`
	// Rating      []RatingResponse   `json:"rating" gorm:"many2many:movie_ratings"`
}

func (MovieResponse) TableName() string {
	return "movies"
}

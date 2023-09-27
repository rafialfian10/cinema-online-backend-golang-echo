package models

import (
	"time"
)

type Rating struct {
	ID        int           `json:"id" gorm:"primary_key:auto_increment"`
	Star      int           `json:"star" gorm:"type: int"`
	MovieID   int           `json:"movie_id" form:"movie_id" gorm:"column:movie_id"`
	Movie     MovieResponse `json:"movie" gorm:"foreignkey:MovieID"`
	UserID    int           `json:"user_id" form:"user_id"`
	User      UserResponse  `json:"user"`
	CreatedAt time.Time     `json:"-"`
	UpdatedAt time.Time     `json:"-"`
}

type RatingResponse struct {
	ID      int           `json:"id" gorm:"primary_key:auto_increment"`
	Star    int           `json:"star"`
	MovieID int           `json:"movie_id"`
	Movie   MovieResponse `json:"movie"`
	UserID  int           `json:"user_id"`
	User    UserResponse  `json:"user"`
}

func (RatingResponse) TableName() string {
	return "ratings"
}

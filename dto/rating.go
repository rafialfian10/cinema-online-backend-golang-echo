package dto

type CreateRatingRequest struct {
	Star    int `json:"star" form:"star" validate:"required"`
	MovieID int `json:"movie_id" form:"movie_id" validate:"required"`
	UserID  int `json:"user_id" form:"user_id" validate:"required"`
}

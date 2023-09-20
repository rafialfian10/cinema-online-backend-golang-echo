package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type RatingRepository interface {
	FindRatings() ([]models.Rating, error)
	GetRating(ID int) (models.Rating, error)
	CreateRating(rating models.Rating) (models.Rating, error)
	DeleteRating(rating models.Rating, ID int) (models.Rating, error)
	GetMovie(movieId int) (models.Movie, error)
}

type ratingRepository struct {
	db *gorm.DB
}

func RepositoryRating(db *gorm.DB) *ratingRepository {
	return &ratingRepository{db}
}

func (r *ratingRepository) FindRatings() ([]models.Rating, error) {
	var ratings []models.Rating
	err := r.db.Preload("Movie.User.Premi").Preload("User.Premi").Find(&ratings).Error

	return ratings, err
}

func (r *ratingRepository) GetRating(ID int) (models.Rating, error) {
	var rating models.Rating
	err := r.db.Preload("Movie.User.Premi").Preload("User.Premi").First(&rating, ID).Error

	return rating, err
}

func (r *ratingRepository) CreateRating(rating models.Rating) (models.Rating, error) {
	err := r.db.Create(&rating).Error

	return rating, err
}

func (r *ratingRepository) DeleteRating(rating models.Rating, ID int) (models.Rating, error) {
	err := r.db.Delete(&rating, ID).Scan(&rating).Error

	return rating, err
}

func (r *ratingRepository) GetMovie(movieId int) (models.Movie, error) {
	var movie models.Movie
	err := r.db.First(&movie, movieId).Error

	return movie, err
}

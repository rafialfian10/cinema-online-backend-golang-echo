package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type MovieRepository interface {
	FindMovies() ([]models.Movie, error)
	GetMovie(ID int) (models.Movie, error)
	CreateMovie(movie models.Movie) (models.Movie, error)
	UpdateMovie(movie models.Movie) (models.Movie, error)
	DeleteMovie(movie models.Movie, ID int) (models.Movie, error)
}

type movieRepository struct {
	db *gorm.DB
}

func RepositoryMovie(db *gorm.DB) *movieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) FindMovies() ([]models.Movie, error) {
	var movies []models.Movie
	err := r.db.Find(&movies).Error

	return movies, err
}

func (r *movieRepository) GetMovie(ID int) (models.Movie, error) {
	var movie models.Movie
	err := r.db.First(&movie, ID).Error

	return movie, err
}

func (r *movieRepository) CreateMovie(movie models.Movie) (models.Movie, error) {
	err := r.db.Create(&movie).Error

	return movie, err
}

func (r *movieRepository) UpdateMovie(movie models.Movie) (models.Movie, error) {
	err := r.db.Debug().Model(&movie).Updates(movie).Error
	return movie, err
}

func (r *movieRepository) DeleteMovie(movie models.Movie, ID int) (models.Movie, error) {
	err := r.db.Raw("DELETE FROM movies WHERE id=?", ID).Scan(&movie).Error

	return movie, err
}

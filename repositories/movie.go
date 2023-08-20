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
	FindCategoriesById(categoriesId []int) ([]models.Category, error)
	DeleteMovieCategoryByMovieId(movie models.Movie) (models.Movie, error)
	DeleteThumbnailByID(ID int) error
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

func (r *movieRepository) FindCategoriesById(categoriesId []int) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories, categoriesId).Error

	return categories, err
}

func (r *movieRepository) DeleteMovieCategoryByMovieId(movie models.Movie) (models.Movie, error) {
	r.db.Exec("DELETE FROM movie_categories WHERE movie_id=?", movie.ID)
	err := r.db.Preload("User").Preload("Category").First(&movie, movie.ID).Error // add this code

	return movie, err
}

func (r *movieRepository) DeleteThumbnailByID(ID int) error {
	return r.db.Model(&models.Movie{}).Where("id = ?", ID).UpdateColumn("thumbnail", gorm.Expr("NULL")).Error
}

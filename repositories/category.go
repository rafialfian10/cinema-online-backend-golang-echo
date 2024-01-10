package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindCategories() ([]models.Category, error)
	GetCategory(ID int) (models.Category, error)
	CreateCategory(category models.Category) (models.Category, error)
	UpdateCategory(category models.Category) (models.Category, error)
	DeleteCategory(category models.Category, ID int) (models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func RepositoryCategory(db *gorm.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("id DESC").Find(&categories).Error

	return categories, err
}

func (r *categoryRepository) GetCategory(ID int) (models.Category, error) {
	var category models.Category
	err := r.db.First(&category, ID).Error

	return category, err
}

func (r *categoryRepository) CreateCategory(category models.Category) (models.Category, error) {
	err := r.db.Create(&category).Error

	return category, err
}

func (r *categoryRepository) UpdateCategory(category models.Category) (models.Category, error) {
	err := r.db.Save(&category).Error

	return category, err
}

func (r *categoryRepository) DeleteCategory(category models.Category, ID int) (models.Category, error) {
	err := r.db.Delete(&category, ID).Scan(&category).Error

	return category, err
}

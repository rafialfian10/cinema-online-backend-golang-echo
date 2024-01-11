package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUsers() ([]models.User, error)
	GetUser(ID int) (models.User, error)
	CreateUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(user models.User, ID int) (models.User, error)
	GetProfile(userId int) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func RepositoryUser(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) FindUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Premi").Order("id DESC").Find(&users).Error

	return users, err
}

func (r *userRepository) GetUser(ID int) (models.User, error) {
	var user models.User
	err := r.db.Preload("Premi").First(&user, ID).Error

	return user, err
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Debug().Model(&user).Updates(user).Error
	return user, err
}

func (r *userRepository) DeleteUser(user models.User, ID int) (models.User, error) {
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("user_id = ?", ID).Delete(&models.Premi{}).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	if err := tx.Delete(&user, ID).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	return user, tx.Commit().Error
}

func (r *userRepository) GetProfile(userId int) (models.User, error) {
	var profile models.User
	err := r.db.Preload("Premi").First(&profile, userId).Error

	return profile, err
}

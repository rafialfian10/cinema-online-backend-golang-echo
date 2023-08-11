package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type VerifyEmailRepository interface {
	CreateVerificationToken(userID int, token string) error
	GetUserByVerificationToken(token string) (models.User, error)
	UpdateEmailVerificationStatus(userID int, verified bool) error
}

type verifyEmailRepository struct {
	db *gorm.DB
}

func RepositoryVerifyEmail(db *gorm.DB) VerifyEmailRepository {
	return &verifyEmailRepository{db}
}

func (r *verifyEmailRepository) CreateVerificationToken(userID int, token string) error {
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("verification_token", token).Error
	return err
}

func (r *verifyEmailRepository) GetUserByVerificationToken(token string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "verification_token = ?", token).Error
	return user, err
}

func (r *verifyEmailRepository) UpdateEmailVerificationStatus(userID int, verified bool) error {
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("verified", verified).Error
	return err
}

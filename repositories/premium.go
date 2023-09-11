package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type PremiumRepository interface {
	FindPremiums() ([]models.Premium, error)
	GetPremium(Id int) (models.Premium, error)
	CreatePremium(premium models.Premium) (models.Premium, error)
	UpdatePremium(status bool, orderId int) (models.Premium, error)
	UpdateTokenPremium(token string, Id int) (models.Premium, error)
	DeletePremium(premium models.Premium, ID int) (models.Premium, error)
}

type premiumRepository struct {
	db *gorm.DB
}

func RepositoryPremium(db *gorm.DB) *premiumRepository {
	return &premiumRepository{db}
}

func (r *premiumRepository) FindPremiums() ([]models.Premium, error) {
	var premiums []models.Premium
	err := r.db.Preload("Buyer").Order("id desc").Find(&premiums).Error

	return premiums, err
}

func (r *premiumRepository) GetPremium(Id int) (models.Premium, error) {
	var premium models.Premium
	err := r.db.Preload("Buyer").First(&premium, "id = ?", Id).Error

	return premium, err
}

func (r *premiumRepository) CreatePremium(premium models.Premium) (models.Premium, error) {
	err := r.db.Create(&premium).Error

	return premium, err
}

func (r *premiumRepository) UpdatePremium(status bool, orderId int) (models.Premium, error) {
	var premium models.Premium
	r.db.Preload("Buyer").First(&premium, orderId)

	premium.Status = status
	err := r.db.Save(&premium).Error
	return premium, err
}

func (r *premiumRepository) UpdateTokenPremium(token string, Id int) (models.Premium, error) {
	var premium models.Premium
	r.db.Preload("Buyer").First(&premium, "id = ?", Id)

	// mengubah premium token
	premium.Token = token

	err := r.db.Model(&premium).Updates(premium).Error

	return premium, err
}

func (r *premiumRepository) DeletePremium(premium models.Premium, ID int) (models.Premium, error) {
	err := r.db.Delete(&premium, ID).Scan(&premium).Error

	return premium, err
}

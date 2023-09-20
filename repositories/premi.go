package repositories

import (
	"cinemaonline/models"
	"time"

	"gorm.io/gorm"
)

type PremiRepository interface {
	FindPremis() ([]models.Premi, error)
	GetPremi(Id int) (models.Premi, error)
	GetPremiOrderId(orderId int) (models.Premi, error)
	UpdatePremiUserStatus(status bool, orderId int) (models.Premi, error)
	UpdatePremiUser(premi models.Premi, userId int) (models.Premi, error)
	UpdatePremiExpired(premi models.Premi, userId int) (models.Premi, error)
	UpdateTokenPremi(token string, Id int) (models.Premi, error)
	DeletePremi(premi models.Premi, ID int) (models.Premi, error)
}

type premiRepository struct {
	db *gorm.DB
}

func RepositoryPremi(db *gorm.DB) *premiRepository {
	return &premiRepository{db}
}

func (r *premiRepository) FindPremis() ([]models.Premi, error) {
	var premis []models.Premi
	err := r.db.Preload("User").Find(&premis).Error

	return premis, err
}

func (r *premiRepository) GetPremi(Id int) (models.Premi, error) {
	var premi models.Premi
	err := r.db.Preload("User").First(&premi, "id = ?", Id).Error

	return premi, err
}

func (r *premiRepository) GetPremiOrderId(orderId int) (models.Premi, error) {
	var premi models.Premi
	err := r.db.Preload("User").First(&premi, "order_id = ?", orderId).Error

	return premi, err
}

func (r *premiRepository) UpdatePremiUser(premiUpdate models.Premi, Id int) (models.Premi, error) {
	var premi models.Premi
	r.db.First(&premi, "id = ?", Id)

	premi = premiUpdate
	err := r.db.Save(&premi).Error

	return premi, err
}

func (r *premiRepository) UpdatePremiUserStatus(status bool, orderId int) (models.Premi, error) {
	var premi models.Premi
	r.db.First(&premi, orderId)

	premi.Status = status

	if premi.Status {
		premi.ActivatedAt = time.Now()
		premi.ExpiredAt = time.Now().AddDate(0, 1, 0) // insert 30 day
	} else {
		premi.ActivatedAt = time.Time{}
		premi.ExpiredAt = time.Time{}
	}

	err := r.db.Save(&premi).Error
	return premi, err
}

func (r *premiRepository) UpdatePremiExpired(premiUpdate models.Premi, Id int) (models.Premi, error) {
	var premi models.Premi
	r.db.First(&premi, "id = ?", Id)

	premi = premiUpdate
	err := r.db.Save(&premi).Error

	return premi, err
}

func (r *premiRepository) UpdateTokenPremi(token string, Id int) (models.Premi, error) {
	var premi models.Premi
	r.db.First(&premi, "id = ?", Id)

	// mengubah premium token
	premi.Token = token
	err := r.db.Model(&premi).Updates(premi).Error

	return premi, err
}

func (r *premiRepository) DeletePremi(premi models.Premi, ID int) (models.Premi, error) {
	err := r.db.Delete(&premi, ID).Scan(&premi).Error

	return premi, err
}

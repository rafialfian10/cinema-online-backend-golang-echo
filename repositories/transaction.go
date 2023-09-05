package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactionsByUser(UserId int) ([]models.Transaction, error)
	// FindTransactions(userId int) ([]models.Transaction, error)
	GetTransaction(transactionId int) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(status string, orderId int) (models.Transaction, error)
	UpdateTokenTransaction(token string, Id int) (models.Transaction, error)
	DeleteTransaction(transaction models.Transaction) (models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func RepositoryTransaction(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) FindTransactionsByUser(userId int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("buyer_id=?", userId).Preload("Movie.User").Preload("Buyer").Preload("Seller").Order("id desc").Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) GetTransaction(transactionId int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Movie").Preload("Buyer").Preload("Seller").First(&transaction, transactionId).Error

	return transaction, err
}

func (r *transactionRepository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Movie").Preload("Buyer").Preload("Seller").Create(&transaction).Error

	return transaction, err
}

func (r *transactionRepository) UpdateTransaction(status string, orderId int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Movie").Preload("Buyer").Preload("Seller").First(&transaction, orderId)

	if status != transaction.Status && status == "success" {
		var movie models.Movie
		r.db.First(&movie, transaction.Movie.ID)
		r.db.Save(&movie)
	}

	transaction.Status = status
	err := r.db.Save(&transaction).Error
	return transaction, err
}

func (r *transactionRepository) UpdateTokenTransaction(token string, Id int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Movie.Category").Preload("Movie").Preload("User").First(&transaction, "id = ?", Id)

	// mengubah transaction token
	transaction.Token = token

	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *transactionRepository) DeleteTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Trip").Delete(&transaction).Error

	return transaction, err
}

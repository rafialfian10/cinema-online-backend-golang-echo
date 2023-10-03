package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactionsByUser(UserId int) ([]models.Transaction, error)
	FindTransactions() ([]models.Transaction, error)
	GetTransaction(transactionId int) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(status string, orderId int) (models.Transaction, error)
	UpdateTokenTransaction(token string, transactionId int) (models.Transaction, error)
	DeleteTransaction(transaction models.Transaction, ID int) (models.Transaction, error)
	GetMovie(movieId int) (models.Movie, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func RepositoryTransaction(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) FindTransactionsByUser(userId int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Where("buyer_id=?", userId).Preload("Movie.User.Premi").Preload("Movie.Rating.User.Premi").Preload("Buyer.Premi").Preload("Seller.Premi").Find(&transactions).Error

	// for i := range transactions {
	// 	var movie models.Movie
	// 	if err := r.db.Model(&transactions[i].Movie).Association("Category").Find(&movie.Category); err != nil {
	// 		return nil, err
	// 	}

	// 	var categoryResponses []models.CategoryResponse
	// 	for _, category := range movie.Category {
	// 		categoryResponses = append(categoryResponses, models.CategoryResponse{
	// 			ID:   category.ID,
	// 			Name: category.Name,
	// 		})
	// 	}

	// 	transactions[i].Movie.Category = categoryResponses
	// }

	return transactions, err
}

func (r *transactionRepository) FindTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Movie.User.Premi").Preload("Movie.Rating.User.Premi").Preload("Buyer.Premi").Preload("Seller.Premi").Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) GetTransaction(transactionId int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Movie.User.Premi").Preload("Movie.Rating.User.Premi").Preload("Buyer.Premi").Preload("Seller.Premi").First(&transaction, "id = ?", transactionId).Error

	return transaction, err
}

func (r *transactionRepository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("Movie.User").Preload("Buyer").Preload("Seller").Create(&transaction).Error

	return transaction, err
}

func (r *transactionRepository) UpdateTransaction(status string, orderId int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Movie.User").Preload("Buyer").Preload("Seller").First(&transaction, orderId)

	transaction.Status = status
	err := r.db.Save(&transaction).Error
	return transaction, err
}

func (r *transactionRepository) UpdateTokenTransaction(token string, transactionId int) (models.Transaction, error) {
	var transaction models.Transaction
	r.db.Preload("Movie.User").Preload("Buyer").Preload("Seller").First(&transaction, "id = ?", transactionId)

	// mengubah transaction token
	transaction.Token = token
	err := r.db.Model(&transaction).Updates(transaction).Error

	return transaction, err
}

func (r *transactionRepository) DeleteTransaction(transaction models.Transaction, ID int) (models.Transaction, error) {
	err := r.db.Raw("DELETE FROM transactions WHERE id=?", ID).Scan(&transaction).Error

	return transaction, err
}

func (r *transactionRepository) GetMovie(movieId int) (models.Movie, error) {
	var movie models.Movie
	err := r.db.First(&movie, movieId).Error

	return movie, err
}

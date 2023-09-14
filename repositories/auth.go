package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(user models.User) (models.User, error)
	FindUserByUsernameOrEmail(username, email string) (models.User, error)
	Login(email string) (models.User, error)
	CheckAuth(ID int) (models.User, error)
	CreateUserPremi(premi models.Premi) (models.Premi, error)
	GetPremi(Id int) (models.Premi, error)
}

type authRepository struct {
	db *gorm.DB
}

// membuat function RepositoryAuth. parameter pointer ke gorm, return repository{db}. ini akan dipanggil di routes
func RepositoryAuth(db *gorm.DB) *authRepository {
	return &authRepository{db}
}

// function register
func (r *authRepository) Register(user models.User) (models.User, error) {
	err := r.db.Preload("Premi").Create(&user).Error

	return user, err
}

// function check data username & email
func (r *authRepository) FindUserByUsernameOrEmail(username, email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "username=? OR email=?", username, email).Error

	return user, err
}

// function login
func (r *authRepository) Login(email string) (models.User, error) {
	var user models.User

	// ambil data user yang email user == request email
	err := r.db.First(&user, "email=?", email).Error

	return user, err
}

func (r *authRepository) CheckAuth(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

func (r *authRepository) CreateUserPremi(premi models.Premi) (models.Premi, error) {
	err := r.db.Create(&premi).Error

	return premi, err
}

func (r *authRepository) GetPremi(Id int) (models.Premi, error) {
	var premi models.Premi
	err := r.db.First(&premi, "id = ?", Id).Error

	return premi, err
}

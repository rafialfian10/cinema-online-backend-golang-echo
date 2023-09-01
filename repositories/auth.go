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
	Getuser(ID int) (models.User, error)
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
	err := r.db.Create(&user).Error

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
	err := r.db.Preload("movies").First(&user, ID).Error

	return user, err
}

// function get user id
func (r *authRepository) Getuser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

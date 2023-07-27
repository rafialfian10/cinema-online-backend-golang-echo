package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(user models.User) (models.User, error)
	FindUserByUsernameOrEmail(username, email string) (models.User, error)
	Login(email string) (models.User, error)
	Getuser(ID int) (models.User, error)
}

// membuat function RepositoryAuth. parameter pointer ke gorm, return repository{db}. ini akan dipanggil di routes
func RepositoryAuth(db *gorm.DB) *repository {
	return &repository{db}
}

// function register
func (r *repository) Register(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

// function check data username & email
func (r *repository) FindUserByUsernameOrEmail(username, email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "username=? OR email=?", username, email).Error

	return user, err
}

// function login
func (r *repository) Login(email string) (models.User, error) {
	var user models.User

	// ambil data user yang email user == request email
	err := r.db.First(&user, "email=?", email).Error

	return user, err
}

// function get user id
func (r *repository) Getuser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

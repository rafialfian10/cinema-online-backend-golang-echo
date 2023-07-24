package repositories

import (
	"cinemaonline/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(user models.User) (models.User, error)
	Login(email string) (models.User, error)
	Getuser(ID int) (models.User, error)
}

// membuat function RepositoryAuth. parameter pointer ke gorm, return repository{db}. ini akan dipanggil di routes
func RepositoryAuth(db *gorm.DB) *repository {
	return &repository{db}
}

// membuat struct method Register(memanggil struct dengan struct function)
func (r *repository) Register(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error

	return user, err
}

// membuat struct method Login(memanggil struct dengan struct function)
func (r *repository) Login(email string) (models.User, error) {
	var user models.User

	// ambil data user yang email user == request email
	err := r.db.First(&user, "email=?", email).Error

	return user, err
}

func (r *repository) Getuser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error

	return user, err
}

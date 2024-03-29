package database

import (
	"cinemaonline/models"
	"cinemaonline/pkg/mysql"
	"fmt"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Category{},
		&models.Rating{},
		&models.Transaction{},
		&models.Premi{},
	)

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}

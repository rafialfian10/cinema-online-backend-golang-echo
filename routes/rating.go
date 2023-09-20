package routes

import (
	"cinemaonline/handlers"
	"cinemaonline/pkg/middleware"
	"cinemaonline/pkg/mysql"
	"cinemaonline/repositories"

	"github.com/labstack/echo/v4"
)

func RatingRoutes(e *echo.Group) {
	ratingRepository := repositories.RepositoryRating(mysql.DB)
	h := handlers.HandlerRating(ratingRepository)

	e.GET("/ratings", h.FindRatings)
	e.GET("/rating/:id", h.GetRating)
	e.POST("/rating", middleware.Auth(h.CreateRating))
	e.DELETE("/rating/:id", middleware.Auth(h.DeleteRating))
}

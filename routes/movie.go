package routes

import (
	"cinemaonline/handlers"
	"cinemaonline/pkg/middleware"
	"cinemaonline/pkg/mysql"
	"cinemaonline/repositories"

	"github.com/labstack/echo/v4"
)

func MovieRoutes(e *echo.Group) {
	movieRepository := repositories.RepositoryMovie(mysql.DB)
	h := handlers.HandlerMovie(movieRepository)

	e.GET("/movies", h.FindMovies)
	e.GET("/movie/:id", h.GetMovie)
	e.POST("/movie", middleware.UploadFile(h.CreateMovie))
	e.PATCH("/movie/:id", middleware.UploadFile(h.UpdateMovie))
	e.DELETE("/movie/:id", h.DeleteMovie)
}
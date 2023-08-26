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
	e.POST("/movie", middleware.UploadVideo(middleware.UploadImage(h.CreateMovie)))
	e.PATCH("/movie/:id", middleware.UploadVideo(middleware.UploadImage(h.UpdateMovie)))
	e.DELETE("/movie/:id", h.DeleteMovie)
	e.DELETE("/movie/:id/thumbnail", h.DeleteThumbnail)
	e.DELETE("/movie/:id/trailer", h.DeleteTrailer)
}

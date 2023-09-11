package handlers

import (
	dto "cinemaonline/dto"
	"cinemaonline/models"
	"cinemaonline/repositories"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var path_thumbnail = "http://localhost:5000/uploads/thumbnail/"
var path_trailer = "http://localhost:5000/uploads/trailer/"

type handlerMovie struct {
	MovieRepository repositories.MovieRepository
}

func HandlerMovie(MovieRepository repositories.MovieRepository) *handlerMovie {
	return &handlerMovie{MovieRepository}
}

// function get all movie
func (h *handlerMovie) FindMovies(c echo.Context) error {
	movies, err := h.MovieRepository.FindMovies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	for i, movie := range movies {
		movies[i].Thumbnail = path_thumbnail + movie.Thumbnail
		movies[i].Trailer = path_trailer + movie.Trailer
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: movies})
}

// function get movie by id
func (h *handlerMovie) GetMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var movie models.Movie
	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	movie.Thumbnail = path_thumbnail + movie.Thumbnail
	movie.Trailer = path_trailer + movie.Trailer

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertMovieResponse(movie)})

}

// function create user
func (h *handlerMovie) CreateMovie(c echo.Context) error {
	var err error
	dataThumbnail := c.Get("dataThumbnail").(string)
	dataTrailer := c.Get("dataTrailer").(string)
	// fmt.Println("this is data file", dataThumbnail)
	// fmt.Println("this is data file", dataTrailer)

	price, _ := strconv.Atoi(c.FormValue("price"))
	categoryIdString := c.FormValue("category_id")

	if categoryIdString == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "Error: category_id form value is missing."})
	}

	var categoriesId []int
	err = json.Unmarshal([]byte(categoryIdString), &categoriesId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	if len(categoriesId) == 0 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "Error: category_id form value is missing."})
	}

	request := dto.CreateMovieRequest{
		Title:       c.FormValue("title"),
		ReleaseDate: c.FormValue("release_date"),
		CategoryID:  categoriesId,
		Price:       price,
		Link:        c.FormValue("link"),
		Description: c.FormValue("description"),
		Thumbnail:   dataThumbnail,
		Trailer:     dataTrailer,
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	ReleaseDate, _ := time.Parse("2006-01-02", c.FormValue("release_date"))
	categories, _ := h.MovieRepository.FindCategoriesById(request.CategoryID)

	movie := models.Movie{
		Title:       request.Title,
		ReleaseDate: ReleaseDate,
		Category:    categories,
		Price:       request.Price,
		Link:        request.Link,
		Description: request.Description,
		Thumbnail:   request.Thumbnail,
		Trailer:     request.Trailer,
		UserID:      int(userId),
	}

	movie, err = h.MovieRepository.CreateMovie(movie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	movie, _ = h.MovieRepository.GetMovie(movie.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertMovieResponse(movie)})
}

// function update movie
func (h *handlerMovie) UpdateMovie(c echo.Context) error {
	var err error
	dataThumbnail := c.Get("dataThumbnail").(string)
	dataTrailer := c.Get("dataTrailer").(string)

	price, _ := strconv.Atoi(c.FormValue("price"))

	var categoriesId []int
	categoryIdString := c.FormValue("category_id")
	err = json.Unmarshal([]byte(categoryIdString), &categoriesId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	request := dto.UpdateMovieRequest{
		Title:       c.FormValue("title"),
		ReleaseDate: c.FormValue("release_date"),
		CategoryID:  categoriesId,
		Price:       price,
		Link:        c.FormValue("link"),
		Description: c.FormValue("description"),
		Thumbnail:   dataThumbnail,
		Trailer:     dataTrailer,
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	id, _ := strconv.Atoi(c.Param("id"))
	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Title != "" {
		movie.Title = request.Title
	}

	date, _ := time.Parse("2006-01-02", request.ReleaseDate)
	time := time.Now()
	if date != time {
		movie.ReleaseDate = date
	}

	if len(request.CategoryID) == 0 {
		data, err := h.MovieRepository.DeleteMovieCategoryByMovieId(movie)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
	}

	categories, _ := h.MovieRepository.FindCategoriesById(request.CategoryID)
	movie.Category = categories

	if request.Price != 0 {
		movie.Price = request.Price
	}

	if request.Link != "" {
		movie.Link = request.Link
	}

	if request.Description != "" {
		movie.Description = request.Description
	}

	if request.Thumbnail != "" {
		movie.Thumbnail = request.Thumbnail
	}

	if request.Trailer != "" {
		movie.Trailer = request.Trailer
	}

	data, err := h.MovieRepository.UpdateMovie(movie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

// function delete movie
func (h *handlerMovie) DeleteMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.MovieRepository.DeleteMovie(movie, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertMovieResponse(data)})
}

// function delete thumbnail by id movie
func (h *handlerMovie) DeleteThumbnail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Delete the thumbnail using repository function
	if err := h.MovieRepository.DeleteThumbnailByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Get the updated movie data after deleting thumbnail
	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertMovieResponse(movie)})
}

// function delete thumbnail by id movie
func (h *handlerMovie) DeleteTrailer(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Delete the thumbnail using repository function
	if err := h.MovieRepository.DeleteTrailerByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// Get the updated movie data after deleting trailer
	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertMovieResponse(movie)})
}

// function convert movie response
func convertMovieResponse(movie models.Movie) models.MovieResponse {
	var result models.MovieResponse
	result.ID = movie.ID
	result.Title = movie.Title
	result.ReleaseDate = movie.ReleaseDate
	result.Price = movie.Price
	result.Link = movie.Link
	result.Description = movie.Description
	result.Thumbnail = movie.Thumbnail
	result.Trailer = movie.Trailer
	result.UserID = movie.UserID
	result.User = movie.User

	for _, cat := range movie.Category {
		categoryResponse := models.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
		result.Category = append(result.Category, categoryResponse)
	}

	return result
}

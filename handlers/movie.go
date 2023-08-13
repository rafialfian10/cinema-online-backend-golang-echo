package handlers

import (
	dto "cinemaonline/dto"
	"cinemaonline/models"
	"cinemaonline/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

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
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: movies})
}

// function get movie by id
func (h *handlerMovie) GetMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertMovieResponse(movie)})
}

// function create user
func (h *handlerMovie) CreateMovie(c echo.Context) error {
	var err error
	dataFile := c.Get("dataFile").(string)
	// fmt.Println("this is data file", dataFile)

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
		CategoryID:  categoriesId,
		Price:       price,
		Link:        c.FormValue("link"),
		Description: c.FormValue("description"),
		Thumbnail:   dataFile,
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	// userLogin := c.Get("userLogin")
	// userId := userLogin.(jwt.MapClaims)["id"].(float64)

	categories, _ := h.MovieRepository.FindCategoriesById(request.CategoryID)

	movie := models.Movie{
		Title:       request.Title,
		Category:    categories,
		Price:       request.Price,
		Link:        request.Link,
		Description: request.Description,
		Thumbnail:   request.Thumbnail,
		// UserID:      int(userId),
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
	dataFile := c.Get("dataFile").(string)

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
		CategoryID:  categoriesId,
		Price:       price,
		Link:        c.FormValue("link"),
		Description: c.FormValue("description"),
		Thumbnail:   dataFile,
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

	if len(request.CategoryID) == 0 {
		data, err := h.MovieRepository.DeleteMovieCategoryByMovieId(movie)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
		}

		return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
	}

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

	categories, _ := h.MovieRepository.FindCategoriesById(request.CategoryID)
	movie.Category = categories

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

// function convert movie response
func convertMovieResponse(movie models.Movie) models.MovieResponse {
	return models.MovieResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Category:    movie.Category,
		Price:       movie.Price,
		Link:        movie.Link,
		Description: movie.Description,
		Thumbnail:   movie.Thumbnail,
	}
}

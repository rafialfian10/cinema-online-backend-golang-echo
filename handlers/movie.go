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
var path_full_movie = "http://localhost:5000/uploads/full_movie/"

type handlerMovie struct {
	MovieRepository repositories.MovieRepository
}

func HandlerMovie(MovieRepository repositories.MovieRepository) *handlerMovie {
	return &handlerMovie{MovieRepository}
}

func (h *handlerMovie) FindMovies(c echo.Context) error {
	movies, err := h.MovieRepository.FindMovies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	for i, movie := range movies {
		movies[i].Thumbnail = path_thumbnail + movie.Thumbnail
		movies[i].Trailer = path_trailer + movie.Trailer
		movies[i].FullMovie = path_full_movie + movie.FullMovie
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMultipleMovieResponse(movies)})
}

func (h *handlerMovie) GetMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var movie models.Movie
	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	movie.Thumbnail = path_thumbnail + movie.Thumbnail
	movie.Trailer = path_trailer + movie.Trailer
	movie.FullMovie = path_full_movie + movie.FullMovie

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMovieResponse(movie)})
}

func (h *handlerMovie) CreateMovie(c echo.Context) error {
	var err error
	dataThumbnail := c.Get("dataThumbnail").(string)
	dataTrailer := c.Get("dataTrailer").(string)
	dataFullMovie := c.Get("dataFullMovie").(string)
	// fmt.Println("this is data file", dataThumbnail)
	// fmt.Println("this is data file", dataTrailer)
	// fmt.Println("this is data file", dataFullMovie)

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
		FullMovie:   dataFullMovie,
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
		FullMovie:   request.FullMovie,
		UserID:      int(userId),
	}

	movie, err = h.MovieRepository.CreateMovie(movie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	movie, _ = h.MovieRepository.GetMovie(movie.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMovieResponse(movie)})
}

func (h *handlerMovie) UpdateMovie(c echo.Context) error {
	var err error
	dataThumbnail := c.Get("dataThumbnail").(string)
	dataTrailer := c.Get("dataTrailer").(string)
	dataFullMovie := c.Get("dataFullMovie").(string)

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
		FullMovie:   dataFullMovie,
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

	if request.FullMovie != "" {
		movie.FullMovie = request.FullMovie
	}

	data, err := h.MovieRepository.UpdateMovie(movie)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

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

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMovieResponse(data)})
}

func (h *handlerMovie) DeleteThumbnail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.MovieRepository.DeleteThumbnailByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMovieResponse(movie)})
}

func (h *handlerMovie) DeleteTrailer(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.MovieRepository.DeleteTrailerByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMovieResponse(movie)})
}

func (h *handlerMovie) DeleteFullMovie(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.MovieRepository.DeleteFullMovieByID(id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	movie, err := h.MovieRepository.GetMovie(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMovieResponse(movie)})
}

func ConvertMovieResponse(movie models.Movie) models.MovieResponse {
	var result models.MovieResponse
	result.ID = movie.ID
	result.Title = movie.Title
	result.ReleaseDate = movie.ReleaseDate
	result.CategoryID = movie.CategoryID
	result.Price = movie.Price
	result.Link = movie.Link
	result.Description = movie.Description
	result.Thumbnail = movie.Thumbnail
	result.Trailer = movie.Trailer
	result.FullMovie = movie.FullMovie
	result.UserID = movie.UserID
	result.User = movie.User
	result.Rating = make([]models.RatingResponse, len(movie.Rating))
	for i, rating := range movie.Rating {
		result.Rating[i] = models.RatingResponse{
			ID:      rating.ID,
			Star:    rating.Star,
			MovieID: rating.MovieID,
			UserID:  rating.UserID,
			User:    rating.User,
		}
	}

	for _, cat := range movie.Category {
		categoryResponse := models.CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}
		result.Category = append(result.Category, categoryResponse)
	}

	return result
}

func ConvertMultipleMovieResponse(movies []models.Movie) []models.MovieResponse {
	var result []models.MovieResponse

	for _, movie := range movies {
		movies := models.MovieResponse{
			ID:          movie.ID,
			Title:       movie.Title,
			ReleaseDate: movie.ReleaseDate,
			CategoryID:  movie.CategoryID,
			Price:       movie.Price,
			Link:        movie.Link,
			Description: movie.Description,
			Thumbnail:   movie.Thumbnail,
			Trailer:     movie.Trailer,
			FullMovie:   movie.FullMovie,
			UserID:      movie.UserID,
			User:        movie.User,
			// RatingID:    movie.RatingID,
		}

		movies.Rating = []models.RatingResponse{}

		for _, rating := range movie.Rating {
			movies.Rating = append(movies.Rating, models.RatingResponse{
				ID:      rating.ID,
				Star:    rating.Star,
				MovieID: rating.MovieID,
				UserID:  rating.UserID,
				User:    rating.User,
			})
		}

		for _, cat := range movie.Category {
			categoryResponse := models.CategoryResponse{
				ID:   cat.ID,
				Name: cat.Name,
			}
			movies.Category = append(movies.Category, categoryResponse)
		}

		result = append(result, movies)
	}

	return result
}

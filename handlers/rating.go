package handlers

import (
	dto "cinemaonline/dto"
	"cinemaonline/models"
	"cinemaonline/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type handlerRating struct {
	RatingRepository repositories.RatingRepository
}

func HandlerRating(RatingRepository repositories.RatingRepository) *handlerRating {
	return &handlerRating{RatingRepository}
}

// function get all rating
func (h *handlerRating) FindRatings(c echo.Context) error {
	ratings, err := h.RatingRepository.FindRatings()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: ConvertMultipleRatingResponse(ratings)})
}

// function get by id rating
func (h *handlerRating) GetRating(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	rating, err := h.RatingRepository.GetRating(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertRatingResponse(rating)})
}

// function create rating
func (h *handlerRating) CreateRating(c echo.Context) error {
	request := new(dto.CreateRatingRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	userLogin := c.Get("userLogin")
	userId := userLogin.(jwt.MapClaims)["id"].(float64)

	star, _ := strconv.Atoi(c.FormValue("star"))

	rating := models.Rating{
		Star:    star,
		UserID:  int(userId),
		MovieID: request.MovieID,
	}

	data, err := h.RatingRepository.CreateRating(rating)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	data, _ = h.RatingRepository.GetRating(data.ID)

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: convertRatingResponse(data)})
}

// function delete rating
func (h *handlerRating) DeleteRating(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	rating, err := h.RatingRepository.GetRating(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.RatingRepository.DeleteRating(rating, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

// convert rating
func convertRatingResponse(rating models.Rating) models.RatingResponse {
	var result models.RatingResponse
	result.ID = rating.ID
	result.Star = rating.Star
	result.MovieID = rating.MovieID
	result.Movie = rating.Movie
	result.UserID = rating.UserID
	result.User = rating.User

	return result
}

// function convert multiple rating
func ConvertMultipleRatingResponse(ratings []models.Rating) []models.RatingResponse {
	var result []models.RatingResponse

	for _, rating := range ratings {
		ratings := models.RatingResponse{
			ID:      rating.ID,
			Star:    rating.Star,
			MovieID: rating.MovieID,
			Movie:   rating.Movie,
			UserID:  rating.UserID,
			User:    rating.User,
		}

		result = append(result, ratings)
	}

	return result
}
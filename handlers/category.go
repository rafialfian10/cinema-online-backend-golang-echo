package handlers

import (
	dto "cinemaonline/dto"
	"cinemaonline/models"
	"cinemaonline/repositories"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type handlerCategory struct {
	CategoryRepository repositories.CategoryRepository
}

func HandlerCategory(CategoryRepository repositories.CategoryRepository) *handlerCategory {
	return &handlerCategory{CategoryRepository}
}

// function get all category
func (h *handlerCategory) FindCategories(c echo.Context) error {
	categories, err := h.CategoryRepository.FindCategories()
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: categories})
}

// function get by id category
func (h *handlerCategory) GetCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	category, err := h.CategoryRepository.GetCategory(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: category})
}

// function create category
func (h *handlerCategory) CreateCategory(c echo.Context) error {
	request := new(dto.CreateCategoryRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	category := models.Category{
		Name: request.Name,
	}

	data, err := h.CategoryRepository.CreateCategory(category)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

// function update category
func (h *handlerCategory) UpdateCategory(c echo.Context) error {
	request := new(dto.UpdateCategoryRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	id, _ := strconv.Atoi(c.Param("id"))

	user, err := h.CategoryRepository.GetCategory(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	data, err := h.CategoryRepository.UpdateCategory(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

// function delete category
func (h *handlerCategory) DeleteCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	category, err := h.CategoryRepository.GetCategory(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	data, err := h.CategoryRepository.DeleteCategory(category, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

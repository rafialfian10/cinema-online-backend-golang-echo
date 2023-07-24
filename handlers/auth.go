package handlers

import (
	dto "cinemaonline/dto"
	"cinemaonline/models"
	"cinemaonline/pkg/bcrypt"
	jwtToken "cinemaonline/pkg/jwt"
	"cinemaonline/repositories"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// Handle struct
type handlerAuth struct {
	AuthRepository repositories.AuthRepository
}

// Handle Authentication
func HandlerAuth(AuthRepository repositories.AuthRepository) *handlerAuth {
	return &handlerAuth{AuthRepository}
}

// function Register
func (h *handlerAuth) Register(c echo.Context) error {
	request := new(dto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	password, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: password,
		Role:     "user",
	}

	data, err := h.AuthRepository.Register(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

// function register admin
func (h *handlerAuth) RegisterAdmin(c echo.Context) error {
	request := new(dto.RegisterRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	hashedPassword, err := bcrypt.HashingPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user := models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: hashedPassword,
		Role:     "admin",
	}

	data, err := h.AuthRepository.Register(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Status: http.StatusInternalServerError, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: data})
}

// function Login
func (h *handlerAuth) Login(c echo.Context) error {
	request := new(dto.LoginRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	user := models.User{
		Email:    request.Email,
		Password: request.Password,
	}

	// Check email
	user, err := h.AuthRepository.Login(user.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()})
	}

	// Check password
	isValid := bcrypt.CheckPasswordHash(request.Password, user.Password)
	if !isValid {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Status: http.StatusBadRequest, Message: "wrong email or password"})
	}

	//generate token
	claims := jwt.MapClaims{}
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 2 hours expired

	token, errGenerateToken := jwtToken.GenerateToken(&claims)
	if errGenerateToken != nil {
		log.Println(errGenerateToken)
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	loginResponse := dto.LoginResponse{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Token:    token,
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Status: http.StatusOK, Data: loginResponse})
}

// function check auth
func (h *handlerAuth) CheckAuth(c echo.Context) error {
	userInfo := c.Get("userInfo").(*jwt.Token)
	claims := userInfo.Claims.(jwt.MapClaims)
	userId := int(claims["id"].(float64))

	// Check User by Id
	user, err := h.AuthRepository.Getuser(userId)
	if err != nil {
		response := dto.ErrorResult{Status: http.StatusBadRequest, Message: err.Error()}
		return c.JSON(http.StatusBadRequest, response)
	}

	CheckAuthResponse := dto.CheckAuth{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	response := dto.SuccessResult{Status: http.StatusOK, Data: CheckAuthResponse}
	return c.JSON(http.StatusOK, response)
}

package middleware

import (
	dto "cinemaonline/dto"
	jwtToken "cinemaonline/pkg/jwt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ResultAdmin struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func AuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Status: http.StatusBadRequest, Message: "unauthorized"})
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, ResultAdmin{Status: http.StatusUnauthorized, Message: "unauthorized"})
		}

		role := claims["role"].(string)
		if role != "admin" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Status: http.StatusUnauthorized, Message: "unauthorized, you're not admin !"})
		}

		c.Set("userLogin", claims)
		return next(c)
	}
}

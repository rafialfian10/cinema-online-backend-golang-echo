package routes

import (
	"cinemaonline/handlers"
	"cinemaonline/pkg/middleware"
	"cinemaonline/pkg/mysql"
	"cinemaonline/repositories"

	"github.com/labstack/echo/v4"
)

func PremiumRoutes(e *echo.Group) {
	premiumRepository := repositories.RepositoryPremium(mysql.DB)
	h := handlers.HandlerPremium(premiumRepository)

	e.GET("/premiums", h.FindPremiums)
	e.GET("/premium/:id", h.GetPremium)
	e.POST("/premium", middleware.Auth(h.CreatePremium))
	e.POST("/notification_premium", h.NotificationPremium)
	e.PATCH("/premium/:id", middleware.Auth(h.UpdatePremiumByAdmin))
	e.DELETE("/premium/:id", middleware.Auth(h.DeletePremium))
}

package routes

import (
	"cinemaonline/handlers"
	"cinemaonline/pkg/middleware"
	"cinemaonline/pkg/mysql"
	"cinemaonline/repositories"

	"github.com/labstack/echo/v4"
)

func PremiRoutes(e *echo.Group) {
	premiRepository := repositories.RepositoryPremi(mysql.DB)
	h := handlers.HandlerPremi(premiRepository)

	e.GET("/premis", h.FindPremis)
	e.GET("/premi/:id", h.GetPremi)
	e.POST("/notification_premi", h.NotificationPremi)
	e.PATCH("/premi_user/:id", middleware.Auth(h.UpdatePremiumByUser))
	e.PATCH("/premi_admin/:id", middleware.Auth(h.UpdatePremiByAdmin))
	e.DELETE("/premi/:id", middleware.Auth(h.DeletePremi))
}

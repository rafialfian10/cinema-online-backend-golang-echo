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

	e.GET("/premis", middleware.Auth(h.FindPremis))
	e.GET("/premi/:id", middleware.Auth(h.GetPremi))
	e.POST("/notification_transaction_premi", h.NotificationPremi)
	e.PATCH("/premi_user/:id", middleware.Auth(h.UpdatePremiByUser))
	e.PATCH("/premi_admin/:id", middleware.Auth(h.UpdatePremiByAdmin))
	e.PATCH("/premi_expired/:id", h.UpdatePremiExpired)
	e.DELETE("/premi/:id", middleware.Auth(h.DeletePremi))
}

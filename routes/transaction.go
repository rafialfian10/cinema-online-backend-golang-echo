package routes

import (
	"cinemaonline/handlers"
	"cinemaonline/pkg/middleware"
	"cinemaonline/pkg/mysql"
	"cinemaonline/repositories"

	"github.com/labstack/echo/v4"
)

func TransactionRoutes(e *echo.Group) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	e.GET("/transactions_by_user", middleware.Auth(h.FindTransactionsByUser))
	e.GET("/transactions", middleware.Auth(h.FindTransactions))
	e.GET("/transaction/:id", middleware.Auth(h.GetTransaction))
	e.POST("/transaction", middleware.Auth(h.CreateTransaction))
	// e.POST("/notification_movie", h.Notification)
	e.PATCH("/transaction/:id", middleware.Auth(h.UpdateTransactionByAdmin))
	e.DELETE("/transaction/:id", middleware.Auth(h.DeleteTransaction))
}

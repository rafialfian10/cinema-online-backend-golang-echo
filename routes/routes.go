package routes

import "github.com/labstack/echo/v4"

func RouteInit(e *echo.Group) {
	AuthRoutes(e)
	UserRoutes(e)
	MovieRoutes(e)
	CategoryRoutes(e)
	TransactionRoutes(e)
	PremiRoutes(e)
}

package route

import (
	handler "github.com/HasanNugroho/golang-starter/internal/handler/auth"
	"github.com/labstack/echo/v4"
)

func NewAuthRoute(router *echo.Group, handler *handler.AuthHandler) {
	route := router.Group("/v1/auth")
	{
		// route.Use(middleware.AuthMiddleware(app))
		route.POST("/login", handler.Login)
		route.POST("/refresh", handler.RefreshToken)

	}
}

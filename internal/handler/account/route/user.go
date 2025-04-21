package route

import (
	handler "github.com/HasanNugroho/golang-starter/internal/handler/account"
	"github.com/HasanNugroho/golang-starter/internal/middleware"
	"github.com/labstack/echo/v4"
)

func NewUserRoute(router *echo.Group, handler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	userRoutes := router.Group("/v1/users")
	{
		userRoutes.Use(authMiddleware.AuthRequired())

		userRoutes.POST("", handler.Create)
		userRoutes.GET("/", handler.FindAll)
		userRoutes.GET("/:id", handler.FindById)
		userRoutes.GET("/me", handler.GetCurrentUser)
		userRoutes.PUT("/:id", handler.Update)
		userRoutes.DELETE("/:id", handler.Delete)

	}
}

package route

import (
	"github.com/HasanNugroho/golang-starter/internal/handler"
	"github.com/labstack/echo/v4"
)

func NewUserRoute(router *echo.Group, handler *handler.UserHandler) {
	userRoutes := router.Group("/v1/users")
	// userRoutes.Use(middleware.AuthMiddleware(app))
	{
		userRoutes.POST("", handler.Create)
		userRoutes.GET("/", handler.FindAll)
		userRoutes.GET("/:id", handler.FindById)
		userRoutes.PUT("/:id", handler.Update)
		userRoutes.DELETE("/:id", handler.Delete)

	}
}

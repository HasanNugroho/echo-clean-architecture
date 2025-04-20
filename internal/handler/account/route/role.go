package route

import (
	handler "github.com/HasanNugroho/golang-starter/internal/handler/account"
	"github.com/HasanNugroho/golang-starter/internal/middleware"
	"github.com/labstack/echo/v4"
)

func NewRoleRoute(router *echo.Group, handler *handler.RoleHandler, authMiddleware *middleware.AuthMiddleware) {
	route := router.Group("/v1/roles")
	route.Use(authMiddleware.AuthRequired())
	{
		// route.Use(middleware.AuthMiddleware(app))
		route.POST("", handler.Create)
		route.GET("", handler.FindAll)
		route.GET("/:id", handler.FindById)
		route.PUT("/:id", handler.Update)
		route.DELETE("/:id", handler.Delete)
		route.POST("/assign", handler.AssignUser)
		route.POST("/unassign", handler.UnAssignUser)

	}
}

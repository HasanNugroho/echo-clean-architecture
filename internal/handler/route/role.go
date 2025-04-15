package route

import (
	"github.com/HasanNugroho/golang-starter/internal/handler"
	"github.com/labstack/echo/v4"
)

func NewRoleRoute(router *echo.Group, handler *handler.RoleHandler) {
	route := router.Group("/v1/roles")
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

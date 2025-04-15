package main

import (
	"fmt"

	"github.com/HasanNugroho/golang-starter/cmd/docs"
	"github.com/HasanNugroho/golang-starter/internal"
	"github.com/HasanNugroho/golang-starter/internal/configs"
	"github.com/HasanNugroho/golang-starter/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title       Starter Golang API
// @version     1.0
// @description This is a sample server.

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:7000
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to load config: %v", err)
	}

	router := echo.New()
	internal.Init(config, router)

	// Middleware
	router.Use(middleware.ErrorHandler())

	// Swagger Setup
	loadSwagger(router, config)

	if err := router.Start(fmt.Sprintf(":%s", config.Server.ServerPort)); err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

}

func loadSwagger(r *echo.Echo, cfg *configs.Config) {
	docs.SwaggerInfo.Title = cfg.AppName
	docs.SwaggerInfo.Version = cfg.Version
	docs.SwaggerInfo.Description = cfg.AppName
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.ServerHost, cfg.Server.ServerPort)
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Swagger endpoint
	r.GET("/swagger/*", echoSwagger.WrapHandler)
}

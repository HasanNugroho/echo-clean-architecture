package internal

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HasanNugroho/golang-starter/internal/app"
	"github.com/HasanNugroho/golang-starter/internal/configs"
	accountHandler "github.com/HasanNugroho/golang-starter/internal/handler/account"
	accountRoute "github.com/HasanNugroho/golang-starter/internal/handler/account/route"
	authHandler "github.com/HasanNugroho/golang-starter/internal/handler/auth"
	authRoute "github.com/HasanNugroho/golang-starter/internal/handler/auth/route"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/middleware"
	"github.com/labstack/echo/v4"
)

func Init(config *configs.Config, router *echo.Echo) {
	logger := configs.InitLogger(config)

	// Inisialisasi MongoDB
	mongoDB, err := config.Database.InitMongo(logger)
	if err != nil {
		logger.Fatal().Msg(err.Error())
		panic(1)
	}

	// Inisialisasi Redis
	redisClient, err := config.Redis.InitRedis()
	if err != nil {
		logger.Fatal().Msg(err.Error())
		panic(1)
	}

	helper.SetJWTHelper(config.Security.JWTSecretKey, time.Duration(config.Security.JWTExpired)*time.Minute, time.Duration(config.Security.JWTRefreshTokenExpired)*time.Hour, redisClient)

	container, err := app.BuildContainer(config, mongoDB, logger)
	if err != nil {
		logger.Fatal().Msg(err.Error())
		panic(1)
	}

	apiGroup := router.Group("/api")
	authMiddleware := container.Get("authMiddleware").(*middleware.AuthMiddleware)

	roleHandler := container.Get("roleHandler").(*accountHandler.RoleHandler)
	userHandler := container.Get("userHandler").(*accountHandler.UserHandler)
	authHandler := container.Get("authHandler").(*authHandler.AuthHandler)

	// Daftarkan route
	accountRoute.NewRoleRoute(apiGroup, roleHandler, authMiddleware)
	accountRoute.NewUserRoute(apiGroup, userHandler, authMiddleware)
	authRoute.NewAuthRoute(apiGroup, authHandler)

	// Siapkan fungsi shutdown untuk melakukan cleanup (misal: shutdown Redis dan container)
	shutdownFunc := func() {
		configs.ShutdownRedis(redisClient)
		container.Delete()
	}

	// Jalankan goroutine untuk menangani sinyal terminasi
	go func() {
		terminateSignals := make(chan os.Signal, 1)
		signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

		sig := <-terminateSignals
		logger.Info().Msgf("Received stop signal: %v. Shutting down...", sig)
		shutdownFunc()
		os.Exit(0)
	}()
}

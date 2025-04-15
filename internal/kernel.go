package internal

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/HasanNugroho/golang-starter/internal/app"
	"github.com/HasanNugroho/golang-starter/internal/configs"
	"github.com/HasanNugroho/golang-starter/internal/handler"
	"github.com/HasanNugroho/golang-starter/internal/handler/route"
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

	container, err := app.BuildContainer(config, mongoDB, logger)
	if err != nil {
		logger.Fatal().Msg(err.Error())
		panic(1)
	}
	// Penting: jangan defer container.Delete() di sini!
	// Karena kita ingin container tetap aktif selama aplikasi berjalan, lalu dihapus saat shutdown.

	apiGroup := router.Group("/api")
	roleHandler := container.Get("roleHandler").(*handler.RoleHandler)
	userHandler := container.Get("userHandler").(*handler.UserHandler)

	// Daftarkan route
	route.NewRoleRoute(apiGroup, roleHandler)
	route.NewUserRoute(apiGroup, userHandler)

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

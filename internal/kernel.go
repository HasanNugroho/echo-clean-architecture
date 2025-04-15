package internal

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HasanNugroho/golang-starter/internal/configs"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func Init(router *echo.Echo) {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Msg("❌ Failed to initialize config: " + err.Error())
		panic(1)
	}

	logger := configs.InitLogger(config)

	mongoDB, err := config.Database.InitMongo(logger)
	if err != nil {
		logger.Fatal().Msg(err.Error())
		panic(1)
	}

	redisClient, err := config.Redis.InitRedis()
	if err != nil {
		logger.Fatal().Msg(err.Error())
		panic(1)
	}

	// init apps
	InitializeApp(router, mongoDB, logger)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			logger.Info().Msgf("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			configs.ShutdownRedis(redisClient)
			stop = true
		}
	}

	time.Sleep(5 * time.Second)
}

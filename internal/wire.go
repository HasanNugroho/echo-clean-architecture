package internal

import (
	"github.com/HasanNugroho/golang-starter/internal/repository"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	// wire.Bind(new(repository.IUserRepository), new(*repository.UserRepository)),
	// users.NewUserService, // Pastikan `UserService` di dalam package `users`
	// wire.Bind(new(users.IUserService), new(*users.UserService)), // Bind ke IUserService
	// users.NewUserHandler,
)

func InitializeApp(e *echo.Echo, db *mongo.Database, logger *zerolog.Logger) error {
	wire.Build(
		userSet,
	)
	return nil
}

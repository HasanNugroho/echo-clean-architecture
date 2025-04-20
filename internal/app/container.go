package app

import (
	"github.com/HasanNugroho/golang-starter/internal/configs"
	accounthandler "github.com/HasanNugroho/golang-starter/internal/handler/account"
	authhandler "github.com/HasanNugroho/golang-starter/internal/handler/auth"
	"github.com/HasanNugroho/golang-starter/internal/middleware"
	accountrepository "github.com/HasanNugroho/golang-starter/internal/repository/account"
	accountservice "github.com/HasanNugroho/golang-starter/internal/service/account"
	authservice "github.com/HasanNugroho/golang-starter/internal/service/auth"
	"github.com/rs/zerolog"
	"github.com/sarulabs/di/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// BuildContainer menginisialisasi dependency injection container untuk fitur Role dan User
func BuildContainer(cfg *configs.Config, mongoDB *mongo.Database, logger *zerolog.Logger) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	// Register logger (dipakai di semua layer)
	builder.Add(di.Def{
		Name: "logger",
		Build: func(ctn di.Container) (interface{}, error) {
			return logger, nil
		},
	})

	// Register MongoDB instance
	builder.Add(di.Def{
		Name: "mongoDB",
		Build: func(ctn di.Container) (interface{}, error) {
			return mongoDB, nil
		},
	})

	// --- ROLE FEATURE ---

	// RoleRepository
	builder.Add(di.Def{
		Name: "roleRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			mongoDB := ctn.Get("mongoDB").(*mongo.Database)
			log := ctn.Get("logger").(*zerolog.Logger)
			return accountrepository.NewRoleRepository(mongoDB, log), nil
		},
	})

	// RoleService
	builder.Add(di.Def{
		Name: "roleService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("roleRepository").(*accountrepository.RoleRepository)
			log := ctn.Get("logger").(*zerolog.Logger)
			roleService, err := accountservice.NewRoleService(repo, log)
			if err != nil {
				return nil, err
			}
			return roleService, nil
		},
	})

	// RoleHandler
	builder.Add(di.Def{
		Name: "roleHandler",
		Build: func(ctn di.Container) (interface{}, error) {
			roleSvc := ctn.Get("roleService").(*accountservice.RoleService)
			return accounthandler.NewRoleHandler(roleSvc), nil
		},
	})

	// --- USER FEATURE ---

	// UserRepository
	builder.Add(di.Def{
		Name: "userRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			mongoDB := ctn.Get("mongoDB").(*mongo.Database)
			log := ctn.Get("logger").(*zerolog.Logger)
			return accountrepository.NewUserRepository(mongoDB, log), nil
		},
	})

	// UserService
	builder.Add(di.Def{
		Name: "userService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("userRepository").(accountrepository.IUserRepository)
			rolerepo := ctn.Get("roleRepository").(*accountrepository.RoleRepository)
			log := ctn.Get("logger").(*zerolog.Logger)
			userService := accountservice.NewUserService(repo, rolerepo, log)
			return userService, nil
		},
	})

	// UserHandler
	builder.Add(di.Def{
		Name: "userHandler",
		Build: func(ctn di.Container) (interface{}, error) {
			userSvc := ctn.Get("userService").(accountservice.IUserService)
			return accounthandler.NewUserHandler(userSvc), nil
		},
	})

	// --- AUTH FEATURE ---
	builder.Add(di.Def{
		Name: "authService",
		Build: func(ctn di.Container) (interface{}, error) {
			log := ctn.Get("logger").(*zerolog.Logger)
			userSvc := ctn.Get("userService").(accountservice.IUserService)
			authService := authservice.NewAuthService(userSvc, log, cfg)
			return authService, nil
		},
	})

	// AuthHandler
	builder.Add(di.Def{
		Name: "authHandler",
		Build: func(ctn di.Container) (interface{}, error) {
			authSvc := ctn.Get("authService").(authservice.IAuthService)
			return authhandler.NewAuthHandler(authSvc), nil
		},
	})

	builder.Add(di.Def{
		Name: "authMiddleware",
		Build: func(ctn di.Container) (interface{}, error) {
			userSvc := ctn.Get("userService").(accountservice.IUserService)
			log := ctn.Get("logger").(*zerolog.Logger)

			return middleware.NewAuthMiddleware(log, userSvc), nil
		},
	})

	return builder.Build(), nil
}

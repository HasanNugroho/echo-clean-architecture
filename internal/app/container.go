package app

import (
	"github.com/HasanNugroho/golang-starter/internal/configs"
	"github.com/HasanNugroho/golang-starter/internal/handler"
	"github.com/HasanNugroho/golang-starter/internal/repository"
	"github.com/HasanNugroho/golang-starter/internal/service"
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
			return repository.NewRoleRepository(mongoDB, log), nil
		},
	})

	// RoleService
	builder.Add(di.Def{
		Name: "roleService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("roleRepository").(*repository.RoleRepository)
			log := ctn.Get("logger").(*zerolog.Logger)
			// NewRoleService memuat data permission dari YAML melalui helper
			roleService, err := service.NewRoleService(repo, log)
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
			roleSvc := ctn.Get("roleService").(*service.RoleService)
			return handler.NewRoleHandler(roleSvc), nil
		},
	})

	// --- USER FEATURE ---

	// UserRepository
	builder.Add(di.Def{
		Name: "userRepository",
		Build: func(ctn di.Container) (interface{}, error) {
			mongoDB := ctn.Get("mongoDB").(*mongo.Database)
			log := ctn.Get("logger").(*zerolog.Logger)
			return repository.NewUserRepository(mongoDB, log), nil
		},
	})

	// UserService
	builder.Add(di.Def{
		Name: "userService",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("userRepository").(repository.IUserRepository)
			log := ctn.Get("logger").(*zerolog.Logger)
			userService := service.NewUserService(repo, log)
			return userService, nil
		},
	})

	// UserHandler
	builder.Add(di.Def{
		Name: "userHandler",
		Build: func(ctn di.Container) (interface{}, error) {
			userSvc := ctn.Get("userService").(service.IUserService)
			return handler.NewUserHandler(userSvc), nil
		},
	})

	return builder.Build(), nil
}

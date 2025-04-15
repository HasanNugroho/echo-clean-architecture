package repository

import (
	"context"

	"github.com/HasanNugroho/golang-starter/internal/model"
)

type (
	IUserRepository interface {
		Create(ctx context.Context, user *model.User) error
		FindByEmail(ctx context.Context, email string) (*model.User, error)
		FindById(ctx context.Context, id string) (*model.User, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.User, int, error)
		Update(ctx context.Context, id string, user *model.User) error
		Delete(ctx context.Context, id string) error
	}

	IRoleRepository interface {
		Create(ctx context.Context, role *model.Role) error
		FindById(ctx context.Context, id string) (*model.Role, error)
		FindManyByID(ctx context.Context, ids []string) (*[]model.Role, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.Role, int, error)
		Update(ctx context.Context, id string, role *model.Role) error
		Delete(ctx context.Context, id string) error
		AssignUser(ctx context.Context, userId string, roleId string) error
		UnassignUser(ctx context.Context, userId string, roleId string) error
	}
)

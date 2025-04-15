package service

import (
	"context"

	"github.com/HasanNugroho/golang-starter/internal/model"
)

type (
	IUserService interface {
		Create(ctx context.Context, user *model.CreateUserRequest) error
		FindById(ctx context.Context, id string) (*model.UserResponse, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.UserResponse, int64, error)
		Update(ctx context.Context, id string, user *model.UpdateUserRequest) error
		Delete(ctx context.Context, id string) error
	}

	IRoleService interface {
		Create(ctx context.Context, user *model.CreateRoleRequest) error
		FindById(ctx context.Context, id string) (*model.Role, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.Role, int64, error)
		Update(ctx context.Context, id string, role *model.UpdateRoleRequest) error
		Delete(ctx context.Context, id string) error
		AssignUser(ctx context.Context, payload *model.AssignRoleModel) error
		UnassignUser(ctx context.Context, payload *model.AssignRoleModel) error
	}
)

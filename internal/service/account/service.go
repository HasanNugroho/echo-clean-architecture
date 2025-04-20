package account

import (
	"context"

	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/model/account"
)

type (
	IUserService interface {
		Create(ctx context.Context, user *account.CreateUserRequest) error
		FindById(ctx context.Context, id string) (*account.User, error)
		FindByEmail(ctx context.Context, email string) (*account.User, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]account.UserResponse, int64, error)
		Update(ctx context.Context, id string, user *account.UpdateUserRequest) error
		Delete(ctx context.Context, id string) error
	}

	IRoleService interface {
		Create(ctx context.Context, user *account.CreateRoleRequest) error
		FindById(ctx context.Context, id string) (*account.Role, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]account.Role, int64, error)
		Update(ctx context.Context, id string, role *account.UpdateRoleRequest) error
		Delete(ctx context.Context, id string) error
		AssignUser(ctx context.Context, payload *account.AssignRoleModel) error
		UnassignUser(ctx context.Context, payload *account.AssignRoleModel) error
	}
)

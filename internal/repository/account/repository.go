package account

import (
	"context"

	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/model/account"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	IUserRepository interface {
		Create(ctx context.Context, user *account.User) error
		FindByEmail(ctx context.Context, email string) (*account.User, error)
		FindById(ctx context.Context, id string) (*account.User, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]account.User, int, error)
		Update(ctx context.Context, id string, user *account.User) error
		Delete(ctx context.Context, id string) error
	}

	IRoleRepository interface {
		Create(ctx context.Context, role *account.Role) error
		FindById(ctx context.Context, id string) (*account.Role, error)
		FindManyByID(ctx context.Context, ids []bson.ObjectID) (*[]account.Role, error)
		FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]account.Role, int, error)
		Update(ctx context.Context, id string, role *account.Role) error
		Delete(ctx context.Context, id string) error
		AssignUser(ctx context.Context, userId string, roleId string) error
		UnassignUser(ctx context.Context, userId string, roleId string) error
	}
)

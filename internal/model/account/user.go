package account

import (
	"time"

	"github.com/HasanNugroho/golang-starter/internal/helper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type (
	User struct {
		ID          bson.ObjectID   `bson:"_id,omitempty" json:"id"`
		Email       string          `bson:"email" json:"email"`
		Name        string          `bson:"name" json:"name"`
		Password    string          `bson:"password" json:"password"`
		Roles       []bson.ObjectID `bson:"roles" json:"roles"`
		RolesDetail *[]Role         `bson:"-"`
		CreatedAt   time.Time       `bson:"created_at,omitempty" json:"created_at,omitempty"`
		UpdatedAt   time.Time       `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	}
)

type (
	UserResponse struct {
		ID        string    `bson:"_id,omitempty" json:"id"`
		Email     string    `bson:"email" json:"email"`
		Name      string    `bson:"name" json:"name"`
		Roles     *[]Role   `bson:"roles" json:"roles"`
		CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
		UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	}

	CreateUserRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required,min=6"`
	}

	UpdateUserRequest struct {
		Email    string `json:"email" validate:"email"`
		Name     string `json:"name" validate:""`
		Password string `json:"password" validate:"min=6"`
	}
)

func (u *User) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}

func (u *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID.Hex(),
		Email:     u.Email,
		Name:      u.Name,
		Roles:     u.RolesDetail,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (r *UserResponse) ToUser() *User {
	id, _ := bson.ObjectIDFromHex(r.ID)
	return &User{
		ID:          id,
		Email:       r.Email,
		Name:        r.Name,
		RolesDetail: r.Roles,
		Password:    "",
	}
}

func (u *User) IsHasAccess(permissions []string) bool {
	permSet := make(map[string]struct{})

	for _, role := range *u.RolesDetail {
		for _, p := range role.Permissions {
			permSet[p] = struct{}{}
		}
	}

	if _, ok := permSet["manage:system"]; ok {
		return true
	}

	defaultPerms, err := helper.LoadStringListFromYAML("./internal/constant/data.yaml", "default_permission")
	if err != nil {
		return false
	}
	for p := range defaultPerms {
		permSet[p] = struct{}{}
	}

	for _, p := range permissions {
		if _, ok := permSet[p]; ok {
			return true
		}
	}

	return false
}

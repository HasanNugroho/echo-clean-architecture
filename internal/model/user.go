package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	User struct {
		ID        bson.ObjectID   `bson:"_id,omitempty" json:"id"`
		Email     string          `bson:"email" json:"email"`
		Name      string          `bson:"name" json:"name"`
		Password  string          `bson:"password" json:"password"`
		Roles     []bson.ObjectID `bson:"roles" json:"roles"`
		CreatedAt time.Time       `bson:"created_at,omitempty" json:"created_at,omitempty"`
		UpdatedAt time.Time       `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	}
)

type (
	UserResponse struct {
		ID        string          `bson:"_id,omitempty" json:"id"`
		Email     string          `bson:"email" json:"email"`
		Name      string          `bson:"name" json:"name"`
		Roles     []bson.ObjectID `bson:"roles" json:"roles"`
		CreatedAt time.Time       `bson:"created_at,omitempty" json:"created_at,omitempty"`
		UpdatedAt time.Time       `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
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

func (u *User) IsValid() bool {
	return u.Email != ""
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) SetRoles(roles []bson.ObjectID) {
	u.Roles = roles
}

func (u *User) SetCreatedAt(createdAt time.Time) {
	u.CreatedAt = createdAt
}

func (u *User) SetUpdatedAt(updatedAt time.Time) {
	u.UpdatedAt = updatedAt
}

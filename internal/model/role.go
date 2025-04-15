package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	Role struct {
		ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
		Name        string        `bson:"name" json:"name"`
		Permissions []string      `bson:"permissions" json:"permissions"`
		CreatedAt   time.Time     `bson:"created_at,omitempty" json:"created_at,omitempty"`
		UpdatedAt   time.Time     `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	}
)

type (
	RoleResponse struct {
		ID          bson.ObjectID `bson:"_id" json:"id"`
		Name        string        `bson:"name" json:"name"`
		Permissions []string      `bson:"permissions" json:"permission"`
	}

	CreateRoleRequest struct {
		Name        string   `json:"name" validate:"required"`
		Permissions []string `json:"permission" validate:"required"`
	}

	UpdateRoleRequest struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permission"`
	}

	AssignRoleModel struct {
		UserID string `json:"user_id"`
		RoleID string `json:"role_id"`
	}
)

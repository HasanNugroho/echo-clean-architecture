package service

import (
	"context"
	"time"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/repository"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService struct {
	repo   repository.IUserRepository
	logger *zerolog.Logger
}

func NewUserService(repo repository.IUserRepository, logger *zerolog.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (u *UserService) Create(ctx context.Context, user *model.CreateUserRequest) error {
	_, err := u.repo.FindByEmail(ctx, user.Email)

	if err == nil {
		return errs.BadRequest("email exist", err)
	}

	password, err := helper.HashPassword([]byte(user.Password))
	if err != nil {
		u.logger.Error().Err(err).Msg("failed to hash password")
		return err
	}

	payload := model.User{
		Email:     user.Email,
		Name:      user.Name,
		Roles:     []bson.ObjectID{},
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err = u.repo.Create(ctx, &payload); err != nil {
		u.logger.Error().Err(err).Fields(payload).Msg("failed to create data")
		return err
	}

	return nil
}

func (u *UserService) FindById(ctx context.Context, id string) (*model.UserResponse, error) {
	user, err := u.repo.FindById(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Str("userID", id).Msg("error from repo")
		return &model.UserResponse{}, err
	}

	response := model.UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Name:      user.Name,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return &response, nil
}

func (u *UserService) FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.UserResponse, int64, error) {
	users, totalItems, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		u.logger.Error().Err(err).
			Str("search", filter.Search).
			Int("page", filter.Page).
			Int("limit", filter.Limit).
			Msg("error from repo")

		return &[]model.UserResponse{}, 0, err
	}

	var usersResponse []model.UserResponse
	for _, user := range *users {
		usersResponse = append(usersResponse, model.UserResponse{
			ID:        user.ID.Hex(),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return &usersResponse, int64(totalItems), nil
}

func (u *UserService) Update(ctx context.Context, id string, user *model.UpdateUserRequest) error {
	existingUser, err := u.repo.FindById(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Str("user", id).Msg("failed to find user for update")
		return err
	}

	if user.Email != "" {
		existingUser.Email = user.Email
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}

	if user.Password != "" {
		hashedPassword, err := helper.HashPassword([]byte(user.Password))
		if err != nil {
			u.logger.Error().Err(err).Msg("failed to hash password")
			return err
		}
		existingUser.Password = hashedPassword
	}

	if err := u.repo.Update(ctx, id, existingUser); err != nil {
		u.logger.Error().Err(err).Fields(existingUser).Msg("failed to update data")
		return err
	}

	return nil
}

func (u *UserService) Delete(ctx context.Context, id string) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		u.logger.Error().Err(err).Str("user", id).Msg("failed to delete data")
	}
	return err
}

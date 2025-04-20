package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/HasanNugroho/golang-starter/internal/configs"
	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/helper"
	"github.com/HasanNugroho/golang-starter/internal/model/auth"
	"github.com/HasanNugroho/golang-starter/internal/service/account"
	"github.com/rs/zerolog"
)

type AuthService struct {
	userservice account.IUserService
	logger      *zerolog.Logger
	config      *configs.Config
}

func NewAuthService(userservice account.IUserService, logger *zerolog.Logger, config *configs.Config) *AuthService {
	return &AuthService{
		userservice: userservice,
		logger:      logger,
		config:      config,
	}
}

func (a *AuthService) Login(ctx context.Context, request auth.LoginRequest) (auth.AuthResponse, error) {
	user, err := a.userservice.FindByEmail(ctx, request.Email)
	if err != nil {
		fmt.Println(err)
		return auth.AuthResponse{}, errs.Unauthorized("Incorrect email or password", err)
	}

	if !user.VerifyPassword(request.Password) {
		return auth.AuthResponse{}, errs.Unauthorized("Incorrect email or password", errors.New("incorrect email or password"))

	}

	accessToken, err := helper.GenerateToken(user.ID.Hex())
	if err != nil {
		return auth.AuthResponse{}, errs.Unauthorized("failed to generate token", err)
	}

	return auth.AuthResponse{
		Token: accessToken,
		Data: map[string]string{
			"user_id": user.ID.Hex(),
			"email":   user.Email,
		},
	}, nil
}

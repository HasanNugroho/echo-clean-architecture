package auth

import (
	"context"

	"github.com/HasanNugroho/golang-starter/internal/model/auth"
)

type (
	IAuthService interface {
		Login(ctx context.Context, request auth.LoginRequest) (auth.AuthResponse, error)
		RefreshToken(ctx context.Context, request auth.RenewalTokenRequest) (auth.AuthResponse, error)
	}
)

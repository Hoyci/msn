package auth

import (
	"context"
	"msn/pkg/common/dto"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	Logout(ctx context.Context) error
	RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewTokenResponse, error)
}

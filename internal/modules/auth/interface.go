package auth

import (
	"context"
	"msn/internal/infra/jwt"
	"msn/pkg/common/dto"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	Logout(ctx context.Context) error
	// RenewAccessToken(ctx context.Context, refreshToken string) (*dto.RenewTokenResponse, error)
}

type TokenProvider interface {
	GenerateAccessToken(user dto.UserResponse) (string, *jwt.Claims, error)
	GenerateRefreshToken(user dto.UserResponse) (string, *jwt.Claims, error)
	VerifyRefreshToken(token string) (*jwt.Claims, error)
}

type SessionManager interface {
	CreateSession(ctx context.Context, userID, jti string) (*dto.SessionResponse, error)
	DeactivateSession(ctx context.Context, userID string) error
}

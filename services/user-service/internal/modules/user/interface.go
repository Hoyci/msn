package user

import (
	"context"
	"msn/pkg/common/dto"
	"msn/services/user-service/internal/infra/database/model"
)

type Repository interface {
	Insert(ctx context.Context, user model.User) error
	// Update(ctx context.Context, user model.User) error
	// GetByID(ctx context.Context, userId string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	// Delete(ctx context.Context, userId string) error
}

type Service interface {
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	// GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error)
}

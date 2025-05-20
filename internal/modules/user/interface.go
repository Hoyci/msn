package user

import (
	"context"
	"msn/internal/infra/database/model"
	"msn/pkg/common/dto"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) error
	// Update(ctx context.Context, user model.User) error
	GetByID(ctx context.Context, userId string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	// Delete(ctx context.Context, userId string) error
	GetUserRoleByName(ctx context.Context, name string) (*model.UserRole, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error)
}

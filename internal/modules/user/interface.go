package user

import (
	"context"
	"msn/pkg/common/dto"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	// Update(ctx context.Context, user model.User) error
	GetByID(ctx context.Context, userId string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	// Delete(ctx context.Context, userId string) error
	RoleExists(ctx context.Context, roleID string) (bool, error)
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error)
}

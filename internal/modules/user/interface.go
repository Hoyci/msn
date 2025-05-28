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
	GetEnrichedByEmail(ctx context.Context, email string) (*dto.EnrichedUserResponse, error)
	GetProfessionalUsers(ctx context.Context) ([]*dto.ProfessionalUserResponse, error)
	// Delete(ctx context.Context, userId string) error
}

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
	GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error)
	GetProfessionalUsers(ctx context.Context) ([]*dto.ProfessionalUserResponse, error)
}

package user

import (
	"context"
	"errors"
	"fmt"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/dbutil"

	"github.com/lib/pq"
)

type ServiceConfig struct {
	UserRepo Repository
}

type service struct {
	userRepo Repository
}

func NewService(c ServiceConfig) Service {
	return &service{
		userRepo: c.UserRepo,
	}
}

func (s service) CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to get user by email")
	} else if userRecord != nil {
		return nil, fault.NewConflict("e-mail already taken")
	}

	newUser, err := New(input.Name, input.Email, input.Password, input.ConfirmPassword, input.AvatarUrl)
	if err != nil {
		return nil, fault.NewUnprocessableEntity("failed to create user entity")
	}
	model := newUser.Model()

	if err = s.userRepo.Insert(ctx, model); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // 23505 is the code for unique constraint violation
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			return nil, fault.NewConflict(fmt.Sprintf("%s already taken", field))
		}
		return nil, fault.NewBadRequest("failed to insert user")
	}

	user := dto.UserResponse{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		AvatarURL: model.AvatarURL,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
	}

	return &user, nil
}

func (s service) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		CreatedAt: userRecord.CreatedAt,
		UpdatedAt: userRecord.UpdatedAt,
		DeletedAt: userRecord.DeletedAt,
	}

	return &user, nil
}

func (s service) GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error) {
	userRecord, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if userRecord == nil {
		return nil, fault.NewNotFound("user not found")
	}

	user := dto.UserResponse{
		ID:        userRecord.ID,
		Name:      userRecord.Name,
		Email:     userRecord.Email,
		AvatarURL: userRecord.AvatarURL,
		CreatedAt: userRecord.CreatedAt,
		UpdatedAt: userRecord.UpdatedAt,
		DeletedAt: userRecord.DeletedAt,
	}

	return &user, nil
}

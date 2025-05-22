package user

import (
	"context"
	"errors"
	"fmt"
	"msn/internal/infra/logging"
	"msn/internal/modules/category"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/common/valueobjects"
	"msn/pkg/utils/dbutil"

	"github.com/lib/pq"
)

type ServiceConfig struct {
	UserRepo     UserRepository
	CategoryRepo category.Repository
}

type service struct {
	userRepo     UserRepository
	categoryRepo category.Repository
}

func NewService(c ServiceConfig) UserService {
	return &service{
		userRepo:     c.UserRepo,
		categoryRepo: c.CategoryRepo,
	}
}

func (s service) CreateUser(ctx context.Context, input dto.CreateUser) (*dto.UserResponse, error) {
	logger := logging.FromContext(ctx)

	logger.DebugContext(ctx, "user_creation_attempt",
		"email", input.Email,
		"name", input.Name,
	)

	existing, err := s.userRepo.GetByEmail(context.Background(), input.Email)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to validate email")
	}
	if existing != nil {
		return nil, fault.NewConflict("email already taken")
	}

	exists, err := s.userRepo.RoleExists(context.Background(), input.UserRoleID)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to validate user role")
	}
	if !exists {
		return nil, fault.NewBadRequest("invalid user role")
	}

	password, err := valueobjects.NewPassword(input.Password)
	if err != nil {
		return nil, fault.NewBadRequest(err.Error())
	}

	user, err := New(
		input.Name,
		input.Email,
		password.Hash,
		input.UserRoleID,
		input.AvatarUrl,
		input.SubcategoryID,
	)
	if err != nil {
		logger.DebugContext(ctx, "invalid user entity", "error", err)
		return nil, fault.NewUnprocessableEntity(err.Error())
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		logger.ErrorContext(ctx, "db_error",
			"operation", "userRepo.Insert",
			"email", input.Email,
			"error", err,
		)

		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			field := dbutil.ExtractFieldFromDetail(pqErr.Detail)
			logger.WarnContext(ctx, "unique_constraint_violation",
				"field", field,
				"email", input.Email,
			)
			return nil, fault.NewConflict(fmt.Sprintf("%s already taken", field))
		}

		return nil, fault.NewBadRequest("failed to create user")
	}

	logger.InfoContext(ctx, "user_created",
		"user_id", user.ID,
		"email", user.Email,
	)

	return &dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		UserRoleID:    user.UserRoleID,
		SubcategoryID: user.SubcategoryID,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (s service) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if user == nil {
		return nil, fault.NewNotFound("user not found")
	}

	return &dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		UserRoleID:    user.UserRoleID,
		SubcategoryID: user.SubcategoryID,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (s service) GetUserByID(ctx context.Context, userId string) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve user")
	}
	if user == nil {
		return nil, fault.NewNotFound("user not found")
	}

	return &dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Email:         user.Email,
		UserRoleID:    user.UserRoleID,
		SubcategoryID: user.SubcategoryID,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
	}, nil
}

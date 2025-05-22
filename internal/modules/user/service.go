package user

import (
	"context"
	"errors"
	"fmt"
	"msn/internal/infra/logging"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/utils/dbutil"

	Categories "msn/internal/modules/categories"

	"github.com/lib/pq"
)

type ServiceConfig struct {
	UserRepo     UserRepository
	CategoryRepo Categories.Repository
}

type service struct {
	userRepo     UserRepository
	categoryRepo Categories.Repository
}

func NewUserService(c ServiceConfig) UserService {
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

	user, err := New(
		input.Name,
		input.Email,
		input.Password,
		input.ConfirmPassword,
		input.UserRoleID,
		input.AvatarUrl,
		input.SubcategoryID,
		s.userRepo,
		s.categoryRepo,
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
		"user_id", user.ID(),
		"email", user.Email(),
	)

	return &dto.UserResponse{
		ID:            user.ID(),
		Name:          user.Name(),
		Email:         user.Email(),
		UserRoleID:    user.UserRoleID(),
		SubcategoryID: user.SubcategoryID(),
		AvatarURL:     user.AvatarURL(),
		CreatedAt:     user.CreatedAt(),
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
		ID:            user.ID(),
		Name:          user.Name(),
		Email:         user.Email(),
		UserRoleID:    user.UserRoleID(),
		SubcategoryID: user.SubcategoryID(),
		AvatarURL:     user.AvatarURL(),
		CreatedAt:     user.CreatedAt(),
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
		ID:            user.ID(),
		Name:          user.Name(),
		Email:         user.Email(),
		UserRoleID:    user.UserRoleID(),
		SubcategoryID: user.SubcategoryID(),
		AvatarURL:     user.AvatarURL(),
		CreatedAt:     user.CreatedAt(),
	}, nil
}

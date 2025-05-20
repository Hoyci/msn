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

	userRecord, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		logger.ErrorContext(ctx, "db_error",
			"operation", "userRepo.GetByEmail",
			"email", input.Email,
			"error", err,
		)
		return nil, fault.NewBadRequest("failed to validate email")
	}

	if userRecord != nil {
		logger.DebugContext(ctx, "email_conflict",
			"email", input.Email,
		)
		return nil, fault.NewConflict("e-mail already taken")
	}

	userRole, err := s.userRepo.GetUserRoleByName(ctx, input.UserRole)
	if err != nil {
		logger.ErrorContext(ctx, "get_user_role_error",
			"operation", "userRepo.GetUserRoleByName",
			"email", input.UserRole,
			"error", err,
		)
		if userRole == nil {
			return nil, fault.NewBadRequest("user role not found")
		}
		return nil, fault.NewInternalServerError("failed to get user role")
	}

	if userRole.Name == "client" && input.SubcategoryID != nil {
		return nil, fault.NewUnprocessableEntity("clients cannot have subcategories")
	}

	if userRole.Name == "professional" && input.SubcategoryID == nil {
		return nil, fault.NewUnprocessableEntity("professionals must provide a subcategory")
	}

	if input.SubcategoryID != nil {
		subcategory, err := s.categoryRepo.GetSubcategoryByID(ctx, *input.SubcategoryID)
		if err != nil {
			logger.ErrorContext(ctx, "get_subcategory_by_id",
				"operation", "categoryRepo.GetSubcategoryByID",
				"id", input.SubcategoryID,
				"error", err,
			)
			return nil, fault.NewInternalServerError("failed to get subcategory")
		}
		if subcategory == nil {
			return nil, fault.NewBadRequest("subcategory not found")
		}
	}

	newUser, err := New(input.Name, input.Email, input.Password, input.ConfirmPassword, userRole.ID, input.AvatarUrl, input.SubcategoryID)
	if err != nil {
		logger.DebugContext(ctx, "invalid_user_entity",
			"email", input.Email,
			"error", err,
		)
		return nil, fault.NewUnprocessableEntity(err.Error())
	}

	model := newUser.Model()
	if err = s.userRepo.Insert(ctx, model); err != nil {
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
		"user_id", model.ID,
		"email", model.Email,
	)

	return &dto.UserResponse{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		UserRoleID:    model.UserRoleID,
		SubcategoryID: model.SubcategoryID,
		AvatarURL:     model.AvatarURL,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		DeletedAt:     model.DeletedAt,
	}, nil
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
		ID:            userRecord.ID,
		Name:          userRecord.Name,
		Email:         userRecord.Email,
		AvatarURL:     userRecord.AvatarURL,
		UserRoleID:    userRecord.UserRoleID,
		SubcategoryID: userRecord.SubcategoryID,
		CreatedAt:     userRecord.CreatedAt,
		UpdatedAt:     userRecord.UpdatedAt,
		DeletedAt:     userRecord.DeletedAt,
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
		ID:            userRecord.ID,
		Name:          userRecord.Name,
		Email:         userRecord.Email,
		UserRoleID:    userRecord.UserRoleID,
		SubcategoryID: userRecord.SubcategoryID,
		AvatarURL:     userRecord.AvatarURL,
		CreatedAt:     userRecord.CreatedAt,
		UpdatedAt:     userRecord.UpdatedAt,
		DeletedAt:     userRecord.DeletedAt,
	}

	return &user, nil
}

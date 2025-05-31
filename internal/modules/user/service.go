package user

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"msn/internal/config"
	"msn/internal/infra/logging"
	"msn/internal/infra/storage"
	"msn/internal/modules/category"
	"msn/internal/modules/role"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/common/valueobjects"
	"msn/pkg/utils/dbutil"

	"github.com/lib/pq"
)

type ServiceConfig struct {
	UserRepo      UserRepository
	CategoryRepo  category.Repository
	RoleRepo      role.Repository
	StorageClient *storage.StorageClient
}

type service struct {
	userRepo      UserRepository
	categoryRepo  category.Repository
	roleRepo      role.Repository
	storageClient *storage.StorageClient
}

func NewService(c ServiceConfig) UserService {
	return &service{
		userRepo:      c.UserRepo,
		categoryRepo:  c.CategoryRepo,
		roleRepo:      c.RoleRepo,
		storageClient: c.StorageClient,
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

	role, err := s.roleRepo.GetRoleByName(context.Background(), input.Role)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to validate user role")
	}

	password, err := valueobjects.NewPassword(input.Password)
	if err != nil {
		return nil, fault.NewBadRequest(err.Error())
	}

	user, err := New(
		input.Name,
		input.Email,
		password.Hash,
		role.ID,
		input.SubcategoryID,
	)
	if err != nil {
		logger.DebugContext(ctx, "invalid user entity", "error", err)
		return nil, fault.NewUnprocessableEntity(err.Error())
	}

	avatarUrl, err := s.UploadUserPicture(ctx, user.ID, input.FileHeader)
	if err != nil {
		return nil, fault.NewInternalServerError(err.Error())
	}

	user.AvatarURL = avatarUrl

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
		RoleID:        user.RoleID,
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
		RoleID:        user.RoleID,
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
		RoleID:        user.RoleID,
		SubcategoryID: user.SubcategoryID,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
	}, nil
}

func (s service) GetProfessionalUsers(ctx context.Context) ([]*dto.ProfessionalUserResponse, error) {
	professionals, err := s.userRepo.GetProfessionalUsers(ctx)
	if err != nil {
		return nil, fault.NewBadRequest("failed to retrieve professional users")
	}

	if professionals == nil {
		professionals = []*dto.ProfessionalUserResponse{}
	}

	return professionals, nil
}

func (s service) UploadUserPicture(ctx context.Context, userID string, fileHeader *multipart.FileHeader) (string, error) {
	objectName := fmt.Sprintf("user_%s_%s", "profile", userID)
	avatarKey, err := s.storageClient.UploadFile("user-profile", objectName, fileHeader)
	if err != nil {
		return "", err
	}
	avatarURL := fmt.Sprintf("%s/%s/%s", config.GetConfig().StorageURL, "user-profile", avatarKey)

	return avatarURL, nil
}

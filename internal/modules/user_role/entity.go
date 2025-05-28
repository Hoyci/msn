package userrole

import (
	"msn/internal/infra/database/models"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type UserRole struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func New(
	name string,
) (*UserRole, error) {
	userRole := UserRole{
		ID:        uid.New("user_role"),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
		DeletedAt: nil,
	}

	if err := userRole.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &userRole, nil
}

func NewFromModel(m models.UserRole) *UserRole {
	return &UserRole{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func (ur *UserRole) ToModel() models.UserRole {
	return models.UserRole{
		ID:        ur.ID,
		Name:      ur.Name,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
		DeletedAt: ur.DeletedAt,
	}
}

func (ur *UserRole) validate() error {
	if ur.Name == "" {
		return fault.NewBadRequest("role name is required")
	}

	return nil
}

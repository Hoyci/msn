package role

import (
	"msn/internal/infra/database/models"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type Role struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func New(
	name string,
) (*Role, error) {
	userRole := Role{
		ID:        uid.New("role"),
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

func NewFromModel(m models.Role) *Role {
	return &Role{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

func (ur *Role) ToModel() models.Role {
	return models.Role{
		ID:        ur.ID,
		Name:      ur.Name,
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
		DeletedAt: ur.DeletedAt,
	}
}

func (ur *Role) validate() error {
	if ur.Name == "" {
		return fault.NewBadRequest("role name is required")
	}

	return nil
}

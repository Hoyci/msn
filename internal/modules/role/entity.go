package role

import (
	"msn/internal/infra/database/models"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type Role struct {
	id        string
	name      string
	createdAt time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

func New(
	name string,
) (*Role, error) {
	role := Role{
		id:        uid.New("role"),
		name:      name,
		createdAt: time.Now(),
		updatedAt: nil,
		deletedAt: nil,
	}

	if err := role.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &role, nil
}

func NewFromModel(m models.Role) *Role {
	return &Role{
		id:        m.ID,
		name:      m.Name,
		createdAt: m.CreatedAt,
		updatedAt: m.UpdatedAt,
		deletedAt: m.DeletedAt,
	}
}

func (r *Role) ToModel() models.Role {
	return models.Role{
		ID:        r.ID(),
		Name:      r.Name(),
		CreatedAt: r.CreatedAt(),
		UpdatedAt: r.UpdatedAt(),
		DeletedAt: r.DeletedAt(),
	}
}

func (r *Role) validate() error {
	if r.Name() == "" {
		return fault.NewBadRequest("role name is required")
	}

	return nil
}

func FromID(id string) *Role {
	return &Role{id: id}
}

func (r *Role) ID() string {
	return r.id
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Role) UpdatedAt() *time.Time {
	return r.updatedAt
}

func (r *Role) DeletedAt() *time.Time {
	return r.deletedAt
}

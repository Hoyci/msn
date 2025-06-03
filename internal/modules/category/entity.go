package category

import (
	"msn/internal/infra/database/models"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type Category struct {
	id         string
	name       string
	created_at time.Time
	updated_at *time.Time
	deleted_at *time.Time
}

func NewFromModel(m models.Category) *Category {
	return &Category{
		id:         m.ID,
		name:       m.Name,
		created_at: m.CreatedAt,
		updated_at: m.UpdatedAt,
		deleted_at: m.DeletedAt,
	}
}

func New(name string) (*Category, error) {
	newCategory := Category{
		id:         uid.New("category"),
		name:       name,
		created_at: time.Now(),
		updated_at: nil,
		deleted_at: nil,
	}

	if err := newCategory.validate(); err != nil {
		return nil, fault.New(
			"failed to create category entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &newCategory, nil
}

func (c *Category) Model() models.Category {
	return models.Category{
		ID:        c.id,
		Name:      c.name,
		CreatedAt: c.created_at,
		UpdatedAt: c.updated_at,
		DeletedAt: c.deleted_at,
	}
}

func (c *Category) validate() error {
	if c.name == "" {
		return fault.New("category name is required")
	}
	return nil
}

func FromID(id string) *Category {
	return &Category{id: id}
}

func (c *Category) ID() string {
	return c.id
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) CreatedAt() time.Time {
	return c.created_at
}

func (c *Category) DeletedAt() *time.Time {
	return c.deleted_at
}

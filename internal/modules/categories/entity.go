package categories

import (
	"msn/internal/infra/database/model"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type category struct {
	id         string
	name       string
	created_at time.Time
	updated_at *time.Time
	deleted_at *time.Time
}

func NewFromModel(m model.Category) *category {
	return &category{
		id:         m.ID,
		name:       m.Name,
		created_at: m.CreatedAt,
		updated_at: m.UpdatedAt,
		deleted_at: m.DeletedAt,
	}
}

func New(name string) (*category, error) {
	c := category{
		id:         uid.New("category"),
		name:       name,
		created_at: time.Now(),
		updated_at: nil,
		deleted_at: nil,
	}

	if err := c.validate(); err != nil {
		return nil, fault.New(
			"failed to create category entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &c, nil
}

func (c *category) Model() model.Category {
	return model.Category{
		ID:        c.id,
		Name:      c.name,
		CreatedAt: c.created_at,
		UpdatedAt: c.updated_at,
		DeletedAt: c.deleted_at,
	}
}

func (c *category) validate() error {
	if c.name == "" {
		return fault.New("category name is required")
	}
	return nil
}

package subcategory

import (
	"msn/internal/infra/database/models"
	"msn/internal/modules/category"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type Subcategory struct {
	id        string
	name      string
	category  category.Category
	createdAt time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

func New(
	name,
	categoryID string,
) (*Subcategory, error) {
	subcategory := &Subcategory{
		id:        uid.New("subcategory"),
		name:      name,
		category:  *category.FromID(categoryID),
		createdAt: time.Now(),
		updatedAt: nil,
		deletedAt: nil,
	}

	if err := subcategory.validate(); err != nil {
		return nil, fault.New(
			"failed to create subcategory entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return subcategory, nil
}

func NewFromModel(m models.Subcategory) *Subcategory {
	return &Subcategory{
		id:        m.ID,
		name:      m.Name,
		category:  *category.FromID(m.CategoryID),
		createdAt: m.CreatedAt,
		updatedAt: m.UpdatedAt,
		deletedAt: m.DeletedAt,
	}
}

func (s *Subcategory) ToModel() models.Subcategory {
	return models.Subcategory{
		ID:         s.id,
		Name:       s.name,
		CategoryID: s.category.ID(),
		CreatedAt:  s.createdAt,
		UpdatedAt:  s.updatedAt,
		DeletedAt:  s.deletedAt,
	}
}

func (s *Subcategory) validate() error {
	if s.name == "" {
		return fault.NewBadRequest("name is required")
	}
	if s.category.ID() == "" {
		return fault.NewBadRequest("category_id is required")
	}

	return nil
}

func FromID(id string) *Subcategory {
	return &Subcategory{id: id}
}

func (s *Subcategory) ID() string {
	return s.id
}

func (s *Subcategory) Name() string {
	return s.name
}

func (s *Subcategory) Category() category.Category {
	return s.category
}

func (s *Subcategory) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Subcategory) DeletedAt() *time.Time {
	return s.deletedAt
}

package subcategory

import (
	"msn/internal/infra/database/models"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

type Subcategory struct {
	ID         string
	Name       string
	CategoryID string
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
}

func New(
	name,
	categoryID string,
) (*Subcategory, error) {
	subcategory := Subcategory{
		ID:         uid.New("subcategory"),
		Name:       name,
		CategoryID: categoryID,
		CreatedAt:  time.Now(),
		UpdatedAt:  nil,
		DeletedAt:  nil,
	}

	if err := subcategory.validate(); err != nil {
		return nil, fault.New(
			"failed to create subcategory entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &subcategory, nil
}

func NewFromModel(m models.Subcategory) *Subcategory {
	return &Subcategory{
		ID:         m.ID,
		Name:       m.Name,
		CategoryID: m.CategoryID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
	}
}

func (s *Subcategory) ToModel() models.Subcategory {
	return models.Subcategory{
		ID:         s.ID,
		Name:       s.Name,
		CategoryID: s.CategoryID,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		DeletedAt:  s.DeletedAt,
	}
}

func (s *Subcategory) validate() error {
	if s.Name == "" {
		return fault.NewBadRequest("name is required")
	}
	if s.CategoryID == "" {
		return fault.NewBadRequest("category_id is required")
	}

	return nil
}

package categories

import (
	"context"
	"msn/internal/infra/database/model"
	"msn/pkg/common/dto"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*dto.Category, error)
	GetWithSubs(ctx context.Context) ([]*dto.Category, error)
	GetSubcategoryByID(ctx context.Context, ID string) (*model.Subcategory, error)
	SubcategoryExists(ctx context.Context, subcategoryID string) (bool, error)
}

type Service interface {
	GetCategories(ctx context.Context, includeSubs bool) ([]*dto.Category, error)
}

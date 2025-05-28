package category

import (
	"context"
	"msn/pkg/common/dto"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]*dto.Category, error)
	GetCategoriesWithSubcategories(ctx context.Context) ([]*dto.Category, error)
}

type Service interface {
	GetCategories(ctx context.Context, includeSubs bool) ([]*dto.Category, error)
}

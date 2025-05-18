package categories

import (
	"context"
	"msn/internal/infra/database/model"
	"msn/pkg/common/dto"
)

type Repository interface {
	GetCategoriesWithUserCount(ctx context.Context) ([]*dto.CategoryWithUserCount, error)
	GetCategoriesWithSubcategories(ctx context.Context) ([]*dto.CategoryWithSubcategories, error)
	GetSubcategoryByID(ctx context.Context, ID string) (*model.Subcategory, error)
}

type Service interface {
	GetCategoriesWithUserCount(ctx context.Context) (*dto.CategoryResponse, error)
	GetCategoriesWithSubcategories(ctx context.Context) (*dto.CategoryWithSubResponse, error)
}

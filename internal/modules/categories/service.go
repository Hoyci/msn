package categories

import (
	"context"
	"msn/internal/infra/logging"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
)

type ServiceConfig struct {
	CategoriesRepo Repository
}

type service struct {
	categoriesRepo Repository
}

func NewService(c ServiceConfig) Service {
	return &service{
		categoriesRepo: c.CategoriesRepo,
	}
}

func (s *service) GetCategories(ctx context.Context, includeSubs bool) ([]*dto.Category, error) {
	logger := logging.FromContext(ctx)

	logger.DebugContext(
		ctx,
		"get_categories_attempt",
		"includeSubs", includeSubs,
	)

	if includeSubs {
		records, err := s.categoriesRepo.GetWithSubs(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "db_error",
				"operation", "categoriesRepo.GetWithSubs",
				"error", err,
			)

			return nil, fault.NewInternalServerError("failed to retrieve categories with subs")
		}
		return records, nil
	}

	records, err := s.categoriesRepo.GetCategories(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "db_error",
			"operation", "categoriesRepo.GetWithSubs",
			"error", err,
		)

		return nil, fault.NewInternalServerError("failed to retrieve categories")
	}
	return records, nil
}

//
// func (s service) GetCategoriesWithUserCount(ctx context.Context) (*dto.CategoryResponse,
// error) {
// 	logger := logging.FromContext(ctx)
//
// 	logger.DebugContext(ctx, "get_all_categories_attempt")
//
// 	records, err := s.categoriesRepo.GetCategoriesWithUserCount(ctx)
// 	if err != nil {
// 		logger.ErrorContext(ctx, "db_error",
// 			"operation", "categoriesRepo.GetCategories",
// 			"error", err,
// 		)
// 		return nil, fault.NewBadRequest("failed to retrieve categories")
// 	}
//
// 	var categories []*dto.CategoryWithUserCount
// 	for _, c := range records {
// 		categories = append(categories, &dto.CategoryWithUserCount{
// 			Category: dto.Category{
// 				ID:        c.Category.ID,
// 				Name:      c.Category.Name,
// 				Icon:      c.Category.Icon,
// 				CreatedAt: c.Category.CreatedAt,
// 			},
// 			SubcategoriesUserCount: c.SubcategoriesUserCount,
// 		})
// 	}
//
// 	logger.InfoContext(ctx, "categories_retrieved",
// 		"count", len(categories),
// 	)
//
// 	return &dto.CategoryResponse{
// 		Categories: categories,
// 	}, nil
// }
//
// func (s service) GetCategoriesWithSubcategories(ctx context.Context) (*dto.CategoryWithSubResponse, error) {
// 	logger := logging.FromContext(ctx)
//
// 	logger.DebugContext(ctx, "get_categories_with_subcategories_attempt")
//
// 	records, err := s.categoriesRepo.GetCategoriesWithSubcategories(ctx)
// 	if err != nil {
// 		logger.ErrorContext(ctx, "db_error",
// 			"operation", "categoriesRepo.GetCategoriesWithSubcategories",
// 			"error", err,
// 		)
// 		return nil, fault.NewBadRequest("failed to retrieve categories with subcategories")
// 	}
//
// 	var categories []*dto.CategoryWithSubcategories
// 	for _, c := range records {
// 		categories = append(categories, &dto.CategoryWithSubcategories{
// 			Category: dto.Category{
// 				ID:        c.Category.ID,
// 				Name:      c.Category.Name,
// 				Icon:      c.Category.Icon,
// 				CreatedAt: c.Category.CreatedAt,
// 				UpdatedAt: c.Category.UpdatedAt,
// 				DeletedAt: c.Category.DeletedAt,
// 			},
// 			Subcategories: c.Subcategories,
// 		})
// 	}
//
// 	logger.InfoContext(ctx, "categories_with_subcategories_retrieved",
// 		"count", len(categories),
// 	)
//
// 	return &dto.CategoryWithSubResponse{
// 		Categories: categories,
// 	}, nil
// }

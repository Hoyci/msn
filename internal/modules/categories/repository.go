package categories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"msn/internal/infra/database/model"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

func (r repo) GetCategoriesWithUserCount(ctx context.Context) ([]*dto.CategoryWithUserCount, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var categories []*dto.CategoryWithUserCount
	err := r.db.SelectContext(
		ctx,
		&categories,
		`SELECT 
			c.*,
  		COUNT(u.id) AS subcategories_user_count
		FROM categories c
		LEFT JOIN subcategories s ON s.category_id = c.id
		LEFT JOIN users u ON u.subcategory_id = s.id
		GROUP BY c.id
		ORDER BY subcategories_user_count DESC;`,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve categories", fault.WithError(err))
	}

	return categories, nil
}

func (r repo) GetCategoriesWithSubcategories(ctx context.Context) ([]*dto.CategoryWithSubcategories, error) {
	rows, err := r.db.QueryxContext(ctx, `
		SELECT
			c.id,
			c.name,
			c.icon,
			c.created_at,
			c.updated_at,
			c.deleted_at,
			COALESCE(
				JSON_AGG(
					json_build_object('id', s.id, 'name', s.name)
				) FILTER (WHERE s.deleted_at IS NULL),
				'[]'
			) AS subcategories
		FROM categories c
		LEFT JOIN subcategories s ON s.category_id = c.id
		WHERE c.deleted_at IS NULL
		GROUP BY c.id, c.name, c.icon, c.created_at, c.updated_at, c.deleted_at
		ORDER BY c.name ASC;`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*dto.CategoryWithSubcategories
	for rows.Next() {
		var cat dto.CategoryWithSubcategories
		var raw []byte
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.Icon,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.DeletedAt,
			&raw,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(raw, &cat.Subcategories); err != nil {
			return nil, err
		}

		result = append(result, &cat)
	}
	return result, nil
}

func (r repo) GetSubcategoryByID(ctx context.Context, ID string) (*model.Subcategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var subcategory model.Subcategory
	err := r.db.GetContext(ctx, &subcategory, "SELECT * FROM subcategories WHERE id = $1", ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve subcategory by id", fault.WithError(err))
	}

	return &subcategory, nil
}

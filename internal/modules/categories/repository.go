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

func (r repo) GetSubcategoryByID(ctx context.Context, categoryID string) (*model.Subcategory, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var subcategory model.Subcategory
	err := r.db.GetContext(ctx, &subcategory, "SELECT * FROM subcategories WHERE id = $1", categoryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve subcategory by id", fault.WithError(err))
	}

	return &subcategory, nil
}

func (r repo) GetCategories(ctx context.Context) ([]*dto.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var categories []*dto.Category
	err := r.db.SelectContext(
		ctx,
		&categories,
		`
		SELECT 
			c.id,
			c.name,
			c.icon,
			COUNT(DISTINCT u.id) as users
		FROM categories c
		LEFT JOIN subcategories s ON s.category_id = c.id
		LEFT JOIN users u ON u.subcategory_id = s.id
		GROUP BY c.id
		ORDER BY users DESC
		`,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve categories", fault.WithError(err))
	}

	return categories, nil
}

func (r repo) GetWithSubs(ctx context.Context) ([]*dto.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var categories []*dto.Category
	rows, err := r.db.QueryxContext(
		ctx,
		`SELECT
			c.id,
			c.name,
			c.icon,
			COALESCE(
				JSON_AGG(
					json_build_object('id', s.id, 'name', s.name)
					ORDER BY s.name
				) FILTER (WHERE s.deleted_at IS NULL),
				'[]'
		) AS subs
		FROM categories c
		LEFT JOIN subcategories s ON s.category_id = c.id
		GROUP BY c.id`,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve categories with subs", fault.WithError(err))
	}
	defer rows.Close()

	for rows.Next() {
		var category dto.Category
		var rawSubs []byte

		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Icon,
			&rawSubs,
		); err != nil {
			return nil, fault.New("failed to scan category", fault.WithError(err))
		}

		if err := json.Unmarshal(rawSubs, &category.Subs); err != nil {
			return nil, fault.New("failed to unmarshal subcategories", fault.WithError(err))
		}

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fault.New("error during rows iteration", fault.WithError(err))
	}

	return categories, nil
}

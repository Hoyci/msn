package models

import "time"

type Subcategory struct {
	ID         string     `db:"id"`
	Name       string     `db:"name"`
	CategoryID string     `db:"category_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

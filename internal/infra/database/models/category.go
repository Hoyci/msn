package models

import "time"

type Category struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Icon      string     `db:"icon"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

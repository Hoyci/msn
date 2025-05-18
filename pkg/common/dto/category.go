package dto

import (
	"time"
)

type Category struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Icon      string     `db:"icon" json:"icon"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type Subcategory struct {
	ID         string     `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	CategoryID string     `db:"category_id" json:"category_id"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updated_at,omitempty"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

type SubcategoryMinimal struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryWithUserCount struct {
	Category               `db:",inline" json:",inline"`
	SubcategoriesUserCount int `db:"subcategories_user_count" json:"subcategories_user_count"`
}

type CategoryWithSubcategories struct {
	Category      `db:",inline" json:",inline"`
	Subcategories []SubcategoryMinimal `json:"subcategories"`
}

type CategoryResponse struct {
	Categories []*CategoryWithUserCount `json:"categories"`
}

type CategoryWithSubResponse struct {
	Categories []*CategoryWithSubcategories `json:"categories"`
}

package models

import (
	"time"
)

type User struct {
	ID            string     `db:"id"`
	Name          string     `db:"name"`
	Email         string     `db:"email"`
	Password      string     `db:"password"`
	AvatarURL     string     `db:"avatar_url"`
	RoleID        string     `db:"role_id"`
	SubcategoryID *string    `db:"subcategory_id"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

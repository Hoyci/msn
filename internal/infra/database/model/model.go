package model

import (
	"time"
)

type User struct {
	ID            string     `db:"id"`
	Name          string     `db:"name"`
	Email         string     `db:"email"`
	Password      string     `db:"password"`
	AvatarURL     *string    `db:"avatar_url"`
	UserRoleID    string     `db:"user_role_id"`
	SubcategoryID *string    `db:"subcategory_id"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

type UserRole struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Session struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	JTI       string    `db:"jti"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type Category struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Icon      string     `db:"icon"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Subcategory struct {
	ID         string     `db:"id"`
	Name       string     `db:"name"`
	CategoryID string     `db:"category_id"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

package dto

import "time"

type CreateUser struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	AvatarUrl       *string `json:"avatar_url"`
	UserRole        string  `json:"user_role"`
	SubcategoryID   *string `json:"subcategory_id,omitempty"`
}

type UserResponse struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	AvatarURL     *string    `json:"avatar_url"`
	UserRoleID    string     `json:"user_role_id"`
	SubcategoryID *string    `json:"subcategory_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

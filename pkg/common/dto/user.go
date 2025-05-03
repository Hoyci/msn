package dto

import "time"

type CreateUser struct {
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	ConfirmPassword string  `json:"confirm_password"`
	AvatarUrl       *string `json:"avatar_url"`
}

type UserResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	AvatarURL *string    `json:"avatar_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

package dto

import (
	"mime/multipart"
	"time"
)

type CreateAvatar struct {
	FileHeader *multipart.FileHeader
}

type CreateUser struct {
	Name            string                `json:"name"`
	Email           string                `json:"email"`
	Password        string                `json:"password"`
	ConfirmPassword string                `json:"confirm_password"`
	FileHeader      *multipart.FileHeader `json:"file_header"`
	UserRole        string                `json:"role"`
	SubcategoryID   *string               `json:"subcategory_id,omitempty"`
}

type UserResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	AvatarURL     string    `json:"avatar_url"`
	RoleID        string    `json:"role_id"`
	SubcategoryID *string   `json:"subcategory_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserTokenResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EnrichedUserResponse struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	AvatarURL      string       `json:"avatar_url"`
	HashedPassword string       `json:"-"`
	CreatedAt      time.Time    `json:"created_at"`
	DeletedAt      *time.Time   `json:"deleted_at"`
	Role           *Role        `json:"role,omitempty"`
	Subcategory    *Subcategory `json:"subcategory,omitempty"`
	Category       *Category    `json:"category,omitempty"`
}

type ProfessionalUserResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Email       string       `json:"email"`
	AvatarURL   string       `json:"avatar_url"`
	Subcategory *Subcategory `json:"subcategory"`
}

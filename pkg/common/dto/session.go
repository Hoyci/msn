package dto

import "time"

type CreateSession struct {
	UserID string `json:"user_id"`
	JTI    string `json:"jti"`
}

type SessionResponse struct {
	ID        string    `json:"id"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

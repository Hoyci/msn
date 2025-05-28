package models

import "time"

type Session struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	JTI       string    `db:"jti"`
	Active    bool      `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

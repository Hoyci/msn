package model

import "time"

type User struct {
	ID        string     `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	AvatarURL *string    `db:"avatar_url"`
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

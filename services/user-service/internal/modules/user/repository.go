package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"msn/pkg/common/fault"
	"msn/services/user-service/internal/infra/database/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repository {
	return &repo{db: db}
}

// func (r repo) Update(ctx context.Context, user model.User) error {
// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	defer cancel()
//
// 	var query = `
// 		UPDATE users
// 		SET
// 			name = :name,
// 			username = :username,
// 			email = :email,
// 			password = :password,
// 			avatar_url = :avatar_url,
// 			enabled = :enabled,
// 			locked = :locked,
// 			updated = :updated
// 		WHERE id = :id
// 	`
//
// 	_, err := r.db.NamedExecContext(ctx, query, user)
// 	if err != nil {
// 		return fault.New("failed to update user", fault.WithError(err))
// 	}
//
// 	return nil
// }

func (r repo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user by email", fault.WithError(err))
	}

	return &user, nil
}

// func (r repo) Delete(ctx context.Context, userId string) error {
// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	defer cancel()
//
// 	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userId)
// 	if err != nil {
// 		return fault.New("failed to delete user", fault.WithError(err))
// 	}
//
// 	return nil
// }

func (r repo) GetByID(ctx context.Context, userId string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var user model.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user", fault.WithError(err))
	}

	return &user, nil
}

func (r repo) Insert(ctx context.Context, user model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			id,
			name,
			email,
			password,
			avatar_url,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			:id,
			:name,
			:email,
			:password,
			:avatar_url,
			:created_at,
			:updated_at,
			:deleted_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		fmt.Println(err)
		return fault.New("failed to insert user", fault.WithError(err))
	}

	return nil
}

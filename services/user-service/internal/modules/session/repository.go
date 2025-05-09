package session

import (
	"context"
	"database/sql"
	"errors"
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

func (r *repo) GetAllByUserID(ctx context.Context, userID string) ([]model.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sessions := make([]model.Session, 0)
	err := r.db.SelectContext(
		ctx,
		&sessions,
		"SELECT * FROM sessions WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fault.New("failed to retrieve sessions by user ID", fault.WithError(err))
	}

	return sessions, nil
}

func (r *repo) GetActiveByUserID(ctx context.Context, userID string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var session model.Session
	err := r.db.GetContext(
		ctx,
		&session,
		"SELECT * FROM sessions WHERE user_id = $1 AND active = true LIMIT 1",
		userID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New(
			"failed to retrieve active session",
			fault.WithError(err),
		)
	}

	return &session, nil
}

func (r *repo) GetByJTI(ctx context.Context, JTI string) (*model.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var session model.Session
	err := r.db.GetContext(
		ctx,
		&session,
		"SELECT * FROM sessions WHERE jti = $1",
		JTI,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve session by JTI", fault.WithError(err))
	}

	return &session, nil
}

func (r *repo) Insert(ctx context.Context, session model.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		INSERT INTO sessions (
			id,
			user_id,
			jti,
			active,
			created_at,
			updated_at,
			expires_at
		)	VALUES (
			:id,
			:user_id,
			:jti,
			:active,
			:created_at,
			:updated_at,
			:expires_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fault.New("failed to insert session", fault.WithError(err))
	}

	return nil
}

func (r *repo) Update(ctx context.Context, session model.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		UPDATE sessions
		SET
			active = :active,
			jti = :jti,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, session)
	if err != nil {
		return fault.New("failed to update session", fault.WithError(err))
	}

	return nil
}

func (r repo) DeactivateAll(ctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET active = false WHERE user_id = $1", userId)
	if err != nil {
		return fault.New("failed to update session", fault.WithError(err))
	}

	return nil
}

package sessionRepository

import (
	"context"
	"database/sql"
	"errors"
	"msn/internal/infra/database/models"
	"msn/internal/modules/session"
	"msn/pkg/common/fault"
	"time"

	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) session.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) GetAllByUserID(ctx context.Context, userID string) ([]*session.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	dbSessions := make([]models.Session, 0)
	err := r.db.SelectContext(
		ctx,
		&dbSessions,
		"SELECT * FROM sessions WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fault.New("failed to retrieve sessions by user ID", fault.WithError(err))
	}

	result := make([]*session.Session, len(dbSessions))
	for i, ms := range dbSessions {
		result[i] = session.NewFromModel(ms)
	}

	return result, nil
}

func (r *sessionRepository) GetActiveByUserID(ctx context.Context, userID string) (*session.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var sessionModel models.Session
	err := r.db.GetContext(
		ctx,
		&sessionModel,
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

	return session.NewFromModel(sessionModel), nil
}

func (r *sessionRepository) GetByJTI(ctx context.Context, JTI string) (*session.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var sessionModel models.Session
	err := r.db.GetContext(
		ctx,
		&sessionModel,
		"SELECT * FROM sessions WHERE jti = $1",
		JTI,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve session by JTI", fault.WithError(err))
	}

	return session.NewFromModel(sessionModel), nil
}

func (r *sessionRepository) Create(ctx context.Context, session *session.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	modelSession := session.ToModel()

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

	_, err := r.db.NamedExecContext(ctx, query, modelSession)
	if err != nil {
		return fault.New("failed to insert session", fault.WithError(err))
	}

	return nil
}

func (r *sessionRepository) Update(ctx context.Context, session *session.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	modelSession := session.ToModel()
	query := `
		UPDATE sessions
		SET
			active = :active,
			jti = :jti,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, modelSession)
	if err != nil {
		return fault.New("failed to update session", fault.WithError(err))
	}

	return nil
}

func (r sessionRepository) DeactivateAll(ctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET active = false WHERE user_id = $1", userId)
	if err != nil {
		return fault.New("failed to update session", fault.WithError(err))
	}

	return nil
}

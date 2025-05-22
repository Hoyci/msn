package session

import (
	"context"
	"msn/pkg/common/dto"
)

type SessionService interface {
	CreateSession(ctx context.Context, input dto.CreateSession) (*Session, error)
	DeactivateAllSessions(ctx context.Context, userID string) error
	GetActiveSessionByUserID(ctx context.Context, userID string) (*Session, error)
	UpdateSession(ctx context.Context, session *Session) (*Session, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	Update(ctx context.Context, session *Session) error
	GetAllByUserID(ctx context.Context, userID string) ([]*Session, error)
	GetActiveByUserID(ctx context.Context, userID string) (*Session, error)
	GetByJTI(ctx context.Context, JTI string) (*Session, error)
	DeactivateAll(ctx context.Context, userID string) error
}

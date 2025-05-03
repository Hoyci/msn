package session

import (
	"context"
	"msn/pkg/common/dto"
	"msn/services/user-service/internal/infra/database/model"
)

type Service interface {
	CreateSession(ctx context.Context, input dto.CreateSession) (*dto.SessionResponse, error)
}

type Repository interface {
	Insert(ctx context.Context, session model.Session) error
	Update(ctx context.Context, session model.Session) error
	GetAllByUserID(ctx context.Context, userID string) ([]model.Session, error)
	GetByJTI(ctx context.Context, JTI string) (*model.Session, error)
	DeactivateAll(ctx context.Context, userID string) error
}

package session

import (
	"context"
	"msn/internal/modules/user"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
)

type ServiceConfig struct {
	SessionRepo SessionRepository
	UserService user.UserService
}

type service struct {
	sessionRepo SessionRepository
	userService user.UserService
}

func NewService(c ServiceConfig) SessionService {
	return &service{
		sessionRepo: c.SessionRepo,
		userService: c.UserService,
	}
}

func (s service) CreateSession(ctx context.Context, input dto.CreateSession) (*Session, error) {
	sess, err := New(input.UserID, input.JTI)
	if err != nil {
		return nil, fault.NewUnprocessableEntity("failed to create session entity")
	}

	err = s.sessionRepo.Create(ctx, sess)
	if err != nil {
		return nil, fault.NewBadRequest("failed to insert session entity")
	}

	return sess, nil
}

func (s service) DeactivateAllSessions(ctx context.Context, userID string) error {
	err := s.sessionRepo.DeactivateAll(ctx, userID)
	if err != nil {
		return fault.NewBadRequest("failed to deactivate all user sessions")
	}

	return nil
}

func (s service) GetActiveSessionByUserID(ctx context.Context, userID string) (*Session, error) {
	sess, err := s.sessionRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, fault.NewBadRequest("failed to get activer user sessions")
	}

	return sess, nil
}

func (s service) UpdateSession(ctx context.Context, session *Session) (*Session, error) {
	err := s.sessionRepo.Update(ctx, session)
	if err != nil {
		return nil, fault.NewBadRequest("failed to update session")
	}

	return session, nil
}

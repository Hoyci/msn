package session

import (
	"context"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/services/user-service/internal/modules/user"
)

type ServiceConfig struct {
	SessionRepo Repository
	UserService user.Service
}

type service struct {
	sessionRepo Repository
	userService user.Service
}

func NewService(c ServiceConfig) Service {
	return &service{
		sessionRepo: c.SessionRepo,
		userService: c.UserService,
	}
}

func (s service) CreateSession(ctx context.Context, input dto.CreateSession) (*dto.SessionResponse, error) {
	userRecord, err := s.userService.GetUserByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}
	userID := userRecord.ID

	sess, err := New(userID, input.JTI)
	if err != nil {
		return nil, fault.NewUnprocessableEntity("failed to create session entity")
	}

	err = s.sessionRepo.Insert(ctx, sess.Model())
	if err != nil {
		return nil, fault.NewBadRequest("failed to insert session entity")
	}

	res := dto.SessionResponse{
		ID:        sess.ID,
		Active:    sess.Active,
		CreatedAt: sess.CreatedAt,
		UpdatedAt: sess.UpdatedAt,
	}

	return &res, nil
}

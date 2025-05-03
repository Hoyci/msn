package session

import (
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"msn/services/user-service/internal/infra/database/model"
	"time"
)

const (
	ttl = time.Hour * 24 * 30
)

type session struct {
	ID        string
	UserID    string
	JTI       string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func NewFromModel(m model.Session) *session {
	return &session{
		ID:        m.ID,
		UserID:    m.UserID,
		JTI:       m.JTI,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		ExpiresAt: m.ExpiresAt,
	}
}

func New(userID, JTI string) (*session, error) {
	session := session{
		ID:        uid.New("sess"),
		UserID:    userID,
		JTI:       JTI,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}

	if err := session.validate(); err != nil {
		return nil, fault.New(
			"failed to create session entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &session, nil
}

func (s *session) Model() model.Session {
	return model.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		JTI:       s.JTI,
		Active:    s.Active,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}

func (s *session) validate() error {
	if s.UserID == "" {
		return fault.New("UserID is required")
	}

	if s.JTI == "" {
		return fault.New("JTI is required")
	}

	return nil
}

func (s *session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *session) ChangeJTI(JTI string) {
	s.JTI = JTI
	s.UpdatedAt = time.Now()
}

func (s *session) Activate() {
	s.Active = true
	s.UpdatedAt = time.Now()
}

func (s *session) Deactivate() {
	s.Active = false
	s.UpdatedAt = time.Now()
}

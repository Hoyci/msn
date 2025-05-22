package session

import (
	"msn/internal/infra/database/model"
	"msn/pkg/common/fault"
	"msn/pkg/utils/uid"
	"time"
)

const (
	ttl = time.Hour * 24 * 30
)

type Session struct {
	ID        string
	UserID    string
	JTI       string
	Active    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

func New(userID, JTI string) (*Session, error) {
	if userID == "" || JTI == "" {
		return nil, fault.New("userID and JTI are required")
	}

	return &Session{
		ID:        uid.New("sess"),
		UserID:    userID,
		JTI:       JTI,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}, nil
}

func NewFromModel(m model.Session) *Session {
	return &Session{
		ID:        m.ID,
		UserID:    m.UserID,
		JTI:       m.JTI,
		Active:    m.Active,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		ExpiresAt: m.ExpiresAt,
	}
}

func (s *Session) ToModel() model.Session {
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

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}

func (s *Session) ChangeJTI(JTI string) {
	s.JTI = JTI
	s.UpdatedAt = time.Now()
}

func (s *Session) Activate() {
	s.Active = true
	s.UpdatedAt = time.Now()
}

func (s *Session) Deactivate() {
	s.Active = false
	s.UpdatedAt = time.Now()
}

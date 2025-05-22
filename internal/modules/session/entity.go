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
	id        string
	userID    string
	jti       string
	active    bool
	createdAt time.Time
	updatedAt time.Time
	expiresAt time.Time
}

func New(userID, JTI string) (*Session, error) {
	session := Session{
		id:        uid.New("sess"),
		userID:    userID,
		jti:       JTI,
		active:    true,
		createdAt: time.Now(),
		updatedAt: time.Now(),
		expiresAt: time.Now().Add(ttl),
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

func NewFromModel(m model.Session) *Session {
	return &Session{
		id:        m.ID,
		userID:    m.UserID,
		jti:       m.JTI,
		active:    m.Active,
		createdAt: m.CreatedAt,
		updatedAt: m.UpdatedAt,
		expiresAt: m.ExpiresAt,
	}
}

func (s *Session) ToModel() model.Session {
	return model.Session{
		ID:        s.id,
		UserID:    s.userID,
		JTI:       s.jti,
		Active:    s.active,
		CreatedAt: s.createdAt,
		UpdatedAt: s.updatedAt,
		ExpiresAt: s.expiresAt,
	}
}

func (s *Session) validate() error {
	if s.userID == "" {
		return fault.New("UserID is required")
	}

	if s.jti == "" {
		return fault.New("JTI is required")
	}

	return nil
}

func (s *Session) IsExpired() bool {
	return s.expiresAt.Before(time.Now())
}

func (s *Session) ChangeJTI(JTI string) {
	s.jti = JTI
	s.updatedAt = time.Now()
}

func (s *Session) Activate() {
	s.active = true
	s.updatedAt = time.Now()
}

func (s *Session) Deactivate() {
	s.active = false
	s.updatedAt = time.Now()
}

func (s *Session) ID() string           { return s.id }
func (s *Session) UserID() string       { return s.userID }
func (s *Session) JTI() string          { return s.jti }
func (s *Session) Active() bool         { return s.active }
func (s *Session) CreatedAt() time.Time { return s.createdAt }
func (s *Session) UpdatedAt() time.Time { return s.updatedAt }
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }

package user

import (
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/pkg/utils/uid"
	"msn/services/user-service/internal/infra/database/model"
	"time"
)

type user struct {
	id              string
	name            string
	email           string
	password        string
	confirmPassword string
	avatar_url      *string
	created_at      time.Time
	updated_at      *time.Time
	deleted_at      *time.Time
}

func NewFromModel(m model.User) *user {
	return &user{
		id:         m.ID,
		name:       m.Name,
		email:      m.Email,
		password:   m.Password,
		avatar_url: m.AvatarURL,
		created_at: m.CreatedAt,
		updated_at: m.UpdatedAt,
	}
}

func New(name, email, password, confirmPassword string, avatar_url *string) (*user, error) {
	u := user{
		id:              uid.New("user"),
		name:            name,
		email:           email,
		password:        password,
		confirmPassword: confirmPassword,
		avatar_url:      avatar_url,
		created_at:      time.Now(),
		updated_at:      nil,
		deleted_at:      nil,
	}

	if err := u.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	hashedPass, err := crypto.HashPassword(password)
	if err != nil {
		return nil, fault.New("failed to hash password", fault.WithError(err))
	}

	u.password = hashedPass

	return &u, nil
}

func (u *user) Model() model.User {
	return model.User{
		ID:        u.id,
		Name:      u.name,
		Email:     u.email,
		Password:  u.password,
		AvatarURL: u.avatar_url,
		CreatedAt: u.created_at,
		UpdatedAt: u.updated_at,
		DeletedAt: u.deleted_at,
	}
}

func (u *user) validate() error {
	if u.name == "" {
		return fault.New("user name is required")
	}
	if u.password == "" {
		return fault.New("password is required")
	}
	if u.email == "" {
		return fault.New("email is required")
	}
	if u.confirmPassword == "" {
		return fault.New("confirm_password is required")
	}
	if u.password != u.confirmPassword {
		return fault.New("password and confirm_password doesnt match")
	}

	return nil
}

package user

import (
	"msn/internal/infra/database/model"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/pkg/utils/uid"
	"net/mail"
	"time"
)

type user struct {
	id              string
	name            string
	email           string
	password        string
	confirmPassword string
	avatar_url      *string
	user_role_id    string
	subcategory_id  *string
	created_at      time.Time
	updated_at      *time.Time
	deleted_at      *time.Time
}

func NewFromModel(m model.User) *user {
	return &user{
		id:             m.ID,
		name:           m.Name,
		email:          m.Email,
		password:       m.Password,
		avatar_url:     m.AvatarURL,
		user_role_id:   m.UserRoleID,
		subcategory_id: m.SubcategoryID,
		created_at:     m.CreatedAt,
		updated_at:     m.UpdatedAt,
	}
}

func New(name, email, password, confirm_password, user_role_id string, avatar_url, subcategory_id *string) (*user, error) {
	u := user{
		id:              uid.New("user"),
		name:            name,
		email:           email,
		password:        password,
		confirmPassword: confirm_password,
		avatar_url:      avatar_url,
		user_role_id:    user_role_id,
		subcategory_id:  subcategory_id,
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
		ID:            u.id,
		Name:          u.name,
		Email:         u.email,
		Password:      u.password,
		AvatarURL:     u.avatar_url,
		UserRoleID:    u.user_role_id,
		SubcategoryID: u.subcategory_id,
		CreatedAt:     u.created_at,
		UpdatedAt:     u.updated_at,
		DeletedAt:     u.deleted_at,
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
	if _, err := mail.ParseAddress(u.email); err != nil {
		return fault.New("invalid email format")
	}
	if u.confirmPassword == "" {
		return fault.New("confirm_password is required")
	}
	if u.password != u.confirmPassword {
		return fault.New("password and confirm_password doesnt match")
	}
	if u.user_role_id == "" {
		return fault.New("user_role_id is required")
	}

	return nil
}

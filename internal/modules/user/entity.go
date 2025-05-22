package user

import (
	"context"
	"msn/internal/infra/database/model"
	categories "msn/internal/modules/categories"
	"msn/pkg/common/fault"
	"msn/pkg/utils/crypto"
	"msn/pkg/utils/uid"
	"net/mail"
	"time"
	"unicode"
)

type User struct {
	id              string
	name            string
	email           string
	password        string
	confirmPassword string
	avatarURL       *string
	userRoleID      string
	subcategoryID   *string
	createdAt       time.Time
	updatedAt       *time.Time
	deletedAt       *time.Time
}

func New(
	name,
	email,
	password,
	confirmPassword,
	userRoleID string,
	avatarURL,
	subcategoryID *string,
	userRepo UserRepository,
	categoryRepo categories.Repository,
) (*User, error) {
	user := User{
		id:              uid.New("user"),
		name:            name,
		email:           email,
		password:        password,
		confirmPassword: confirmPassword,
		avatarURL:       avatarURL,
		userRoleID:      userRoleID,
		subcategoryID:   subcategoryID,
		createdAt:       time.Now(),
		updatedAt:       nil,
		deletedAt:       nil,
	}

	if err := user.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	existing, err := userRepo.GetByEmail(context.Background(), user.email)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to validate email")
	}
	if existing != nil {
		return nil, fault.NewConflict("email already taken")
	}

	exists, err := userRepo.RoleExists(context.Background(), userRoleID)
	if err != nil {
		return nil, fault.NewInternalServerError("failed to validate user role")
	}
	if !exists {
		return nil, fault.NewBadRequest("invalid user role")
	}

	hashedPass, err := crypto.HashPassword(password)
	if err != nil {
		return nil, fault.New("failed to hash password", fault.WithError(err))
	}

	user.password = hashedPass

	return &user, nil
}

func NewFromModel(m model.User) *User {
	return &User{
		id:            m.ID,
		name:          m.Name,
		email:         m.Email,
		password:      m.Password,
		avatarURL:     m.AvatarURL,
		userRoleID:    m.UserRoleID,
		subcategoryID: m.SubcategoryID,
		createdAt:     m.CreatedAt,
		updatedAt:     m.UpdatedAt,
		deletedAt:     m.DeletedAt,
	}
}

func (u *User) ToModel() model.User {
	return model.User{
		ID:            u.id,
		Name:          u.name,
		Email:         u.email,
		Password:      u.password,
		AvatarURL:     u.avatarURL,
		UserRoleID:    u.userRoleID,
		SubcategoryID: u.subcategoryID,
		CreatedAt:     u.createdAt,
		UpdatedAt:     u.updatedAt,
		DeletedAt:     u.deletedAt,
	}
}

func (u *User) validate() error {
	if u.name == "" {
		return fault.NewBadRequest("user name is required")
	}
	if _, err := mail.ParseAddress(u.email); err != nil {
		return fault.NewBadRequest("invalid email format")
	}
	if u.password == "" {
		return fault.NewBadRequest("password is required")
	}
	if u.password != u.confirmPassword {
		return fault.NewBadRequest("password and confirmation do not match")
	}
	if len(u.password) < 8 {
		return fault.NewBadRequest("password must be at least 8 characters")
	}
	var hasUpper, hasLower, hasNumber, hasSymbol bool
	for _, c := range u.password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsSymbol(c), unicode.IsPunct(c):
			hasSymbol = true
		}
	}
	if !hasUpper || !hasLower || !hasNumber || !hasSymbol {
		return fault.NewBadRequest("password must contain uppercase, lowercase, numbers and symbol")
	}
	if u.userRoleID == "" {
		return fault.NewBadRequest("user role is required")
	}
	if u.userRoleID == "client" && u.subcategoryID != nil {
		return fault.NewBadRequest("clients cannot have subcategories")
	}
	if u.userRoleID == "professional" && u.subcategoryID == nil {
		return fault.NewBadRequest("professionals must have a subcategory")
	}

	return nil
}

func (u *User) ID() string             { return u.id }
func (u *User) Name() string           { return u.name }
func (u *User) Email() string          { return u.email }
func (u *User) UserRoleID() string     { return u.userRoleID }
func (u *User) SubcategoryID() *string { return u.subcategoryID }
func (u *User) AvatarURL() *string     { return u.avatarURL }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) Password() string       { return u.password }
func (u *User) DeletedAt() *time.Time  { return u.deletedAt }

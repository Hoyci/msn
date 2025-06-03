package user

import (
	"msn/internal/infra/database/models"
	"msn/internal/modules/role"
	"msn/internal/modules/subcategory"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"msn/pkg/common/valueobjects"
	"time"
)

type User struct {
	id           string
	name         string
	email        valueobjects.Email
	passwordHash string
	avatarURL    string
	role         role.Role
	subcategory  *subcategory.Subcategory
	createdAt    time.Time
	updatedAt    *time.Time
	deletedAt    *time.Time
}

func New(
	ID, name, rawEmail, hashedPassword, avatarURL, roleID string,
	subcatID *string,
) (*User, error) {
	emailVO, err := valueobjects.NewEmail(rawEmail)
	if err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	now := time.Now()

	user := &User{
		id:           ID,
		name:         name,
		email:        emailVO,
		passwordHash: hashedPassword,
		avatarURL:    avatarURL,
		role:         *role.FromID(roleID),
		subcategory:  nil,
		createdAt:    now,
		updatedAt:    nil,
		deletedAt:    nil,
	}

	if subcatID != nil {
		user.subcategory = subcategory.FromID(*subcatID)
	}

	if err := user.validate(); err != nil {
		return nil, err
	}

	return user, nil
}

func NewFromModel(m models.User) (*User, error) {
	emailVO, err := valueobjects.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	user := &User{
		id:           m.ID,
		name:         m.Name,
		email:        emailVO,
		passwordHash: m.Password,
		avatarURL:    m.AvatarURL,
		role:         *role.FromID(m.RoleID),
		subcategory:  nil,
		createdAt:    m.CreatedAt,
		updatedAt:    m.UpdatedAt,
		deletedAt:    m.DeletedAt,
	}

	if m.SubcategoryID != nil {
		user.subcategory = subcategory.FromID(*m.SubcategoryID)
	}

	return user, nil
}

func (u *User) ToModel() models.User {
	var subcatID *string
	if u.subcategory != nil {
		id := u.subcategory.ID()
		subcatID = &id
	}

	return models.User{
		ID:            u.id,
		Name:          u.name,
		Email:         u.email.Value,
		Password:      u.passwordHash,
		AvatarURL:     u.avatarURL,
		RoleID:        u.role.ID(),
		SubcategoryID: subcatID,
		CreatedAt:     u.createdAt,
		UpdatedAt:     u.updatedAt,
		DeletedAt:     u.deletedAt,
	}
}

func (u *User) validate() error {
	if u.Name() == "" {
		return fault.NewBadRequest("user name is required")
	}

	if u.role.ID() == "" {
		return fault.NewBadRequest("user role is required")
	}

	switch u.role.ID() {
	case "client":
		if u.subcategory != nil {
			return fault.NewBadRequest("clients cannot have subcategories")
		}
	case "professional":
		if u.subcategory == nil {
			return fault.NewBadRequest("professionals must have a subcategory")
		}
	}
	return nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email.Value
}

func (u *User) AvatarURL() string {
	return u.avatarURL
}

func (u *User) Role() role.Role {
	return u.role
}

func (u *User) Subcategory() *subcategory.Subcategory {
	return u.subcategory
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

func (u *User) ToResponse() *dto.UserResponse {
	var subcatID *string
	if u.subcategory != nil {
		id := u.subcategory.ID()
		subcatID = &id
	}

	return &dto.UserResponse{
		ID:            u.ID(),
		Name:          u.Name(),
		Email:         u.Email(),
		AvatarURL:     u.AvatarURL(),
		RoleID:        u.role.ID(),
		SubcategoryID: subcatID,
		CreatedAt:     u.CreatedAt(),
	}
}

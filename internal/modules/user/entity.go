package user

import (
	"msn/internal/infra/database/models"
	"msn/pkg/common/fault"
	"msn/pkg/common/valueobjects"
	"msn/pkg/utils/uid"
	"time"
)

// TODO: transform this into a struct that not export any data
// and create functions like func(u *User) ID { return u.ID } to export data
type User struct {
	ID            string
	Name          string
	Email         string
	Password      string
	AvatarURL     string
	UserRoleID    string
	SubcategoryID *string
	CreatedAt     time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
}

func New(
	name,
	email,
	hashedPassword,
	userRoleID string,
	subcategoryID *string,
) (*User, error) {
	user := User{
		ID:            uid.New("user"),
		Name:          name,
		Email:         email,
		Password:      hashedPassword,
		UserRoleID:    userRoleID,
		SubcategoryID: subcategoryID,
		CreatedAt:     time.Now(),
		UpdatedAt:     nil,
		DeletedAt:     nil,
	}

	if err := user.validate(); err != nil {
		return nil, fault.New(
			"failed to create user entity",
			fault.WithTag(fault.INVALID_ENTITY),
			fault.WithError(err),
		)
	}

	return &user, nil
}

func NewFromModel(m models.User) *User {
	return &User{
		ID:            m.ID,
		Name:          m.Name,
		Email:         m.Email,
		Password:      m.Password,
		AvatarURL:     m.AvatarURL,
		UserRoleID:    m.UserRoleID,
		SubcategoryID: m.SubcategoryID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
}

func (u *User) ToModel() models.User {
	return models.User{
		ID:            u.ID,
		Name:          u.Name,
		Email:         u.Email,
		Password:      u.Password,
		AvatarURL:     u.AvatarURL,
		UserRoleID:    u.UserRoleID,
		SubcategoryID: u.SubcategoryID,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     u.DeletedAt,
	}
}

func (u *User) validate() error {
	if u.Name == "" {
		return fault.NewBadRequest("user name is required")
	}
	if _, err := valueobjects.NewEmail(u.Email); err != nil {
		return err
	}
	if u.UserRoleID == "" {
		return fault.NewBadRequest("user role is required")
	}
	if u.UserRoleID == "client" && u.SubcategoryID != nil {
		return fault.NewBadRequest("clients cannot have subcategories")
	}
	if u.UserRoleID == "professional" && u.SubcategoryID == nil {
		return fault.NewBadRequest("professionals must have a subcategory")
	}

	return nil
}

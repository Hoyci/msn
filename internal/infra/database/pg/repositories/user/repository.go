package userRepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"msn/internal/infra/database/models"
	"msn/internal/modules/user"
	"msn/pkg/common/dto"
	"msn/pkg/common/fault"
	"time"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) user.UserRepository {
	return &userRepo{db: db}
}

// func (r repo) Update(ctx context.Context, user model.User) error {
// 	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
// 	defer cancel()
//
// 	var query = `
// 		UPDATE users
// 		SET
// 			name = :name,
// 			username = :username,
// 			email = :email,
// 			password = :password,
// 			avatar_url = :avatar_url,
// 			enabled = :enabled,
// 			locked = :locked,
// 			updated = :updated
// 		WHERE id = :id
// 	`
//
// 	_, err := r.db.NamedExecContext(ctx, query, user)
// 	if err != nil {
// 		return fault.New("failed to update user", fault.WithError(err))
// 	}
//
// 	return nil
// }

func (r userRepo) GetEnrichedByEmail(ctx context.Context, email string) (*dto.EnrichedUserResponse, error) {
	var out struct {
		ID              string     `db:"id"`
		Name            string     `db:"name"`
		Email           string     `db:"email"`
		AvatarURL       string     `db:"avatar_url"`
		HashedPassword  string     `db:"password"`
		CreatedAt       time.Time  `db:"created_at"`
		DeletedAt       *time.Time `db:"deleted_at"`
		RoleID          string     `db:"role_id"`
		RoleName        string     `db:"role_name"`
		SubcategoryID   *string    `db:"subcategory_id"`
		SubcategoryName *string    `db:"subcategory_name"`
		CategoryID      *string    `db:"category_id"`
		CategoryName    *string    `db:"category_name"`
		CategoryIcon    *string    `db:"category_icon"`
	}

	query := `
    SELECT
      u.id, u.name, u.email, u.avatar_url, u.password, u.created_at, u.deleted_at,
      ur.id   AS role_id,          ur.name AS role_name,
      s.id    AS subcategory_id,   s.name AS subcategory_name,
      c.id    AS category_id,      c.name AS category_name, c.icon as category_icon
    FROM users u
    LEFT JOIN roles ur    ON ur.id = u.role_id
    LEFT JOIN subcategories s  ON s.id  = u.subcategory_id
    LEFT JOIN categories c     ON c.id  = s.category_id
    WHERE u.email = $1
    `
	err := r.db.GetContext(ctx, &out, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to get enriched user", fault.WithError(err))
	}

	enrichedUser := &dto.EnrichedUserResponse{
		ID:             out.ID,
		Name:           out.Name,
		Email:          out.Email,
		AvatarURL:      out.AvatarURL,
		HashedPassword: out.HashedPassword,
		CreatedAt:      out.CreatedAt,
		DeletedAt:      out.DeletedAt,
		UserRole: &dto.UserRole{
			ID:   out.RoleID,
			Name: out.RoleName,
		},
	}

	if out.SubcategoryID != nil {
		enrichedUser.Subcategory = &dto.Subcategory{ID: *out.SubcategoryID, Name: *out.SubcategoryName}
	}
	if out.CategoryID != nil {
		enrichedUser.Category = &dto.Category{ID: *out.CategoryID, Name: *out.CategoryName, Icon: *out.CategoryIcon}
	}

	return enrichedUser, nil
}

func (r userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var modelUser models.User
	err := r.db.GetContext(ctx, &modelUser, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user by email", fault.WithError(err))
	}

	return user.NewFromModel(modelUser), nil
}

func (r userRepo) GetByID(ctx context.Context, userId string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var modelUser models.User
	err := r.db.GetContext(ctx, &modelUser, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user", fault.WithError(err))
	}

	return user.NewFromModel(modelUser), nil
}

func (r userRepo) Create(ctx context.Context, user *user.User) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	modelUser := user.ToModel()

	query := `
		INSERT INTO users (
			id,
			name,
			email,
			password,
			avatar_url,
			role_id,
			subcategory_id,
			created_at,
			updated_at,
			deleted_at
		) VALUES (
			:id,
			:name,
			:email,
			:password,
			:avatar_url,
			:role_id,
			:subcategory_id,
			:created_at,
			:updated_at,
			:deleted_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, modelUser)
	if err != nil {
		return fault.New("failed to insert user", fault.WithError(err))
	}

	return nil
}

func (r userRepo) GetProfessionalUsers(ctx context.Context) ([]*dto.ProfessionalUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var out []struct {
		ID              string `db:"id"`
		Name            string `db:"name"`
		Email           string `db:"email"`
		AvatarURL       string `db:"avatar_url"`
		SubcategoryID   string `db:"subcategory_id"`
		SubcategoryName string `db:"subcategory_name"`
		CategoryID      string `db:"category_id"`
	}

	query := `
		SELECT 
			u.id,
			u.name,
			u.email,
			u.avatar_url,
			s.id as subcategory_id,
			s."name" as subcategory_name,
			s.category_id as category_id
		FROM users u
		LEFT JOIN roles ur ON ur.id = u.role_id
		LEFT JOIN subcategories s ON s.id = u.subcategory_id
		WHERE ur."name" = 'professional';
	`

	err := r.db.SelectContext(ctx, &out, query)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user", fault.WithError(err))
	}

	fmt.Println(out[0])

	var users []*dto.ProfessionalUserResponse

	for _, u := range out {
		users = append(users, &dto.ProfessionalUserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			AvatarURL: u.AvatarURL,
			Subcategory: &dto.Subcategory{
				ID:         u.SubcategoryID,
				Name:       u.SubcategoryName,
				CategoryID: u.CategoryID,
			},
		})
	}

	return users, nil
}

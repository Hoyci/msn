package user

import (
	"context"
	"database/sql"
	"errors"
	"msn/internal/infra/database/model"
	"msn/pkg/common/fault"
	"time"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) UserRepository {
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

func (r userRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var modelUser model.User
	err := r.db.GetContext(ctx, &modelUser, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user by email", fault.WithError(err))
	}

	return NewFromModel(modelUser), nil
}

func (r userRepo) GetByID(ctx context.Context, userId string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var modelUser model.User
	err := r.db.GetContext(ctx, &modelUser, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user", fault.WithError(err))
	}

	return NewFromModel(modelUser), nil
}

func (r userRepo) Create(ctx context.Context, user *User) error {
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
			user_role_id,
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
			:user_role_id,
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

func (r userRepo) GetUserRoleByID(ctx context.Context, ID string) (*model.UserRole, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		SELECT
			*
		FROM user_roles ur
		WHERE ur.id = $1
		AND ur.deleted_at IS NULL;
	`

	var userRole model.UserRole
	err := r.db.GetContext(ctx, &userRole, query, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fault.New("failed to retrieve user", fault.WithError(err))
	}

	return &userRole, nil
}

func (r userRepo) RoleExists(ctx context.Context, roleID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_roles WHERE id = $1 AND deleted_at IS NULL)`
	err := r.db.GetContext(ctx, &exists, query, roleID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

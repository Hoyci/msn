package userrole

import (
	"context"
	"msn/internal/infra/database/models"
	userrole "msn/internal/modules/user_role"
	"time"

	"github.com/jmoiron/sqlx"
)

type userRoleRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) userrole.Repository {
	return &userRoleRepo{db: db}
}

func (u userRoleRepo) GetUserRoleByName(ctx context.Context, userRoleName string) (*userrole.UserRole, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var userRole models.UserRole
	query := `SELECT * FROM user_roles WHERE name = $1 AND deleted_at IS NULL`
	err := u.db.GetContext(ctx, &userRole, query, userRoleName)
	if err != nil {
		return nil, err
	}
	return userrole.NewFromModel(userRole), nil
}

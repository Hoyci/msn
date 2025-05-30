package role

import (
	"context"
	"msn/internal/infra/database/models"
	"msn/internal/modules/role"
	"time"

	"github.com/jmoiron/sqlx"
)

type roleRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) role.Repository {
	return &roleRepo{db: db}
}

func (u roleRepo) GetRoleByName(ctx context.Context, userRoleName string) (*role.Role, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var modelRole models.Role
	query := `SELECT * FROM roles WHERE name = $1 AND deleted_at IS NULL`
	err := u.db.GetContext(ctx, &modelRole, query, userRoleName)
	if err != nil {
		return nil, err
	}
	return role.NewFromModel(modelRole), nil
}

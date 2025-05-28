package userrole

import "context"

type Repository interface {
	GetUserRoleByName(ctx context.Context, userRoleName string) (*UserRole, error)
}

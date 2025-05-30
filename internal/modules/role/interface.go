package role

import "context"

type Repository interface {
	GetRoleByName(ctx context.Context, roleName string) (*Role, error)
}

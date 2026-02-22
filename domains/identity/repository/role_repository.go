package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/google/uuid"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) (*entity.Role, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	List(ctx context.Context) ([]*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) (*entity.Role, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) error
}

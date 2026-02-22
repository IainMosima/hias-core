package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/google/uuid"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *entity.Permission) (*entity.Permission, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	List(ctx context.Context) ([]*entity.Permission, error)
	ListByRole(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error)
}

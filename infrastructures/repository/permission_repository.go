package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/identity/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type permissionRepository struct {
	store db.Store
}

func NewPermissionRepository(store db.Store) domainRepo.PermissionRepository {
	return &permissionRepository{store: store}
}

func (r *permissionRepository) Create(ctx context.Context, permission *entity.Permission) (*entity.Permission, error) {
	dbPerm, err := r.store.CreatePermission(ctx, db.CreatePermissionParams{
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}
	return sqlcPermissionToDomain(dbPerm), nil
}

func (r *permissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	dbPerm, err := r.store.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	return sqlcPermissionToDomain(dbPerm), nil
}

func (r *permissionRepository) List(ctx context.Context) ([]*entity.Permission, error) {
	dbPerms, err := r.store.ListPermissions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	perms := make([]*entity.Permission, len(dbPerms))
	for i, p := range dbPerms {
		perms[i] = sqlcPermissionToDomain(p)
	}
	return perms, nil
}

func (r *permissionRepository) ListByRole(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error) {
	dbPerms, err := r.store.ListPermissionsByRole(ctx, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions by role: %w", err)
	}
	perms := make([]*entity.Permission, len(dbPerms))
	for i, p := range dbPerms {
		perms[i] = sqlcPermissionToDomain(p)
	}
	return perms, nil
}

func sqlcPermissionToDomain(p db.Permission) *entity.Permission {
	return &entity.Permission{
		ID:          p.ID,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

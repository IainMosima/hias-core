package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/identity/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type roleRepository struct {
	store db.Store
}

func NewRoleRepository(store db.Store) domainRepo.RoleRepository {
	return &roleRepository{store: store}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	dbRole, err := r.store.CreateRole(ctx, db.CreateRoleParams{
		Name:        role.Name,
		Description: role.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return sqlcRoleToDomain(dbRole), nil
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	dbRole, err := r.store.GetRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return sqlcRoleToDomain(dbRole), nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	dbRole, err := r.store.GetRoleByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	return sqlcRoleToDomain(dbRole), nil
}

func (r *roleRepository) List(ctx context.Context) ([]*entity.Role, error) {
	dbRoles, err := r.store.ListRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	roles := make([]*entity.Role, len(dbRoles))
	for i, role := range dbRoles {
		roles[i] = sqlcRoleToDomain(role)
	}
	return roles, nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) (*entity.Role, error) {
	dbRole, err := r.store.UpdateRole(ctx, db.UpdateRoleParams{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}
	return sqlcRoleToDomain(dbRole), nil
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteRole(ctx, id)
}

func (r *roleRepository) AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	_, err := r.store.AssignPermissionToRole(ctx, db.AssignPermissionToRoleParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
	return err
}

func (r *roleRepository) RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return r.store.RemovePermissionFromRole(ctx, db.RemovePermissionFromRoleParams{
		RoleID:       roleID,
		PermissionID: permissionID,
	})
}

func sqlcRoleToDomain(r db.Role) *entity.Role {
	return &entity.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

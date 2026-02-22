package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/identity/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type userRepository struct {
	store db.Store
}

func NewUserRepository(store db.Store) domainRepo.UserRepository {
	return &userRepository{store: store}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser, err := r.store.CreateUser(ctx, db.CreateUserParams{
		CognitoSub: stringToPgtypeText(user.CognitoSub),
		Email:      user.Email,
		Name:       user.Name,
		Phone:      stringToPgtypeText(user.Phone),
		NationalID: stringToPgtypeText(user.NationalID),
		RoleID:     user.RoleID,
		Status:     user.Status,
		CreatedBy:  uuidToPgtype(user.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	dbUser, err := r.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) GetByCognitoSub(ctx context.Context, sub string) (*entity.User, error) {
	dbUser, err := r.store.GetUserByCognitoSub(ctx, stringToPgtypeText(sub))
	if err != nil {
		return nil, fmt.Errorf("failed to get user by cognito sub: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) GetByNationalID(ctx context.Context, nationalID string) (*entity.User, error) {
	dbUser, err := r.store.GetUserByNationalID(ctx, stringToPgtypeText(nationalID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user by national ID: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	dbUsers, err := r.store.ListUsers(ctx, db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	users := make([]*entity.User, len(dbUsers))
	for i, u := range dbUsers {
		users[i] = sqlcUserToDomain(u)
	}
	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.store.CountUsers(ctx)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		ID:         user.ID,
		Name:       stringToPgtypeText(user.Name),
		Phone:      stringToPgtypeText(user.Phone),
		NationalID: stringToPgtypeText(user.NationalID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.User, error) {
	dbUser, err := r.store.UpdateUserStatus(ctx, db.UpdateUserStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) UpdateRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) (*entity.User, error) {
	dbUser, err := r.store.UpdateUserRole(ctx, db.UpdateUserRoleParams{
		ID:     id,
		RoleID: roleID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func (r *userRepository) UpdateCognitoSub(ctx context.Context, id uuid.UUID, sub string) (*entity.User, error) {
	dbUser, err := r.store.UpdateUserCognitoSub(ctx, db.UpdateUserCognitoSubParams{
		ID:         id,
		CognitoSub: stringToPgtypeText(sub),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update cognito sub: %w", err)
	}
	return sqlcUserToDomain(dbUser), nil
}

func sqlcUserToDomain(u db.User) *entity.User {
	return &entity.User{
		ID:         u.ID,
		CognitoSub: u.CognitoSub.String,
		Email:      u.Email,
		Name:       u.Name,
		Phone:      u.Phone.String,
		NationalID: u.NationalID.String,
		RoleID:     u.RoleID,
		Status:     u.Status,
		CreatedBy:  pgtypeToUUID(u.CreatedBy),
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/identity/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByCognitoSub(ctx context.Context, sub string) (*entity.User, error)
	GetByNationalID(ctx context.Context, nationalID string) (*entity.User, error)
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.User, error)
	UpdateRole(ctx context.Context, id uuid.UUID, roleID uuid.UUID) (*entity.User, error)
	UpdateCognitoSub(ctx context.Context, id uuid.UUID, sub string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, hash string) error
}

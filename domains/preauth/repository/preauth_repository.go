package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/preauth/entity"
	"github.com/google/uuid"
)

type PreAuthRepository interface {
	Create(ctx context.Context, preauth *entity.PreAuthorization) (*entity.PreAuthorization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PreAuthorization, error)
	GetByAuthCode(ctx context.Context, authCode string) (*entity.PreAuthorization, error)
	List(ctx context.Context, limit, offset int) ([]*entity.PreAuthorization, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.PreAuthorization, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.PreAuthorization, error)
	ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.PreAuthorization, error)
	Count(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.PreAuthorization, error)
	Approve(ctx context.Context, preauth *entity.PreAuthorization) (*entity.PreAuthorization, error)
	Deny(ctx context.Context, id uuid.UUID, reason string, reviewedBy uuid.UUID) (*entity.PreAuthorization, error)
	GetExpiring(ctx context.Context) ([]*entity.PreAuthorization, error)
}

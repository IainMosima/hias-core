package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type PolicyRenewalRepository interface {
	Create(ctx context.Context, renewal *entity.PolicyRenewal) (*entity.PolicyRenewal, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.PolicyRenewal, error)
	GetByPolicyID(ctx context.Context, policyID uuid.UUID) (*entity.PolicyRenewal, error)
	ListPending(ctx context.Context) ([]*entity.PolicyRenewal, error)
	ListExpired(ctx context.Context) ([]*entity.PolicyRenewal, error)
	Update(ctx context.Context, renewal *entity.PolicyRenewal) (*entity.PolicyRenewal, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.PolicyRenewal, error)
}

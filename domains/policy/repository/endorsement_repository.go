package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type EndorsementRepository interface {
	Create(ctx context.Context, endorsement *entity.Endorsement) (*entity.Endorsement, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Endorsement, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.Endorsement, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Endorsement, error)
	Update(ctx context.Context, endorsement *entity.Endorsement) (*entity.Endorsement, error)
}

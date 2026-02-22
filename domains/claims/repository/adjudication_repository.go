package repository

import (
	"context"
	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type AdjudicationRepository interface {
	Create(ctx context.Context, decision *entity.AdjudicationDecision) (*entity.AdjudicationDecision, error)
	GetByClaimID(ctx context.Context, claimID uuid.UUID) (*entity.AdjudicationDecision, error)
	ListByDecision(ctx context.Context, decision string, limit, offset int) ([]*entity.AdjudicationDecision, error)
}

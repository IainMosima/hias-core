package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/google/uuid"
)

type UnderwritingRepository interface {
	Create(ctx context.Context, assessment *entity.UnderwritingAssessment) (*entity.UnderwritingAssessment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.UnderwritingAssessment, error)
	GetByPolicyID(ctx context.Context, policyID uuid.UUID) ([]*entity.UnderwritingAssessment, error)
	GetByMemberID(ctx context.Context, memberID uuid.UUID) (*entity.UnderwritingAssessment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.UnderwritingAssessment, error)
	Update(ctx context.Context, assessment *entity.UnderwritingAssessment) (*entity.UnderwritingAssessment, error)
}

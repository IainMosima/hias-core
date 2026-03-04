package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	"github.com/google/uuid"
)

type LeadActivityRepository interface {
	Create(ctx context.Context, activity *entity.LeadActivity) (*entity.LeadActivity, error)
	ListByLead(ctx context.Context, leadID uuid.UUID) ([]*entity.LeadActivity, error)
}

package repository

import (
	"context"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/google/uuid"
)

type RecoveryWorkflowEventRepository interface {
	Create(ctx context.Context, event *entity.RecoveryWorkflowEvent) (*entity.RecoveryWorkflowEvent, error)
	ListByRecovery(ctx context.Context, recoveryID uuid.UUID) ([]*entity.RecoveryWorkflowEvent, error)
}

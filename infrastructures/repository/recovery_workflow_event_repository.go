package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type recoveryWorkflowEventRepository struct {
	store db.Store
}

func NewRecoveryWorkflowEventRepository(store db.Store) domainRepo.RecoveryWorkflowEventRepository {
	return &recoveryWorkflowEventRepository{store: store}
}

func (r *recoveryWorkflowEventRepository) Create(ctx context.Context, event *entity.RecoveryWorkflowEvent) (*entity.RecoveryWorkflowEvent, error) {
	dbEvent, err := r.store.CreateRecoveryWorkflowEvent(ctx, db.CreateRecoveryWorkflowEventParams{
		RecoveryID:  event.RecoveryID,
		FromStatus:  event.FromStatus,
		ToStatus:    event.ToStatus,
		EventType:   event.EventType,
		Notes:       stringToPgtypeText(event.Notes),
		PerformedBy: event.PerformedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create recovery workflow event: %w", err)
	}
	return sqlcRecoveryWorkflowEventToDomain(dbEvent), nil
}

func (r *recoveryWorkflowEventRepository) ListByRecovery(ctx context.Context, recoveryID uuid.UUID) ([]*entity.RecoveryWorkflowEvent, error) {
	dbEvents, err := r.store.ListRecoveryWorkflowEventsByRecovery(ctx, recoveryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list recovery workflow events by recovery: %w", err)
	}
	events := make([]*entity.RecoveryWorkflowEvent, len(dbEvents))
	for i, e := range dbEvents {
		events[i] = sqlcRecoveryWorkflowEventToDomain(e)
	}
	return events, nil
}

func sqlcRecoveryWorkflowEventToDomain(e db.RecoveryWorkflowEvent) *entity.RecoveryWorkflowEvent {
	return &entity.RecoveryWorkflowEvent{
		ID:          e.ID,
		RecoveryID:  e.RecoveryID,
		FromStatus:  e.FromStatus,
		ToStatus:    e.ToStatus,
		EventType:   e.EventType,
		Notes:       e.Notes.String,
		PerformedBy: e.PerformedBy,
		CreatedAt:   e.CreatedAt,
	}
}

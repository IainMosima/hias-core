package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/audit/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/audit/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

// auditRepository is APPEND ONLY - no Update or Delete methods
type auditRepository struct {
	store db.Store
}

func NewAuditRepository(store db.Store) domainRepo.AuditRepository {
	return &auditRepository{store: store}
}

func (r *auditRepository) Create(ctx context.Context, event *entity.AuditEvent) (*entity.AuditEvent, error) {
	dbEvent, err := r.store.CreateAuditEvent(ctx, db.CreateAuditEventParams{
		UserID:     uuidToPgtype(event.UserID),
		EntityType: event.EntityType,
		EntityID:   event.EntityID,
		Action:     event.Action,
		OldValue:   event.OldValue,
		NewValue:   event.NewValue,
		IpAddress:  stringToPgtypeText(event.IPAddress),
		UserAgent:  stringToPgtypeText(event.UserAgent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create audit event: %w", err)
	}
	return sqlcAuditEventToDomain(dbEvent), nil
}

func (r *auditRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditEvent, error) {
	dbEvent, err := r.store.GetAuditEventByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit event by ID: %w", err)
	}
	return sqlcAuditEventToDomain(dbEvent), nil
}

func (r *auditRepository) List(ctx context.Context, limit, offset int) ([]*entity.AuditEvent, error) {
	dbEvents, err := r.store.ListAuditEvents(ctx, db.ListAuditEventsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list audit events: %w", err)
	}
	events := make([]*entity.AuditEvent, len(dbEvents))
	for i, e := range dbEvents {
		events[i] = sqlcAuditEventToDomain(e)
	}
	return events, nil
}

func (r *auditRepository) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID, limit, offset int) ([]*entity.AuditEvent, error) {
	dbEvents, err := r.store.ListAuditEventsByEntity(ctx, db.ListAuditEventsByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list audit events by entity: %w", err)
	}
	events := make([]*entity.AuditEvent, len(dbEvents))
	for i, e := range dbEvents {
		events[i] = sqlcAuditEventToDomain(e)
	}
	return events, nil
}

func (r *auditRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.AuditEvent, error) {
	dbEvents, err := r.store.ListAuditEventsByUser(ctx, db.ListAuditEventsByUserParams{
		UserID: uuidToPgtype(userID),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list audit events by user: %w", err)
	}
	events := make([]*entity.AuditEvent, len(dbEvents))
	for i, e := range dbEvents {
		events[i] = sqlcAuditEventToDomain(e)
	}
	return events, nil
}

func (r *auditRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountAuditEvents(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit events: %w", err)
	}
	return count, nil
}

func (r *auditRepository) CountByEntity(ctx context.Context, entityType string, entityID uuid.UUID) (int64, error) {
	count, err := r.store.CountAuditEventsByEntity(ctx, db.CountAuditEventsByEntityParams{
		EntityType: entityType,
		EntityID:   entityID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count audit events by entity: %w", err)
	}
	return count, nil
}

func (r *auditRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := r.store.CountAuditEventsByUser(ctx, uuidToPgtype(userID))
	if err != nil {
		return 0, fmt.Errorf("failed to count audit events by user: %w", err)
	}
	return count, nil
}

func sqlcAuditEventToDomain(e db.AuditEvent) *entity.AuditEvent {
	return &entity.AuditEvent{
		ID:         e.ID,
		UserID:     pgtypeToUUID(e.UserID),
		EntityType: e.EntityType,
		EntityID:   e.EntityID,
		Action:     e.Action,
		OldValue:   e.OldValue,
		NewValue:   e.NewValue,
		IPAddress:  e.IpAddress.String,
		UserAgent:  e.UserAgent.String,
		CreatedAt:  e.CreatedAt,
	}
}

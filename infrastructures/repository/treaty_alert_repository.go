package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type treatyAlertRepository struct {
	store db.Store
}

func NewTreatyAlertRepository(store db.Store) domainRepo.TreatyAlertRepository {
	return &treatyAlertRepository{store: store}
}

func (r *treatyAlertRepository) Create(ctx context.Context, alert *entity.TreatyAlert) (*entity.TreatyAlert, error) {
	dbAlert, err := r.store.CreateTreatyAlert(ctx, db.CreateTreatyAlertParams{
		TreatyID:       alert.TreatyID,
		TreatyLayerID:  uuidToPgtype(alert.TreatyLayerID),
		AlertType:      alert.AlertType,
		Severity:       alert.Severity,
		Message:        alert.Message,
		ThresholdValue: alert.ThresholdValue,
		CurrentValue:   alert.CurrentValue,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create treaty alert: %w", err)
	}
	return sqlcTreatyAlertToDomain(dbAlert), nil
}

func (r *treatyAlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TreatyAlert, error) {
	dbAlert, err := r.store.GetTreatyAlertByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get treaty alert by ID: %w", err)
	}
	return sqlcTreatyAlertToDomain(dbAlert), nil
}

func (r *treatyAlertRepository) List(ctx context.Context, limit, offset int) ([]*entity.TreatyAlert, error) {
	dbAlerts, err := r.store.ListTreatyAlerts(ctx, db.ListTreatyAlertsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list treaty alerts: %w", err)
	}
	alerts := make([]*entity.TreatyAlert, len(dbAlerts))
	for i, a := range dbAlerts {
		alerts[i] = sqlcTreatyAlertToDomain(a)
	}
	return alerts, nil
}

func (r *treatyAlertRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID, limit, offset int) ([]*entity.TreatyAlert, error) {
	dbAlerts, err := r.store.ListTreatyAlertsByTreaty(ctx, db.ListTreatyAlertsByTreatyParams{
		TreatyID: treatyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list treaty alerts by treaty: %w", err)
	}
	alerts := make([]*entity.TreatyAlert, len(dbAlerts))
	for i, a := range dbAlerts {
		alerts[i] = sqlcTreatyAlertToDomain(a)
	}
	return alerts, nil
}

func (r *treatyAlertRepository) ListUnacknowledged(ctx context.Context, limit, offset int) ([]*entity.TreatyAlert, error) {
	dbAlerts, err := r.store.ListUnacknowledgedTreatyAlerts(ctx, db.ListUnacknowledgedTreatyAlertsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list unacknowledged treaty alerts: %w", err)
	}
	alerts := make([]*entity.TreatyAlert, len(dbAlerts))
	for i, a := range dbAlerts {
		alerts[i] = sqlcTreatyAlertToDomain(a)
	}
	return alerts, nil
}

func (r *treatyAlertRepository) Acknowledge(ctx context.Context, id uuid.UUID, acknowledgedBy uuid.UUID) (*entity.TreatyAlert, error) {
	dbAlert, err := r.store.AcknowledgeTreatyAlert(ctx, db.AcknowledgeTreatyAlertParams{
		ID:             id,
		AcknowledgedBy: uuidToPgtype(acknowledgedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to acknowledge treaty alert: %w", err)
	}
	return sqlcTreatyAlertToDomain(dbAlert), nil
}

func (r *treatyAlertRepository) CountUnacknowledged(ctx context.Context) (int64, error) {
	count, err := r.store.CountUnacknowledgedTreatyAlerts(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count unacknowledged treaty alerts: %w", err)
	}
	return count, nil
}

func sqlcTreatyAlertToDomain(a db.TreatyAlert) *entity.TreatyAlert {
	return &entity.TreatyAlert{
		ID:             a.ID,
		TreatyID:       a.TreatyID,
		TreatyLayerID:  pgtypeToUUID(a.TreatyLayerID),
		AlertType:      a.AlertType,
		Severity:       a.Severity,
		Message:        a.Message,
		ThresholdValue: a.ThresholdValue,
		CurrentValue:   a.CurrentValue,
		IsAcknowledged: a.IsAcknowledged,
		AcknowledgedBy: pgtypeToUUID(a.AcknowledgedBy),
		AcknowledgedAt: pgtypeTimestamptzToTimePtr(a.AcknowledgedAt),
		CreatedAt:      a.CreatedAt,
	}
}

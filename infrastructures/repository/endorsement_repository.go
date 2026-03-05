package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type endorsementRepository struct {
	store db.Store
}

func NewEndorsementRepository(store db.Store) domainRepo.EndorsementRepository {
	return &endorsementRepository{store: store}
}

func (r *endorsementRepository) Create(ctx context.Context, e *entity.Endorsement) (*entity.Endorsement, error) {
	dbE, err := r.store.CreateEndorsement(ctx, db.CreateEndorsementParams{
		PolicyID:          e.PolicyID,
		EndorsementType:   e.EndorsementType,
		Status:            e.Status,
		EffectiveDate:     e.EffectiveDate,
		Changes:           e.Changes,
		Reason:            stringToPgtypeText(e.Reason),
		PremiumAdjustment: int64ToPgtypeInt8(e.PremiumAdjustment),
		RequestedBy:       e.RequestedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create endorsement: %w", err)
	}
	return sqlcEndorsementToDomain(dbE), nil
}

func (r *endorsementRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Endorsement, error) {
	dbE, err := r.store.GetEndorsementByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get endorsement: %w", err)
	}
	return sqlcEndorsementToDomain(dbE), nil
}

func (r *endorsementRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.Endorsement, error) {
	dbEs, err := r.store.ListEndorsementsByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list endorsements: %w", err)
	}
	result := make([]*entity.Endorsement, len(dbEs))
	for i, e := range dbEs {
		result[i] = sqlcEndorsementToDomain(e)
	}
	return result, nil
}

func (r *endorsementRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Endorsement, error) {
	dbE, err := r.store.UpdateEndorsementStatus(ctx, db.UpdateEndorsementStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update endorsement status: %w", err)
	}
	return sqlcEndorsementToDomain(dbE), nil
}

func (r *endorsementRepository) Update(ctx context.Context, e *entity.Endorsement) (*entity.Endorsement, error) {
	dbE, err := r.store.UpdateEndorsement(ctx, db.UpdateEndorsementParams{
		ID:         e.ID,
		Status:     stringToPgtypeText(e.Status),
		ApprovedBy: uuidToPgtype(e.ApprovedBy),
		ApprovedAt: timePtrToPgtypeTimestamptz(e.ApprovedAt),
		AppliedAt:  timePtrToPgtypeTimestamptz(e.AppliedAt),
		Reason:     stringToPgtypeText(e.Reason),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update endorsement: %w", err)
	}
	return sqlcEndorsementToDomain(dbE), nil
}

func sqlcEndorsementToDomain(e db.Endorsement) *entity.Endorsement {
	var premiumAdj int64
	if e.PremiumAdjustment.Valid {
		premiumAdj = e.PremiumAdjustment.Int64
	}
	var changes json.RawMessage
	if e.Changes != nil {
		changes = e.Changes
	} else {
		changes = json.RawMessage("{}")
	}
	return &entity.Endorsement{
		ID:                e.ID,
		PolicyID:          e.PolicyID,
		EndorsementType:   e.EndorsementType,
		Status:            e.Status,
		EffectiveDate:     e.EffectiveDate,
		Changes:           changes,
		Reason:            e.Reason.String,
		PremiumAdjustment: premiumAdj,
		RequestedBy:       e.RequestedBy,
		ApprovedBy:        pgtypeToUUID(e.ApprovedBy),
		ApprovedAt:        pgtypeTimestamptzToTimePtr(e.ApprovedAt),
		AppliedAt:         pgtypeTimestamptzToTimePtr(e.AppliedAt),
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

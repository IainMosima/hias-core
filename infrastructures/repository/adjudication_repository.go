package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type adjudicationRepository struct {
	store db.Store
}

func NewAdjudicationRepository(store db.Store) domainRepo.AdjudicationRepository {
	return &adjudicationRepository{store: store}
}

func (r *adjudicationRepository) Create(ctx context.Context, decision *entity.AdjudicationDecision) (*entity.AdjudicationDecision, error) {
	dbDecision, err := r.store.CreateAdjudicationDecision(ctx, db.CreateAdjudicationDecisionParams{
		ClaimID:              decision.ClaimID,
		Decision:             decision.Decision,
		PayableAmount:        decision.PayableAmount,
		MemberResponsibility: decision.MemberResponsibility,
		Reasons:              decision.Reasons,
		RuleResults:          decision.RuleResults,
		AdjudicatedBy:        uuidToPgtype(decision.AdjudicatedBy),
		AdjudicatedAt:        decision.AdjudicatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create adjudication decision: %w", err)
	}
	return sqlcAdjudicationToDomain(dbDecision), nil
}

func (r *adjudicationRepository) GetByClaimID(ctx context.Context, claimID uuid.UUID) (*entity.AdjudicationDecision, error) {
	dbDecision, err := r.store.GetAdjudicationByClaimID(ctx, claimID)
	if err != nil {
		return nil, fmt.Errorf("failed to get adjudication by claim ID: %w", err)
	}
	return sqlcAdjudicationToDomain(dbDecision), nil
}

func (r *adjudicationRepository) ListByDecision(ctx context.Context, decision string, limit, offset int) ([]*entity.AdjudicationDecision, error) {
	dbDecisions, err := r.store.ListAdjudicationsByDecision(ctx, db.ListAdjudicationsByDecisionParams{
		Decision: decision,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list adjudications by decision: %w", err)
	}
	decisions := make([]*entity.AdjudicationDecision, len(dbDecisions))
	for i, d := range dbDecisions {
		decisions[i] = sqlcAdjudicationToDomain(d)
	}
	return decisions, nil
}

func sqlcAdjudicationToDomain(d db.AdjudicationDecision) *entity.AdjudicationDecision {
	return &entity.AdjudicationDecision{
		ID:                   d.ID,
		ClaimID:              d.ClaimID,
		Decision:             d.Decision,
		PayableAmount:        d.PayableAmount,
		MemberResponsibility: d.MemberResponsibility,
		Reasons:              d.Reasons,
		RuleResults:          d.RuleResults,
		AdjudicatedBy:        pgtypeToUUID(d.AdjudicatedBy),
		AdjudicatedAt:        d.AdjudicatedAt,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
	}
}

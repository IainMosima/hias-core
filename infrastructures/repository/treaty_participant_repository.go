package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/reinsurance/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type treatyParticipantRepository struct {
	store db.Store
}

func NewTreatyParticipantRepository(store db.Store) domainRepo.TreatyParticipantRepository {
	return &treatyParticipantRepository{store: store}
}

func (r *treatyParticipantRepository) Create(ctx context.Context, participant *entity.TreatyParticipant) (*entity.TreatyParticipant, error) {
	dbParticipant, err := r.store.CreateTreatyParticipant(ctx, db.CreateTreatyParticipantParams{
		TreatyID:        participant.TreatyID,
		ReinsurerName:   participant.ReinsurerName,
		SharePercentage: float64ToPgNumeric(participant.SharePercentage),
		CommissionRate:  float64ToPgNumeric(participant.CommissionRate),
		IsLead:          participant.IsLead,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create treaty participant: %w", err)
	}
	return sqlcTreatyParticipantToDomain(dbParticipant), nil
}

func (r *treatyParticipantRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.TreatyParticipant, error) {
	dbParticipant, err := r.store.GetTreatyParticipantByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get treaty participant by ID: %w", err)
	}
	return sqlcTreatyParticipantToDomain(dbParticipant), nil
}

func (r *treatyParticipantRepository) ListByTreaty(ctx context.Context, treatyID uuid.UUID) ([]*entity.TreatyParticipant, error) {
	dbParticipants, err := r.store.ListTreatyParticipantsByTreaty(ctx, treatyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list treaty participants by treaty: %w", err)
	}
	participants := make([]*entity.TreatyParticipant, len(dbParticipants))
	for i, p := range dbParticipants {
		participants[i] = sqlcTreatyParticipantToDomain(p)
	}
	return participants, nil
}

func (r *treatyParticipantRepository) Update(ctx context.Context, participant *entity.TreatyParticipant) (*entity.TreatyParticipant, error) {
	dbParticipant, err := r.store.UpdateTreatyParticipant(ctx, db.UpdateTreatyParticipantParams{
		ID:              participant.ID,
		ReinsurerName:   participant.ReinsurerName,
		SharePercentage: float64ToPgNumeric(participant.SharePercentage),
		CommissionRate:  float64ToPgNumeric(participant.CommissionRate),
		IsLead:          participant.IsLead,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update treaty participant: %w", err)
	}
	return sqlcTreatyParticipantToDomain(dbParticipant), nil
}

func (r *treatyParticipantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteTreatyParticipant(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete treaty participant: %w", err)
	}
	return nil
}

func (r *treatyParticipantRepository) GetTotalShareByTreaty(ctx context.Context, treatyID uuid.UUID) (float64, error) {
	totalShare, err := r.store.GetTotalShareByTreaty(ctx, treatyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total share by treaty: %w", err)
	}
	return pgNumericToFloat64(totalShare), nil
}

func sqlcTreatyParticipantToDomain(p db.TreatyParticipant) *entity.TreatyParticipant {
	return &entity.TreatyParticipant{
		ID:              p.ID,
		TreatyID:        p.TreatyID,
		ReinsurerName:   p.ReinsurerName,
		SharePercentage: pgNumericToFloat64(p.SharePercentage),
		CommissionRate:  pgNumericToFloat64(p.CommissionRate),
		IsLead:          p.IsLead,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

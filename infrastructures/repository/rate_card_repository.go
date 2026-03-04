package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/provider/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type rateCardRepository struct {
	store db.Store
}

func NewRateCardRepository(store db.Store) domainRepo.RateCardRepository {
	return &rateCardRepository{store: store}
}

func (r *rateCardRepository) Create(ctx context.Context, rateCard *entity.RateCard) (*entity.RateCard, error) {
	dbRateCard, err := r.store.CreateRateCard(ctx, db.CreateRateCardParams{
		ProviderID:    rateCard.ProviderID,
		ProcedureCode: rateCard.ProcedureCode,
		ProcedureName: rateCard.ProcedureName,
		RateAmount:    rateCard.RateAmount,
		EffectiveDate: rateCard.EffectiveDate,
		AgeFrom:       int32(rateCard.AgeFrom),
		AgeTo:         int32(rateCard.AgeTo),
		Gender:        stringToPgtypeText(rateCard.Gender),
		Relationship:  stringToPgtypeText(rateCard.Relationship),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create rate card: %w", err)
	}
	return sqlcRateCardToDomain(dbRateCard), nil
}

func (r *rateCardRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.RateCard, error) {
	dbRateCard, err := r.store.GetRateCardByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate card by ID: %w", err)
	}
	return sqlcRateCardToDomain(dbRateCard), nil
}

func (r *rateCardRepository) GetByProviderAndProcedure(ctx context.Context, providerID uuid.UUID, procedureCode string) (*entity.RateCard, error) {
	dbRateCard, err := r.store.GetRateByProviderAndProcedure(ctx, db.GetRateByProviderAndProcedureParams{
		ProviderID:    providerID,
		ProcedureCode: procedureCode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rate card by provider and procedure: %w", err)
	}
	return sqlcRateCardToDomain(dbRateCard), nil
}

func (r *rateCardRepository) GetByProviderProcedureAndAge(ctx context.Context, providerID uuid.UUID, procedureCode string, age int) (*entity.RateCard, error) {
	dbRateCard, err := r.store.GetRateByProviderProcedureAndAge(ctx, db.GetRateByProviderProcedureAndAgeParams{
		ProviderID:    providerID,
		ProcedureCode: procedureCode,
		AgeFrom:       int32(age),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get rate card by provider, procedure and age: %w", err)
	}
	return sqlcRateCardToDomain(dbRateCard), nil
}

func (r *rateCardRepository) ListByProvider(ctx context.Context, providerID uuid.UUID) ([]*entity.RateCard, error) {
	dbRateCards, err := r.store.ListRateCardsByProvider(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list rate cards by provider: %w", err)
	}
	rateCards := make([]*entity.RateCard, len(dbRateCards))
	for i, rc := range dbRateCards {
		rateCards[i] = sqlcRateCardToDomain(rc)
	}
	return rateCards, nil
}

func (r *rateCardRepository) Update(ctx context.Context, rateCard *entity.RateCard) (*entity.RateCard, error) {
	dbRateCard, err := r.store.UpdateRateCard(ctx, db.UpdateRateCardParams{
		ID:            rateCard.ID,
		RateAmount:    rateCard.RateAmount,
		EffectiveDate: rateCard.EffectiveDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update rate card: %w", err)
	}
	return sqlcRateCardToDomain(dbRateCard), nil
}

func (r *rateCardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteRateCard(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete rate card: %w", err)
	}
	return nil
}

func sqlcRateCardToDomain(rc db.RateCard) *entity.RateCard {
	return &entity.RateCard{
		ID:            rc.ID,
		ProviderID:    rc.ProviderID,
		ProcedureCode: rc.ProcedureCode,
		ProcedureName: rc.ProcedureName,
		RateAmount:    rc.RateAmount,
		EffectiveDate: rc.EffectiveDate,
		AgeFrom:       int(rc.AgeFrom),
		AgeTo:         int(rc.AgeTo),
		Gender:        rc.Gender.String,
		Relationship:  rc.Relationship.String,
		CreatedAt:     rc.CreatedAt,
		UpdatedAt:     rc.UpdatedAt,
	}
}

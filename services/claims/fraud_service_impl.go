package claims

import (
	"context"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type fraudServiceImpl struct {
	fraudFlagRepo claimRepo.FraudFlagRepository
}

func NewFraudService(fraudFlagRepo claimRepo.FraudFlagRepository) service.FraudService {
	return &fraudServiceImpl{fraudFlagRepo: fraudFlagRepo}
}

func (s *fraudServiceImpl) CheckDuplicate(ctx context.Context, claimNumber string, claimID uuid.UUID) (bool, error) {
	count, err := s.fraudFlagRepo.CheckDuplicate(ctx, claimNumber, claimID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *fraudServiceImpl) CheckFrequency(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, claimID uuid.UUID) (bool, error) {
	count, err := s.fraudFlagRepo.CheckFrequency(ctx, memberID, providerID, procedureCode, claimID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *fraudServiceImpl) CheckAmountThreshold(ctx context.Context, providerID uuid.UUID, procedureCode string, amount int64) (bool, error) {
	// Threshold check: flag if amount exceeds 500,000 KES
	return amount > shared.FraudAmountThresholdCents, nil
}

func (s *fraudServiceImpl) FlagClaim(ctx context.Context, flag *entity.FraudFlag) error {
	_, err := s.fraudFlagRepo.Create(ctx, flag)
	return err
}

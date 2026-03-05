package claims

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/claims/service"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type fraudServiceImpl struct {
	fraudFlagRepo claimRepo.FraudFlagRepository
	contractRepo  providerRepo.ContractRepository
	rateCardRepo  providerRepo.RateCardRepository
	providerRepo  providerRepo.ProviderRepository
}

func NewFraudService(
	fraudFlagRepo claimRepo.FraudFlagRepository,
	contractRepo providerRepo.ContractRepository,
	rateCardRepo providerRepo.RateCardRepository,
	providerRepo providerRepo.ProviderRepository,
) service.FraudService {
	return &fraudServiceImpl{
		fraudFlagRepo: fraudFlagRepo,
		contractRepo:  contractRepo,
		rateCardRepo:  rateCardRepo,
		providerRepo:  providerRepo,
	}
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
	return amount > shared.FraudAmountThresholdCents, nil
}

func (s *fraudServiceImpl) CheckExpiredContract(ctx context.Context, providerID uuid.UUID, serviceDate time.Time) (bool, error) {
	if s.contractRepo == nil {
		return false, nil
	}
	contracts, err := s.contractRepo.ListByProvider(ctx, providerID)
	if err != nil {
		return false, err
	}
	for _, c := range contracts {
		if c.Status == string(shared.ContractStatusActive) &&
			serviceDate.After(c.StartDate) && serviceDate.Before(c.EndDate) {
			return false, nil
		}
	}
	return true, nil
}

func (s *fraudServiceImpl) CheckSuspendedProvider(ctx context.Context, providerID uuid.UUID) (bool, error) {
	if s.providerRepo == nil {
		return false, nil
	}
	provider, err := s.providerRepo.GetByID(ctx, providerID)
	if err != nil {
		return false, err
	}
	return provider.Status == string(shared.ProviderStatusSuspended), nil
}

func (s *fraudServiceImpl) CheckRateCardOvercharge(ctx context.Context, providerID uuid.UUID, procedureCode string, amount int64) (bool, error) {
	if s.rateCardRepo == nil {
		return false, nil
	}
	rateCard, err := s.rateCardRepo.GetByProviderAndProcedure(ctx, providerID, procedureCode)
	if err != nil {
		return false, nil // no rate card = can't check
	}
	return amount > rateCard.RateAmount, nil
}

func (s *fraudServiceImpl) CheckRepeatVisit(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, serviceDate time.Time, claimID uuid.UUID) (bool, error) {
	// Reuse frequency check — same member+provider+procedure within recent window
	count, err := s.fraudFlagRepo.CheckFrequency(ctx, memberID, providerID, procedureCode, claimID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *fraudServiceImpl) FlagClaim(ctx context.Context, flag *entity.FraudFlag) error {
	_, err := s.fraudFlagRepo.Create(ctx, flag)
	return err
}

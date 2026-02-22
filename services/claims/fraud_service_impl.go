package claims

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/claims/service"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type fraudServiceImpl struct {
	claimRepo     repository.ClaimRepository
	rateCardRepo  providerRepo.RateCardRepository
	fraudFlagRepo repository.FraudFlagRepository
}

func NewFraudService(
	claimRepo repository.ClaimRepository,
	rateCardRepo providerRepo.RateCardRepository,
	fraudFlagRepo repository.FraudFlagRepository,
) service.FraudService {
	return &fraudServiceImpl{
		claimRepo:     claimRepo,
		rateCardRepo:  rateCardRepo,
		fraudFlagRepo: fraudFlagRepo,
	}
}

// CheckDuplicate looks for an existing claim with the same claim number or
// a matching member+provider+service_date combination within a recent period.
// The excludeID parameter ensures the current claim being processed is not
// matched against itself.
func (f *fraudServiceImpl) CheckDuplicate(ctx context.Context, claimNumber string, claimID uuid.UUID) (bool, error) {
	count, err := f.fraudFlagRepo.CheckDuplicate(ctx, claimNumber, claimID)
	if err != nil {
		return false, fmt.Errorf("failed to check duplicate for claim %s: %w", claimNumber, err)
	}

	isDuplicate := count > 0
	if isDuplicate {
		utils.LogWarn("Fraud: duplicate detected for claim %s (matches: %d)", claimNumber, count)
	}

	return isDuplicate, nil
}

// CheckFrequency detects if the same member has claimed the same procedure
// at the same provider within a 7-day window. This is a common indicator of
// potential claim fraud or billing errors.
func (f *fraudServiceImpl) CheckFrequency(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, claimID uuid.UUID) (bool, error) {
	count, err := f.fraudFlagRepo.CheckFrequency(ctx, memberID, providerID, procedureCode, claimID)
	if err != nil {
		return false, fmt.Errorf("failed to check frequency for member %s, provider %s, procedure %s: %w",
			memberID, providerID, procedureCode, err)
	}

	isFrequent := count > 0
	if isFrequent {
		utils.LogWarn("Fraud: frequency anomaly for member %s - procedure %s at provider %s (%d occurrences in 7 days)",
			memberID, procedureCode, providerID, count)
	}

	return isFrequent, nil
}

// CheckAmountThreshold compares the claimed amount against the provider's
// rate card for the procedure. If the amount exceeds 2x the standard rate,
// it is flagged as suspicious. Returns an error if no rate card exists for
// the procedure at the given provider.
func (f *fraudServiceImpl) CheckAmountThreshold(ctx context.Context, providerID uuid.UUID, procedureCode string, amount int64) (bool, error) {
	rateCard, err := f.rateCardRepo.GetByProviderAndProcedure(ctx, providerID, procedureCode)
	if err != nil {
		return false, fmt.Errorf("no rate card found for provider %s, procedure %s: %w",
			providerID, procedureCode, err)
	}

	threshold := rateCard.RateAmount * 2
	exceeds := amount > threshold

	if exceeds {
		utils.LogWarn("Fraud: amount %d cents for procedure %s exceeds 2x rate card (%d cents) at provider %s",
			amount, procedureCode, rateCard.RateAmount, providerID)
	}

	return exceeds, nil
}

// FlagClaim creates a fraud flag record associated with the given claim.
// These flags are stored for review by claims officers and can be resolved
// manually after investigation.
func (f *fraudServiceImpl) FlagClaim(ctx context.Context, flag *entity.FraudFlag) error {
	_, err := f.fraudFlagRepo.Create(ctx, flag)
	if err != nil {
		utils.LogError("Fraud: failed to create fraud flag for claim %s: %v", flag.ClaimID, err)
		return fmt.Errorf("failed to create fraud flag: %w", err)
	}

	utils.LogInfo("Fraud: flagged claim %s with type=%s severity=%s", flag.ClaimID, flag.FlagType, flag.Severity)
	return nil
}

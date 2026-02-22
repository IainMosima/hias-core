package claims

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/bitbiz/hias-core/domains/claims/service"
	policyDomainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
)

type validatorServiceImpl struct {
	policyRepo   policyDomainRepo.PolicyRepository
	memberRepo   policyDomainRepo.MemberRepository
	providerRepo providerRepo.ProviderRepository
}

func NewValidatorService(
	policyRepo policyDomainRepo.PolicyRepository,
	memberRepo policyDomainRepo.MemberRepository,
	providerRepo providerRepo.ProviderRepository,
) service.ValidatorService {
	return &validatorServiceImpl{
		policyRepo:   policyRepo,
		memberRepo:   memberRepo,
		providerRepo: providerRepo,
	}
}

// ValidateClaim runs a series of validation checks against the claim and its
// line items. It returns (valid, errorMessages, error) where:
//   - valid: true if all checks passed
//   - errorMessages: list of human-readable reasons for each failed check
//   - error: non-nil only if an infrastructure error occurred (DB failure, etc.)
func (v *validatorServiceImpl) ValidateClaim(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) (bool, []string, error) {
	var errors []string

	// ------------------------------------------------------------------
	// Check 1: Required fields present
	// ------------------------------------------------------------------
	if claim.PolicyID.String() == "00000000-0000-0000-0000-000000000000" {
		errors = append(errors, "Policy ID is required")
	}
	if claim.MemberID.String() == "00000000-0000-0000-0000-000000000000" {
		errors = append(errors, "Member ID is required")
	}
	if claim.ProviderID.String() == "00000000-0000-0000-0000-000000000000" {
		errors = append(errors, "Provider ID is required")
	}
	if claim.ServiceDate.IsZero() {
		errors = append(errors, "Service date is required")
	}
	if len(lineItems) == 0 {
		errors = append(errors, "At least one line item is required")
	}
	for i, li := range lineItems {
		if li.ProcedureCode == "" {
			errors = append(errors, fmt.Sprintf("Line item %d: procedure code is required", i+1))
		}
		if li.Quantity < 1 {
			errors = append(errors, fmt.Sprintf("Line item %d: quantity must be at least 1", i+1))
		}
		if li.UnitPrice < 1 {
			errors = append(errors, fmt.Sprintf("Line item %d: unit price must be at least 1 cent", i+1))
		}
	}

	// ------------------------------------------------------------------
	// Check 2: Provider exists and is ACTIVE
	// ------------------------------------------------------------------
	provider, err := v.providerRepo.GetByID(ctx, claim.ProviderID)
	if err != nil {
		utils.LogError("Validator: failed to fetch provider %s: %v", claim.ProviderID, err)
		errors = append(errors, "Provider not found")
	} else if provider.Status != string(shared.ProviderStatusActive) {
		errors = append(errors, fmt.Sprintf("Provider is not active (current status: %s)", provider.Status))
	}

	// ------------------------------------------------------------------
	// Check 3: Member exists
	// ------------------------------------------------------------------
	_, err = v.memberRepo.GetByID(ctx, claim.MemberID)
	if err != nil {
		utils.LogError("Validator: failed to fetch member %s: %v", claim.MemberID, err)
		errors = append(errors, "Member not found")
	}

	// ------------------------------------------------------------------
	// Check 4: Policy exists and is ACTIVE
	// ------------------------------------------------------------------
	policy, err := v.policyRepo.GetByID(ctx, claim.PolicyID)
	if err != nil {
		utils.LogError("Validator: failed to fetch policy %s: %v", claim.PolicyID, err)
		errors = append(errors, "Policy not found")
	} else if policy.Status != string(shared.PolicyStatusActive) {
		errors = append(errors, fmt.Sprintf("Policy is not active (current status: %s)", policy.Status))
	}

	// ------------------------------------------------------------------
	// Check 5: Service date is not in the future
	// ------------------------------------------------------------------
	if !claim.ServiceDate.IsZero() && claim.ServiceDate.After(time.Now()) {
		errors = append(errors, "Service date cannot be in the future")
	}

	valid := len(errors) == 0
	if !valid {
		utils.LogInfo("Validator: claim %s failed validation with %d errors", claim.ClaimNumber, len(errors))
	} else {
		utils.LogInfo("Validator: claim %s passed all validation checks", claim.ClaimNumber)
	}

	return valid, errors, nil
}

package claims

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/bitbiz/hias-core/domains/claims/service"
	memberRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
)

type validatorServiceImpl struct {
	policyRepo   policyRepo.PolicyRepository
	memberRepo   memberRepo.MemberRepository
	providerRepo providerRepo.ProviderRepository
}

func NewValidatorService(
	policyRepo policyRepo.PolicyRepository,
	memberRepo memberRepo.MemberRepository,
	providerRepo providerRepo.ProviderRepository,
) service.ValidatorService {
	return &validatorServiceImpl{
		policyRepo:   policyRepo,
		memberRepo:   memberRepo,
		providerRepo: providerRepo,
	}
}

func (s *validatorServiceImpl) ValidateClaim(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) (bool, []string, error) {
	var errors []string

	// 1. Policy must exist and be ACTIVE
	policy, err := s.policyRepo.GetByID(ctx, claim.PolicyID)
	if err != nil {
		errors = append(errors, "Policy not found")
		return false, errors, nil
	}
	if policy.Status != string(shared.PolicyStatusActive) {
		errors = append(errors, fmt.Sprintf("Policy is %s, must be ACTIVE", policy.Status))
	}

	// 2. Member must exist and be under the policy
	member, err := s.memberRepo.GetByID(ctx, claim.MemberID)
	if err != nil {
		errors = append(errors, "Member not found")
	} else if member.PolicyID != claim.PolicyID {
		errors = append(errors, "Member does not belong to this policy")
	}

	// 3. Provider must exist and be ACTIVE
	provider, err := s.providerRepo.GetByID(ctx, claim.ProviderID)
	if err != nil {
		errors = append(errors, "Provider not found")
	} else if provider.Status != string(shared.ProviderStatusActive) {
		errors = append(errors, fmt.Sprintf("Provider is %s, must be ACTIVE", provider.Status))
	}

	// 4. Must have at least one line item
	if len(lineItems) == 0 {
		errors = append(errors, "Claim must have at least one line item")
	}

	// 5. Total amount must be positive
	if claim.TotalAmount <= 0 {
		errors = append(errors, "Claim total amount must be positive")
	}

	return len(errors) == 0, errors, nil
}

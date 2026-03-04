package claims

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/claims/service"
	memberRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
)

type adjudicatorServiceImpl struct {
	claimRepo     claimRepo.ClaimRepository
	benefitRepo   productRepo.BenefitRepository
	exclusionRepo productRepo.ExclusionRepository
	policyRepo    policyRepo.PolicyRepository
	memberRepo    memberRepo.MemberRepository
	providerRepo  providerRepo.ProviderRepository
	networkRepo   productRepo.ProviderNetworkRepository
	fraudSvc      service.FraudService
}

func NewAdjudicatorService(
	claimRepo claimRepo.ClaimRepository,
	benefitRepo productRepo.BenefitRepository,
	exclusionRepo productRepo.ExclusionRepository,
	policyRepo policyRepo.PolicyRepository,
	memberRepo memberRepo.MemberRepository,
	providerRepo providerRepo.ProviderRepository,
	networkRepo productRepo.ProviderNetworkRepository,
	fraudSvc service.FraudService,
) service.AdjudicatorService {
	return &adjudicatorServiceImpl{
		claimRepo:     claimRepo,
		benefitRepo:   benefitRepo,
		exclusionRepo: exclusionRepo,
		policyRepo:    policyRepo,
		memberRepo:    memberRepo,
		providerRepo:  providerRepo,
		networkRepo:   networkRepo,
		fraudSvc:      fraudSvc,
	}
}

func (s *adjudicatorServiceImpl) Adjudicate(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) (*entity.AdjudicationResult, error) {
	var ruleResults []entity.RuleResult

	// 1. Eligibility check: policy active?
	policy, err := s.policyRepo.GetByID(ctx, claim.PolicyID)
	if err != nil || policy.Status != string(shared.PolicyStatusActive) {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility),
			Rule:     "policy_active",
			Result:   string(shared.RuleResultFail),
			Details:  "Policy is not active",
		})
		return &entity.AdjudicationResult{
			Decision:             string(shared.AdjudicationDecisionReject),
			PayableAmount:        0,
			MemberResponsibility: claim.TotalAmount,
			Reasons:              ruleResults,
		}, nil
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryEligibility),
		Rule:     "policy_active",
		Result:   string(shared.RuleResultPass),
		Details:  "Policy is active",
	})

	// 2. Member eligibility
	member, err := s.memberRepo.GetByID(ctx, claim.MemberID)
	if err != nil {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility), Rule: "member_active", Result: string(shared.RuleResultFail),
			Details: "Member not found",
		})
		return &entity.AdjudicationResult{
			Decision:             string(shared.AdjudicationDecisionReject),
			PayableAmount:        0,
			MemberResponsibility: claim.TotalAmount,
			Reasons:              ruleResults,
		}, nil
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryEligibility), Rule: "member_active", Result: string(shared.RuleResultPass),
		Details: "Member is enrolled",
	})

	// 3. Provider check
	provider, err := s.providerRepo.GetByID(ctx, claim.ProviderID)
	if err != nil || provider.Status != string(shared.ProviderStatusActive) {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility), Rule: "provider_active", Result: string(shared.RuleResultFail),
			Details: "Provider is not active",
		})
		return &entity.AdjudicationResult{
			Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
			MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
		}, nil
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryEligibility), Rule: "provider_active", Result: string(shared.RuleResultPass),
		Details: "Provider is active",
	})

	// 3b. Provider network eligibility check
	if s.networkRepo != nil {
		inNetwork, err := s.networkRepo.CheckEligibility(ctx, policy.PlanID, claim.ProviderID, "")
		if err == nil && !inNetwork {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "provider_network", Result: string(shared.RuleResultFail),
				Details: "Provider is not in the plan's network",
			})
			return &entity.AdjudicationResult{
				Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
				MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
			}, nil
		}
		if err == nil && inNetwork {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "provider_network", Result: string(shared.RuleResultPass),
				Details: "Provider is in network",
			})
		}
	}

	// 4. Coverage & limits check
	benefits, err := s.benefitRepo.ListByPlan(ctx, policy.PlanID)
	if err != nil || len(benefits) == 0 {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryCoverage), Rule: "benefit_exists", Result: string(shared.RuleResultFail),
			Details: "No benefits found for plan",
		})
		return &entity.AdjudicationResult{
			Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
			MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
		}, nil
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryCoverage), Rule: "benefit_exists", Result: string(shared.RuleResultPass),
		Details: "Benefits available",
	})

	// 5. Waiting period check
	for _, b := range benefits {
		if b.WaitingPeriodDays > 0 {
			waitingEnd := member.CreatedAt.AddDate(0, 0, b.WaitingPeriodDays)
			if claim.ServiceDate.Before(waitingEnd) {
				ruleResults = append(ruleResults, entity.RuleResult{
					Category: string(shared.RuleCategoryEligibility), Rule: "waiting_period", Result: string(shared.RuleResultFail),
					Details: fmt.Sprintf("Service date falls within %d-day waiting period for %s", b.WaitingPeriodDays, b.Category),
				})
				return &entity.AdjudicationResult{
					Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
					MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
				}, nil
			}
		}
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryEligibility), Rule: "waiting_period", Result: string(shared.RuleResultPass),
		Details: "Service date is past all waiting periods",
	})

	// 5b. Member age check against benefit min_age/max_age
	memberAge := calculateAge(member.DateOfBirth, claim.ServiceDate)
	for _, b := range benefits {
		if memberAge < b.MinAge || memberAge > b.MaxAge {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "age_eligibility", Result: string(shared.RuleResultFail),
				Details: fmt.Sprintf("Member age %d is outside benefit %s age range (%d-%d)", memberAge, b.Category, b.MinAge, b.MaxAge),
			})
			return &entity.AdjudicationResult{
				Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
				MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
			}, nil
		}
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryEligibility), Rule: "age_eligibility", Result: string(shared.RuleResultPass),
		Details: "Member age is within benefit eligibility range",
	})

	// 6. Exclusion check
	exclusions, _ := s.exclusionRepo.ListByPlan(ctx, policy.PlanID)
	if len(exclusions) > 0 {
		claimDiagCodes := parseJSONStringArray(claim.DiagnosisCodes)
		for _, excl := range exclusions {
			exclCodes := parseJSONStringArray(excl.ICDCodes)
			for _, diagCode := range claimDiagCodes {
				if contains(exclCodes, diagCode) {
					ruleResults = append(ruleResults, entity.RuleResult{
						Category: string(shared.RuleCategoryCoverage), Rule: "exclusion_check", Result: string(shared.RuleResultFail),
						Details: fmt.Sprintf("Diagnosis code %s matches exclusion: %s", diagCode, excl.Description),
					})
					return &entity.AdjudicationResult{
						Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
						MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
					}, nil
				}
			}
		}
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryCoverage), Rule: "exclusion_check", Result: string(shared.RuleResultPass),
		Details: "No exclusions matched",
	})

	// 7. Calculate payable amount based on benefits
	payableAmount := claim.TotalAmount
	var coPayAmount int64
	var memberResponsibility int64

	for _, b := range benefits {
		// Check annual limit
		used, _ := s.claimRepo.GetApprovedAmountForBenefitThisYear(ctx, claim.MemberID, b.Category)
		remaining := b.AnnualLimit - used

		if remaining <= 0 {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "annual_limit", Result: string(shared.RuleResultFail),
				Details: "Annual limit exhausted for " + b.Category,
			})
			continue
		}

		if payableAmount > remaining {
			payableAmount = remaining
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "annual_limit", Result: string(shared.RuleResultPass),
				Details: "Partial approval: limited to remaining balance",
			})
		} else {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "annual_limit", Result: string(shared.RuleResultPass),
				Details: "Within annual limit",
			})
		}

		// Apply co-pay
		if b.CoPayType == string(shared.CoPayTypePercentage) {
			coPayAmount = payableAmount * b.CoPayValue / 100
		} else if b.CoPayType == string(shared.CoPayTypeFixed) {
			coPayAmount = b.CoPayValue
		}

		payableAmount -= coPayAmount
		break
	}

	// 8. Fraud checks
	isDuplicate, _ := s.fraudSvc.CheckDuplicate(ctx, claim.ClaimNumber, claim.ID)
	if isDuplicate {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryFraud), Rule: "duplicate_check", Result: string(shared.RuleResultFlag),
			Details: "Potential duplicate claim detected",
		})
		return &entity.AdjudicationResult{
			Decision: string(shared.AdjudicationDecisionManualReview), PayableAmount: payableAmount,
			MemberResponsibility: coPayAmount + (claim.TotalAmount - payableAmount - coPayAmount),
			Reasons:              ruleResults,
		}, nil
	}

	ruleResults = append(ruleResults, entity.RuleResult{
		Category: string(shared.RuleCategoryFraud), Rule: "duplicate_check", Result: string(shared.RuleResultPass),
		Details: "No duplicate found",
	})

	memberResponsibility = coPayAmount + (claim.TotalAmount - payableAmount - coPayAmount)

	return &entity.AdjudicationResult{
		Decision:             string(shared.AdjudicationDecisionApprove),
		PayableAmount:        payableAmount,
		MemberResponsibility: memberResponsibility,
		Reasons:              ruleResults,
	}, nil
}

func parseJSONStringArray(data json.RawMessage) []string {
	var result []string
	if data == nil {
		return result
	}
	_ = json.Unmarshal(data, &result)
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func calculateAge(dob time.Time, asOf time.Time) int {
	age := asOf.Year() - dob.Year()
	if asOf.YearDay() < dob.YearDay() {
		age--
	}
	return age
}

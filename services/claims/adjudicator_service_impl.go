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
	preauthRepo "github.com/bitbiz/hias-core/domains/preauth/repository"
	productEntity "github.com/bitbiz/hias-core/domains/product/entity"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
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
	contractRepo  providerRepo.ContractRepository
	preauthRepo   preauthRepo.PreAuthRepository
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
	contractRepo providerRepo.ContractRepository,
	preauthRepo preauthRepo.PreAuthRepository,
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
		contractRepo:  contractRepo,
		preauthRepo:   preauthRepo,
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

	// 3a. Provider accreditation check
	if provider.AccreditationStatus != "" && provider.AccreditationStatus != string(shared.AccreditationStatusNone) {
		if provider.AccreditationStatus != string(shared.AccreditationStatusAccredited) {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "provider_accreditation", Result: string(shared.RuleResultFlag),
				Details: fmt.Sprintf("Provider accreditation status is %s (not ACCREDITED)", provider.AccreditationStatus),
			})
		} else {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "provider_accreditation", Result: string(shared.RuleResultPass),
				Details: "Provider is accredited",
			})
		}
	}

	// 3b. Provider network eligibility check
	if s.networkRepo != nil {
		inNetwork, err := s.networkRepo.CheckEligibility(ctx, policy.PlanID, claim.ProviderID, determineBenefitCategory(claim))
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

	// 3c. Contract validity check
	if s.contractRepo != nil {
		contracts, cErr := s.contractRepo.ListByProvider(ctx, claim.ProviderID)
		if cErr == nil {
			hasValidContract := false
			for _, c := range contracts {
				if c.Status == string(shared.ContractStatusActive) &&
					claim.ServiceDate.After(c.StartDate) && claim.ServiceDate.Before(c.EndDate) {
					hasValidContract = true
					break
				}
			}
			if !hasValidContract && len(contracts) > 0 {
				ruleResults = append(ruleResults, entity.RuleResult{
					Category: string(shared.RuleCategoryEligibility), Rule: "contract_valid", Result: string(shared.RuleResultFail),
					Details: "Provider has no valid contract for the service date",
				})
				return &entity.AdjudicationResult{
					Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
					MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
				}, nil
			}
			if hasValidContract {
				ruleResults = append(ruleResults, entity.RuleResult{
					Category: string(shared.RuleCategoryEligibility), Rule: "contract_valid", Result: string(shared.RuleResultPass),
					Details: "Provider has valid contract",
				})
			}
		}
	}

	// 3d. PreAuth validation (if preauth_id provided)
	if claim.PreAuthID != uuid.Nil && s.preauthRepo != nil {
		preauth, paErr := s.preauthRepo.GetByID(ctx, claim.PreAuthID)
		if paErr != nil {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "preauth_valid", Result: string(shared.RuleResultFail),
				Details: "Pre-authorization not found",
			})
			return &entity.AdjudicationResult{
				Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
				MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
			}, nil
		}
		if preauth.Status != string(shared.PreAuthStatusApproved) {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "preauth_valid", Result: string(shared.RuleResultFail),
				Details: fmt.Sprintf("Pre-authorization is %s, not APPROVED", preauth.Status),
			})
			return &entity.AdjudicationResult{
				Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
				MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
			}, nil
		}
		if preauth.ValidityEnd != nil && time.Now().After(*preauth.ValidityEnd) {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "preauth_valid", Result: string(shared.RuleResultFail),
				Details: "Pre-authorization has expired",
			})
			return &entity.AdjudicationResult{
				Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
				MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
			}, nil
		}
		if preauth.ProviderID != claim.ProviderID {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility), Rule: "preauth_valid", Result: string(shared.RuleResultFail),
				Details: "Pre-authorization provider does not match claim provider",
			})
			return &entity.AdjudicationResult{
				Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
				MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
			}, nil
		}
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility), Rule: "preauth_valid", Result: string(shared.RuleResultPass),
			Details: "Pre-authorization is valid",
		})

		// Validate claim procedures match preauth procedures
		preauthProcedures := parseJSONStringArray(preauth.ProcedureCodes)
		if len(preauthProcedures) > 0 && len(lineItems) > 0 {
			for _, li := range lineItems {
				if !contains(preauthProcedures, li.ProcedureCode) {
					ruleResults = append(ruleResults, entity.RuleResult{
						Category: string(shared.RuleCategoryEligibility), Rule: "preauth_procedure",
						Result:  string(shared.RuleResultFlag),
						Details: fmt.Sprintf("Procedure %s not in pre-authorization (authorized: %v)", li.ProcedureCode, preauthProcedures),
					})
				}
			}
		}

		// Warn if claim amount exceeds preauth approved amount
		if preauth.ApprovedAmount > 0 && claim.TotalAmount > preauth.ApprovedAmount {
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "preauth_amount_warning",
				Result:  string(shared.RuleResultFlag),
				Details: fmt.Sprintf("Claim amount %d exceeds PreAuth approved amount %d", claim.TotalAmount, preauth.ApprovedAmount),
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

	// 7. Calculate payable amount based on matching benefit
	payableAmount := claim.TotalAmount
	var coPayAmount int64
	var deductibleAmount int64
	var memberResponsibility int64

	claimCategory := determineBenefitCategory(claim)
	matchedBenefit := findMatchingBenefit(benefits, claimCategory)

	if matchedBenefit == nil {
		// Fallback: try first benefit
		matchedBenefit = benefits[0]
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryCoverage), Rule: "benefit_match",
			Result:  string(shared.RuleResultPass),
			Details: fmt.Sprintf("No exact category match for %s, using default benefit %s", claimCategory, matchedBenefit.Category),
		})
	} else {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryCoverage), Rule: "benefit_match",
			Result:  string(shared.RuleResultPass),
			Details: fmt.Sprintf("Matched benefit category: %s", matchedBenefit.Category),
		})
	}

	// Check annual limit
	used, _ := s.claimRepo.GetApprovedAmountForBenefitThisYear(ctx, claim.MemberID, matchedBenefit.Category)
	remaining := matchedBenefit.AnnualLimit - used

	if remaining <= 0 {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryLimits), Rule: "annual_limit",
			Result:  string(shared.RuleResultFail),
			Details: fmt.Sprintf("Annual limit exhausted for %s (used %d of %d)", matchedBenefit.Category, used, matchedBenefit.AnnualLimit),
		})
		return &entity.AdjudicationResult{
			Decision: string(shared.AdjudicationDecisionReject), PayableAmount: 0,
			MemberResponsibility: claim.TotalAmount, Reasons: ruleResults,
		}, nil
	}

	if payableAmount > remaining {
		payableAmount = remaining
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryLimits), Rule: "annual_limit",
			Result:  string(shared.RuleResultPass),
			Details: fmt.Sprintf("Partial approval: limited to remaining %d of %d", remaining, matchedBenefit.AnnualLimit),
		})
	} else {
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryLimits), Rule: "annual_limit",
			Result:  string(shared.RuleResultPass),
			Details: "Within annual limit",
		})
	}

	// Sub-limit enforcement
	if matchedBenefit.SubLimitType == string(shared.SubLimitTypePerVisit) && matchedBenefit.SubLimitValue > 0 {
		if payableAmount > matchedBenefit.SubLimitValue {
			payableAmount = matchedBenefit.SubLimitValue
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "sub_limit",
				Result:  string(shared.RuleResultPass),
				Details: fmt.Sprintf("Per-visit sub-limit applied: capped at %d", matchedBenefit.SubLimitValue),
			})
		}
	} else if matchedBenefit.SubLimitType == string(shared.SubLimitTypePerItem) && matchedBenefit.SubLimitValue > 0 {
		if payableAmount > matchedBenefit.SubLimitValue {
			payableAmount = matchedBenefit.SubLimitValue
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "sub_limit",
				Result:  string(shared.RuleResultPass),
				Details: fmt.Sprintf("Per-item sub-limit applied: capped at %d", matchedBenefit.SubLimitValue),
			})
		}
	}

	// Apply deductible
	deductibleAmount = matchedBenefit.DeductibleAmount
	if deductibleAmount > 0 {
		payableAmount -= deductibleAmount
		if payableAmount < 0 {
			payableAmount = 0
		}
		ruleResults = append(ruleResults, entity.RuleResult{
			Category: string(shared.RuleCategoryLimits), Rule: "deductible",
			Result:  string(shared.RuleResultPass),
			Details: fmt.Sprintf("Deductible of %d applied", deductibleAmount),
		})
	}

	// Apply co-pay
	if matchedBenefit.CoPayType == string(shared.CoPayTypePercentage) {
		coPayAmount = payableAmount * matchedBenefit.CoPayValue / 100
	} else if matchedBenefit.CoPayType == string(shared.CoPayTypeFixed) {
		coPayAmount = matchedBenefit.CoPayValue
	}
	payableAmount -= coPayAmount

	// 7b. Cap payable at PreAuth approved amount
	if claim.PreAuthID != uuid.Nil && s.preauthRepo != nil {
		preauth, paErr := s.preauthRepo.GetByID(ctx, claim.PreAuthID)
		if paErr == nil && preauth.ApprovedAmount > 0 && payableAmount > preauth.ApprovedAmount {
			payableAmount = preauth.ApprovedAmount
			ruleResults = append(ruleResults, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits), Rule: "preauth_cap", Result: string(shared.RuleResultPass),
				Details: fmt.Sprintf("Payable capped at PreAuth approved amount %d", preauth.ApprovedAmount),
			})
		}
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
			DeductibleApplied:    deductibleAmount,
			CoPayApplied:         coPayAmount,
			BenefitCategory:      claimCategory,
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
		DeductibleApplied:    deductibleAmount,
		CoPayApplied:         coPayAmount,
		BenefitCategory:      claimCategory,
		Reasons:              ruleResults,
	}, nil
}

func determineBenefitCategory(claim *entity.Claim) string {
	if claim.AdmissionDate != nil {
		return string(shared.BenefitCategoryInpatient)
	}
	return string(shared.BenefitCategoryOutpatient)
}

func findMatchingBenefit(benefits []*productEntity.Benefit, category string) *productEntity.Benefit {
	for _, b := range benefits {
		if b.Category == category {
			return b
		}
	}
	return nil
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

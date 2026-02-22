package claims

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimsRepository "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/claims/service"
	policyEntity "github.com/bitbiz/hias-core/domains/policy/entity"
	policyDomainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	productEntity "github.com/bitbiz/hias-core/domains/product/entity"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
)

// adjudicatorServiceImpl implements the core claims adjudication engine.
// It runs a sequential pipeline of rules across four categories:
//
//  1. Eligibility  -- Is the policy/member eligible to claim?
//  2. Coverage     -- Does the plan cover the procedures claimed?
//  3. Limits       -- Does the claim exceed annual or benefit limits?
//  4. Fraud        -- Are there signs of duplicate, frequency, or amount anomalies?
//
// The final decision is:
//   - REJECT  if any eligibility or coverage rule FAILs
//   - MANUAL_REVIEW  if any fraud rule FLAGs the claim
//   - APPROVE  if all rules PASS, with a calculated payable amount
type adjudicatorServiceImpl struct {
	claimRepo     claimsRepository.ClaimRepository
	benefitRepo   productRepo.BenefitRepository
	exclusionRepo productRepo.ExclusionRepository
	policyRepo    policyDomainRepo.PolicyRepository
	memberRepo    policyDomainRepo.MemberRepository
	fraudService  service.FraudService
}

func NewAdjudicatorService(
	claimRepo claimsRepository.ClaimRepository,
	benefitRepo productRepo.BenefitRepository,
	exclusionRepo productRepo.ExclusionRepository,
	policyRepo policyDomainRepo.PolicyRepository,
	memberRepo policyDomainRepo.MemberRepository,
	fraudService service.FraudService,
) service.AdjudicatorService {
	return &adjudicatorServiceImpl{
		claimRepo:     claimRepo,
		benefitRepo:   benefitRepo,
		exclusionRepo: exclusionRepo,
		policyRepo:    policyRepo,
		memberRepo:    memberRepo,
		fraudService:  fraudService,
	}
}

// Adjudicate runs the full adjudication pipeline against a claim and its
// line items. It returns an AdjudicationResult containing the decision,
// payable amount, member responsibility, and the detailed rule results.
func (a *adjudicatorServiceImpl) Adjudicate(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) (*entity.AdjudicationResult, error) {
	var allResults []entity.RuleResult

	utils.LogInfo("Adjudicator: starting adjudication for claim %s (amount: %d cents)", claim.ClaimNumber, claim.TotalAmount)

	// ====================================================================
	// STEP 1: ELIGIBILITY RULES
	// ====================================================================
	eligibilityResults, policy, err := a.runEligibilityRules(ctx, claim)
	if err != nil {
		return nil, fmt.Errorf("eligibility check failed: %w", err)
	}
	allResults = append(allResults, eligibilityResults...)

	// If any eligibility rule failed, reject immediately
	if hasFailure(eligibilityResults) {
		utils.LogInfo("Adjudicator: claim %s REJECTED at eligibility stage", claim.ClaimNumber)
		return &entity.AdjudicationResult{
			Decision:             string(shared.AdjudicationDecisionReject),
			PayableAmount:        0,
			MemberResponsibility: claim.TotalAmount,
			Reasons:              allResults,
		}, nil
	}

	// ====================================================================
	// STEP 2: COVERAGE RULES
	// ====================================================================
	coverageResults, matchedBenefits, err := a.runCoverageRules(ctx, claim, lineItems, policy)
	if err != nil {
		return nil, fmt.Errorf("coverage check failed: %w", err)
	}
	allResults = append(allResults, coverageResults...)

	// If any coverage rule failed, reject
	if hasFailure(coverageResults) {
		utils.LogInfo("Adjudicator: claim %s REJECTED at coverage stage", claim.ClaimNumber)
		return &entity.AdjudicationResult{
			Decision:             string(shared.AdjudicationDecisionReject),
			PayableAmount:        0,
			MemberResponsibility: claim.TotalAmount,
			Reasons:              allResults,
		}, nil
	}

	// ====================================================================
	// STEP 3: LIMITS RULES (also calculates payable amounts)
	// ====================================================================
	limitsResults, payableAmount, memberResponsibility, err := a.runLimitsRules(ctx, claim, lineItems, matchedBenefits)
	if err != nil {
		return nil, fmt.Errorf("limits check failed: %w", err)
	}
	allResults = append(allResults, limitsResults...)

	// ====================================================================
	// STEP 4: FRAUD RULES
	// ====================================================================
	fraudResults, err := a.runFraudRules(ctx, claim, lineItems)
	if err != nil {
		return nil, fmt.Errorf("fraud check failed: %w", err)
	}
	allResults = append(allResults, fraudResults...)

	// ====================================================================
	// STEP 5: FINAL DECISION
	// ====================================================================
	decision := a.determineDecision(allResults, payableAmount, memberResponsibility, claim)

	utils.LogInfo("Adjudicator: claim %s decision=%s payable=%d member_resp=%d (%d rules evaluated)",
		claim.ClaimNumber, decision.Decision, decision.PayableAmount, decision.MemberResponsibility, len(allResults))

	return decision, nil
}

// ---------------------------------------------------------------------------
// Step 1: Eligibility Rules
// ---------------------------------------------------------------------------

func (a *adjudicatorServiceImpl) runEligibilityRules(ctx context.Context, claim *entity.Claim) ([]entity.RuleResult, *policyEntity.Policy, error) {
	var results []entity.RuleResult

	// Rule 1.1: Policy exists and is ACTIVE
	policy, err := a.policyRepo.GetByID(ctx, claim.PolicyID)
	if err != nil {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility),
			Rule:     "policy_exists",
			Result:   string(shared.RuleResultFail),
			Details:  fmt.Sprintf("Policy %s not found", claim.PolicyID),
		})
		return results, nil, nil
	}

	if policy.Status != string(shared.PolicyStatusActive) {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility),
			Rule:     "policy_active",
			Result:   string(shared.RuleResultFail),
			Details:  fmt.Sprintf("Policy status is %s, expected ACTIVE", policy.Status),
		})
	} else {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility),
			Rule:     "policy_active",
			Result:   string(shared.RuleResultPass),
			Details:  "Policy is ACTIVE",
		})
	}

	// Rule 1.2: Member exists and is enrolled (not terminated)
	member, err := a.memberRepo.GetByID(ctx, claim.MemberID)
	if err != nil {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryEligibility),
			Rule:     "member_enrolled",
			Result:   string(shared.RuleResultFail),
			Details:  fmt.Sprintf("Member %s not found", claim.MemberID),
		})
	} else {
		// Check member belongs to the claimed policy
		if member.PolicyID != claim.PolicyID {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility),
				Rule:     "member_enrolled",
				Result:   string(shared.RuleResultFail),
				Details:  fmt.Sprintf("Member does not belong to policy %s", claim.PolicyID),
			})
		} else {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility),
				Rule:     "member_enrolled",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("Member %s is enrolled under policy %s", member.Name, claim.PolicyID),
			})
		}
	}

	// Rule 1.3: Policy covers the service date (within start/end)
	if policy != nil {
		if claim.ServiceDate.Before(policy.StartDate) || claim.ServiceDate.After(policy.EndDate) {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility),
				Rule:     "service_date_in_policy_period",
				Result:   string(shared.RuleResultFail),
				Details:  fmt.Sprintf("Service date %s is outside policy period %s to %s", claim.ServiceDate.Format("2006-01-02"), policy.StartDate.Format("2006-01-02"), policy.EndDate.Format("2006-01-02")),
			})
		} else {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryEligibility),
				Rule:     "service_date_in_policy_period",
				Result:   string(shared.RuleResultPass),
				Details:  "Service date is within policy period",
			})
		}
	}

	return results, policy, nil
}

// ---------------------------------------------------------------------------
// Step 2: Coverage Rules
// ---------------------------------------------------------------------------

func (a *adjudicatorServiceImpl) runCoverageRules(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem, policy *policyEntity.Policy) ([]entity.RuleResult, map[string]*productEntity.Benefit, error) {
	var results []entity.RuleResult
	matchedBenefits := make(map[string]*productEntity.Benefit)

	if policy == nil {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryCoverage),
			Rule:     "policy_available",
			Result:   string(shared.RuleResultFail),
			Details:  "Cannot check coverage without a valid policy",
		})
		return results, matchedBenefits, nil
	}

	// Fetch all benefits for the plan
	benefits, err := a.benefitRepo.ListByPlan(ctx, policy.PlanID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch benefits for plan %s: %w", policy.PlanID, err)
	}

	// Build category lookup: category -> benefit
	benefitByCategory := make(map[string]*productEntity.Benefit)
	for _, b := range benefits {
		benefitByCategory[b.Category] = b
	}

	// Fetch exclusions for the plan
	exclusions, err := a.exclusionRepo.ListByPlan(ctx, policy.PlanID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch exclusions for plan %s: %w", policy.PlanID, err)
	}

	// Build set of excluded ICD codes
	excludedCodes := make(map[string]bool)
	for _, excl := range exclusions {
		var codes []string
		if err := json.Unmarshal(excl.ICDCodes, &codes); err == nil {
			for _, code := range codes {
				excludedCodes[code] = true
			}
		}
	}

	// Rule 2.1: Check each line item's procedure has a matching benefit category
	// For simplicity, we map procedure codes to benefit categories using a heuristic.
	// In a production system, this would be a lookup against a procedure-to-category table.
	for _, li := range lineItems {
		category := mapProcedureToCategory(li.ProcedureCode)
		benefit, exists := benefitByCategory[category]
		if !exists {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryCoverage),
				Rule:     "benefit_exists",
				Result:   string(shared.RuleResultFail),
				Details:  fmt.Sprintf("No benefit found for category '%s' (procedure %s) in the plan", category, li.ProcedureCode),
			})
		} else {
			// Check waiting period
			policyDays := int(claim.ServiceDate.Sub(policy.StartDate).Hours() / 24)
			if benefit.WaitingPeriodDays > 0 && policyDays < benefit.WaitingPeriodDays {
				results = append(results, entity.RuleResult{
					Category: string(shared.RuleCategoryEligibility),
					Rule:     "waiting_period",
					Result:   string(shared.RuleResultFail),
					Details:  fmt.Sprintf("Waiting period of %d days not met for category '%s' (enrolled %d days)", benefit.WaitingPeriodDays, category, policyDays),
				})
			} else {
				results = append(results, entity.RuleResult{
					Category: string(shared.RuleCategoryCoverage),
					Rule:     "benefit_exists",
					Result:   string(shared.RuleResultPass),
					Details:  fmt.Sprintf("Benefit '%s' (category %s) covers procedure %s", benefit.Name, category, li.ProcedureCode),
				})
				matchedBenefits[li.ProcedureCode] = benefit
			}
		}
	}

	// Rule 2.2: Check diagnosis codes are not excluded
	var diagnosisCodes []string
	if err := json.Unmarshal(claim.DiagnosisCodes, &diagnosisCodes); err == nil {
		for _, code := range diagnosisCodes {
			if excludedCodes[code] {
				results = append(results, entity.RuleResult{
					Category: string(shared.RuleCategoryCoverage),
					Rule:     "exclusion_check",
					Result:   string(shared.RuleResultFail),
					Details:  fmt.Sprintf("Diagnosis code %s is excluded from the plan", code),
				})
			} else {
				results = append(results, entity.RuleResult{
					Category: string(shared.RuleCategoryCoverage),
					Rule:     "exclusion_check",
					Result:   string(shared.RuleResultPass),
					Details:  fmt.Sprintf("Diagnosis code %s is not excluded", code),
				})
			}
		}
	}

	return results, matchedBenefits, nil
}

// ---------------------------------------------------------------------------
// Step 3: Limits Rules (calculates payable amounts)
// ---------------------------------------------------------------------------

func (a *adjudicatorServiceImpl) runLimitsRules(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem, matchedBenefits map[string]*productEntity.Benefit) ([]entity.RuleResult, int64, int64, error) {
	var results []entity.RuleResult
	var totalPayable int64
	var totalMemberResponsibility int64

	for _, li := range lineItems {
		lineTotal := li.TotalPrice
		benefit, hasBenefit := matchedBenefits[li.ProcedureCode]

		if !hasBenefit {
			// No benefit matched -- entire line item is member's responsibility
			totalMemberResponsibility += lineTotal
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits),
				Rule:     "no_benefit_match",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("No benefit for procedure %s; %d cents is member responsibility", li.ProcedureCode, lineTotal),
			})
			continue
		}

		// ----- Annual limit check -----
		approvedThisYear, err := a.claimRepo.GetApprovedAmountForBenefitThisYear(ctx, claim.MemberID, benefit.Category)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("failed to get approved amount for benefit %s: %w", benefit.Category, err)
		}

		remainingLimit := benefit.AnnualLimit - approvedThisYear
		if remainingLimit <= 0 {
			// Annual limit fully exhausted
			totalMemberResponsibility += lineTotal
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits),
				Rule:     "annual_limit",
				Result:   string(shared.RuleResultFail),
				Details:  fmt.Sprintf("Annual limit of %d cents for '%s' is fully exhausted (used: %d cents)", benefit.AnnualLimit, benefit.Category, approvedThisYear),
			})
			continue
		}

		// Determine claimable amount (capped at remaining limit)
		claimableAmount := lineTotal
		partialApproval := false
		if claimableAmount > remainingLimit {
			claimableAmount = remainingLimit
			partialApproval = true
		}

		if partialApproval {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits),
				Rule:     "annual_limit",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("Partial approval: %d of %d cents covered (remaining limit: %d cents for '%s')", claimableAmount, lineTotal, remainingLimit, benefit.Category),
			})
		} else {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryLimits),
				Rule:     "annual_limit",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("Within annual limit for '%s' (remaining: %d cents, claiming: %d cents)", benefit.Category, remainingLimit, claimableAmount),
			})
		}

		// ----- CoPay calculation -----
		var coPayAmount int64
		switch benefit.CoPayType {
		case string(shared.CoPayTypePercentage):
			// CoPayValue is percentage (e.g., 20 means 20%)
			coPayAmount = (claimableAmount * benefit.CoPayValue) / 100
		case string(shared.CoPayTypeFixed):
			// CoPayValue is a fixed amount in cents
			coPayAmount = benefit.CoPayValue
			if coPayAmount > claimableAmount {
				coPayAmount = claimableAmount
			}
		default:
			// No copay
			coPayAmount = 0
		}

		payable := claimableAmount - coPayAmount
		memberResp := lineTotal - payable

		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryLimits),
			Rule:     "copay_calculation",
			Result:   string(shared.RuleResultPass),
			Details:  fmt.Sprintf("CoPay (%s=%d): insurer pays %d cents, member pays %d cents for procedure %s", benefit.CoPayType, benefit.CoPayValue, payable, memberResp, li.ProcedureCode),
		})

		totalPayable += payable
		totalMemberResponsibility += memberResp
	}

	return results, totalPayable, totalMemberResponsibility, nil
}

// ---------------------------------------------------------------------------
// Step 4: Fraud Rules
// ---------------------------------------------------------------------------

func (a *adjudicatorServiceImpl) runFraudRules(ctx context.Context, claim *entity.Claim, lineItems []*entity.ClaimLineItem) ([]entity.RuleResult, error) {
	var results []entity.RuleResult

	// Rule 4.1: Duplicate check
	isDuplicate, err := a.fraudService.CheckDuplicate(ctx, claim.ClaimNumber, claim.ID)
	if err != nil {
		return nil, fmt.Errorf("duplicate check failed: %w", err)
	}
	if isDuplicate {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryFraud),
			Rule:     "duplicate_check",
			Result:   string(shared.RuleResultFlag),
			Details:  fmt.Sprintf("Potential duplicate claim detected for claim %s", claim.ClaimNumber),
		})
		// Flag the claim
		_ = a.fraudService.FlagClaim(ctx, &entity.FraudFlag{
			ClaimID:  claim.ID,
			FlagType: string(shared.FraudFlagDuplicate),
			Severity: string(shared.FraudSeverityHigh),
			Details:  "Duplicate claim number or matching member+provider+service_date found",
		})
	} else {
		results = append(results, entity.RuleResult{
			Category: string(shared.RuleCategoryFraud),
			Rule:     "duplicate_check",
			Result:   string(shared.RuleResultPass),
			Details:  "No duplicate claims detected",
		})
	}

	// Rule 4.2: Frequency check (same member+provider+procedure within 7 days)
	for _, li := range lineItems {
		isFrequent, err := a.fraudService.CheckFrequency(ctx, claim.MemberID, claim.ProviderID, li.ProcedureCode, claim.ID)
		if err != nil {
			return nil, fmt.Errorf("frequency check failed for procedure %s: %w", li.ProcedureCode, err)
		}
		if isFrequent {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryFraud),
				Rule:     "frequency_check",
				Result:   string(shared.RuleResultFlag),
				Details:  fmt.Sprintf("Suspicious frequency: procedure %s claimed by same member at same provider within 7 days", li.ProcedureCode),
			})
			_ = a.fraudService.FlagClaim(ctx, &entity.FraudFlag{
				ClaimID:  claim.ID,
				FlagType: string(shared.FraudFlagFrequency),
				Severity: string(shared.FraudSeverityMedium),
				Details:  fmt.Sprintf("Same procedure %s claimed within 7 days for member %s at provider %s", li.ProcedureCode, claim.MemberID, claim.ProviderID),
			})
		} else {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryFraud),
				Rule:     "frequency_check",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("No suspicious frequency for procedure %s", li.ProcedureCode),
			})
		}
	}

	// Rule 4.3: Amount threshold check (> 2x rate card for the procedure)
	for _, li := range lineItems {
		exceeds, err := a.fraudService.CheckAmountThreshold(ctx, claim.ProviderID, li.ProcedureCode, li.TotalPrice)
		if err != nil {
			// Non-fatal: rate card may not exist for all procedures
			utils.LogWarn("Adjudicator: amount threshold check skipped for procedure %s: %v", li.ProcedureCode, err)
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryFraud),
				Rule:     "amount_threshold",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("No rate card found for procedure %s at provider; threshold check skipped", li.ProcedureCode),
			})
			continue
		}
		if exceeds {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryFraud),
				Rule:     "amount_threshold",
				Result:   string(shared.RuleResultFlag),
				Details:  fmt.Sprintf("Amount %d cents for procedure %s exceeds 2x the rate card", li.TotalPrice, li.ProcedureCode),
			})
			_ = a.fraudService.FlagClaim(ctx, &entity.FraudFlag{
				ClaimID:  claim.ID,
				FlagType: string(shared.FraudFlagAmountThreshold),
				Severity: string(shared.FraudSeverityHigh),
				Details:  fmt.Sprintf("Claimed amount %d cents for procedure %s exceeds 2x rate card at provider %s", li.TotalPrice, li.ProcedureCode, claim.ProviderID),
			})
		} else {
			results = append(results, entity.RuleResult{
				Category: string(shared.RuleCategoryFraud),
				Rule:     "amount_threshold",
				Result:   string(shared.RuleResultPass),
				Details:  fmt.Sprintf("Amount for procedure %s is within acceptable range", li.ProcedureCode),
			})
		}
	}

	return results, nil
}

// ---------------------------------------------------------------------------
// Step 5: Final Decision
// ---------------------------------------------------------------------------

func (a *adjudicatorServiceImpl) determineDecision(allResults []entity.RuleResult, payableAmount, memberResponsibility int64, claim *entity.Claim) *entity.AdjudicationResult {
	// Any FAIL in eligibility or coverage -> REJECT
	for _, r := range allResults {
		if r.Result == string(shared.RuleResultFail) &&
			(r.Category == string(shared.RuleCategoryEligibility) || r.Category == string(shared.RuleCategoryCoverage)) {
			return &entity.AdjudicationResult{
				Decision:             string(shared.AdjudicationDecisionReject),
				PayableAmount:        0,
				MemberResponsibility: claim.TotalAmount,
				Reasons:              allResults,
			}
		}
	}

	// Any FLAG from fraud -> MANUAL_REVIEW
	for _, r := range allResults {
		if r.Result == string(shared.RuleResultFlag) && r.Category == string(shared.RuleCategoryFraud) {
			return &entity.AdjudicationResult{
				Decision:             string(shared.AdjudicationDecisionManualReview),
				PayableAmount:        payableAmount,
				MemberResponsibility: memberResponsibility,
				Reasons:              allResults,
			}
		}
	}

	// All PASS -> APPROVE
	return &entity.AdjudicationResult{
		Decision:             string(shared.AdjudicationDecisionApprove),
		PayableAmount:        payableAmount,
		MemberResponsibility: memberResponsibility,
		Reasons:              allResults,
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// hasFailure returns true if any rule result in the slice has a FAIL status.
func hasFailure(results []entity.RuleResult) bool {
	for _, r := range results {
		if r.Result == string(shared.RuleResultFail) {
			return true
		}
	}
	return false
}

// mapProcedureToCategory maps a procedure code prefix to a benefit category.
// This is a simplified heuristic -- in production, a procedure-to-category
// mapping table or external service would be used.
//
// Convention used:
//
//	OPD-xxx  -> outpatient
//	IPD-xxx  -> inpatient
//	DEN-xxx  -> dental
//	OPT-xxx  -> optical
//	MAT-xxx  -> maternity
//	(default) -> outpatient
func mapProcedureToCategory(procedureCode string) string {
	if len(procedureCode) < 3 {
		return string(shared.BenefitCategoryOutpatient)
	}
	prefix := procedureCode[:3]
	switch prefix {
	case "OPD":
		return string(shared.BenefitCategoryOutpatient)
	case "IPD":
		return string(shared.BenefitCategoryInpatient)
	case "DEN":
		return string(shared.BenefitCategoryDental)
	case "OPT":
		return string(shared.BenefitCategoryOptical)
	case "MAT":
		return string(shared.BenefitCategoryMaternity)
	default:
		return string(shared.BenefitCategoryOutpatient)
	}
}

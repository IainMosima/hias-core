package reporting

import (
	"context"
	"fmt"
	"time"

	domainRepo "github.com/bitbiz/hias-core/domains/reporting/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type reportDataRepository struct {
	store db.Store
}

func NewReportDataRepository(store db.Store) domainRepo.ReportDataRepository {
	return &reportDataRepository{store: store}
}

func (r *reportDataRepository) GetClaimsExperienceData(ctx context.Context, start, end time.Time, policyID *uuid.UUID) ([]map[string]interface{}, error) {
	var pid pgtype.UUID
	if policyID != nil {
		pid = pgtype.UUID{Bytes: *policyID, Valid: true}
	}
	rows, err := r.store.GetClaimsExperienceData(ctx, db.GetClaimsExperienceDataParams{
		StartDate: start,
		EndDate:   end,
		PolicyID:  pid,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get claims experience data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		tatFloat, _ := row.AvgTatHours.Float64Value()
		result[i] = map[string]interface{}{
			"policy_number":     row.PolicyNumber,
			"policyholder_name": row.PolicyholderName,
			"total_premium":     row.TotalPremium,
			"total_claims":      row.TotalClaims,
			"approved_claims":   row.ApprovedClaims,
			"rejected_claims":   row.RejectedClaims,
			"loss_ratio":        float64(row.LossRatio),
			"claim_count":       row.ClaimCount,
			"avg_tat_hours":     tatFloat.Float64,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetClaimsRegisterData(ctx context.Context, start, end time.Time, status string, limit, offset int) ([]map[string]interface{}, error) {
	rows, err := r.store.GetClaimsRegisterData(ctx, db.GetClaimsRegisterDataParams{
		Limit:       int32(limit),
		Offset:      int32(offset),
		StartDate:   start,
		EndDate:     end,
		ClaimStatus: pgtype.Text{String: status, Valid: status != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get claims register data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		result[i] = map[string]interface{}{
			"claim_number":    row.ClaimNumber,
			"policy_number":   row.PolicyNumber,
			"member_name":     row.MemberName,
			"provider_name":   row.ProviderName,
			"claim_type":      row.ClaimType,
			"service_date":    row.ServiceDate,
			"total_amount":    row.TotalAmount,
			"approved_amount": row.ApprovedAmount,
			"status":          row.Status,
			"created_at":      row.CreatedAt,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetPremiumDebtorsAgeingData(ctx context.Context) ([]map[string]interface{}, error) {
	rows, err := r.store.GetPremiumDebtorsAgeingData(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get premium debtors ageing data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		result[i] = map[string]interface{}{
			"policy_number":     row.PolicyNumber,
			"policyholder_name": row.PolicyholderName,
			"total_premium":     row.TotalPremium,
			"total_paid":        row.TotalPaid,
			"outstanding":       row.Outstanding,
			"current_bucket":    row.CurrentBucket,
			"days_30":           row.Days30,
			"days_60":           row.Days60,
			"days_90_plus":      row.Days90Plus,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetPremiumRegisterData(ctx context.Context, start, end time.Time, limit, offset int) ([]map[string]interface{}, error) {
	rows, err := r.store.GetPremiumRegisterData(ctx, db.GetPremiumRegisterDataParams{
		Limit:     int32(limit),
		Offset:    int32(offset),
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get premium register data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		var paymentDate interface{}
		if row.PaymentDate.Valid {
			paymentDate = row.PaymentDate.Time
		}
		result[i] = map[string]interface{}{
			"policy_number":     row.PolicyNumber,
			"policyholder_name": row.PolicyholderName,
			"plan_name":         row.PlanName,
			"premium_amount":    row.PremiumAmount,
			"payment_amount":    row.PaymentAmount,
			"payment_date":      paymentDate,
			"payment_status":    row.PaymentStatus,
			"payment_method":    row.PaymentMethod,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetMembershipData(ctx context.Context, policyID *uuid.UUID, status string, limit, offset int) ([]map[string]interface{}, error) {
	var pid pgtype.UUID
	if policyID != nil {
		pid = pgtype.UUID{Bytes: *policyID, Valid: true}
	}
	rows, err := r.store.GetMembershipData(ctx, db.GetMembershipDataParams{
		Limit:        int32(limit),
		Offset:       int32(offset),
		PolicyID:     pid,
		MemberStatus: pgtype.Text{String: status, Valid: status != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get membership data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		var dob interface{}
		if row.DateOfBirth.Valid {
			dob = row.DateOfBirth.Time
		}
		result[i] = map[string]interface{}{
			"member_number":   row.MemberNumber,
			"full_name":       row.FullName,
			"date_of_birth":   dob,
			"gender":          row.Gender,
			"phone":           pgtypeTextToString(row.Phone),
			"email":           pgtypeTextToString(row.Email),
			"relationship":    row.Relationship,
			"policy_number":   row.PolicyNumber,
			"plan_name":       row.PlanName,
			"status":          row.Status,
			"enrollment_date": row.EnrollmentDate,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetProviderPerformanceData(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error) {
	rows, err := r.store.GetProviderPerformanceData(ctx, db.GetProviderPerformanceDataParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get provider performance data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		result[i] = map[string]interface{}{
			"provider_name":    row.ProviderName,
			"tier":             row.Tier,
			"status":           row.Status,
			"total_claims":     row.TotalClaims,
			"total_claimed":    row.TotalClaimed,
			"total_approved":   row.TotalApproved,
			"rejection_rate":   float64(row.RejectionRate),
			"avg_claim_amount": int64(row.AvgClaimAmount),
			"fraud_flag_count": row.FraudFlagCount,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetLossRatioData(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error) {
	rows, err := r.store.GetLossRatioData(ctx, db.GetLossRatioDataParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get loss ratio data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		expenseFloat, _ := row.ExpenseRatio.Float64Value()
		result[i] = map[string]interface{}{
			"plan_name":       row.PlanName,
			"active_policies": row.ActivePolicies,
			"total_members":   row.TotalMembers,
			"earned_premium":  row.EarnedPremium,
			"incurred_claims": row.IncurredClaims,
			"loss_ratio":      float64(row.LossRatio),
			"expense_ratio":   expenseFloat.Float64,
			"combined_ratio":  float64(row.CombinedRatio),
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetRenewalData(ctx context.Context, start, end time.Time) ([]map[string]interface{}, error) {
	rows, err := r.store.GetRenewalData(ctx, db.GetRenewalDataParams{
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get renewal data: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		var expiryDate interface{}
		if row.ExpiryDate.Valid {
			expiryDate = row.ExpiryDate.Time
		}
		result[i] = map[string]interface{}{
			"policy_number":      row.PolicyNumber,
			"policyholder_name":  row.PolicyholderName,
			"plan_name":          row.PlanName,
			"expiry_date":        expiryDate,
			"renewal_status":     row.RenewalStatus,
			"current_premium":    row.CurrentPremium,
			"proposed_premium":   row.ProposedPremium,
			"premium_change_pct": float64(row.PremiumChangePct),
			"member_count":       row.MemberCount,
		}
	}
	return result, nil
}

func (r *reportDataRepository) DrillDownClaimsByPolicy(ctx context.Context, policyID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error) {
	rows, err := r.store.DrillDownClaimsByPolicy(ctx, db.DrillDownClaimsByPolicyParams{
		PolicyID:  policyID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to drill down claims by policy: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		result[i] = map[string]interface{}{
			"claim_number":    row.ClaimNumber,
			"member_name":     row.MemberName,
			"provider_name":   row.ProviderName,
			"claim_type":      row.ClaimType,
			"service_date":    row.ServiceDate,
			"total_amount":    row.TotalAmount,
			"approved_amount": row.ApprovedAmount,
			"co_pay_amount":   row.CoPayAmount,
			"status":          row.Status,
			"created_at":      row.CreatedAt,
		}
	}
	return result, nil
}

func (r *reportDataRepository) DrillDownPaymentsByPolicy(ctx context.Context, policyID uuid.UUID, start, end time.Time) ([]map[string]interface{}, error) {
	rows, err := r.store.DrillDownPaymentsByPolicy(ctx, db.DrillDownPaymentsByPolicyParams{
		PolicyID:  policyID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to drill down payments by policy: %w", err)
	}
	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		var paidAt interface{}
		if row.PaidAt.Valid {
			paidAt = row.PaidAt.Time
		}
		result[i] = map[string]interface{}{
			"reference_number": pgtypeTextToString(row.ReferenceNumber),
			"amount":           row.Amount,
			"method":           row.Method,
			"status":           row.Status,
			"paid_at":          paidAt,
			"invoice_number":   row.InvoiceNumber,
			"invoice_amount":   row.InvoiceAmount,
			"due_date":         row.DueDate,
		}
	}
	return result, nil
}

func (r *reportDataRepository) GetOutstandingPremium(ctx context.Context) (int64, error) {
	amount, err := r.store.GetOutstandingPremium(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get outstanding premium: %w", err)
	}
	return amount, nil
}

func (r *reportDataRepository) GetSLABreachCount(ctx context.Context) (int64, error) {
	count, err := r.store.GetSLABreachCount(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get SLA breach count: %w", err)
	}
	return count, nil
}

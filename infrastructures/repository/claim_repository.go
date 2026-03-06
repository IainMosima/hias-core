package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type claimRepository struct {
	store db.Store
}

func NewClaimRepository(store db.Store) domainRepo.ClaimRepository {
	return &claimRepository{store: store}
}

func (r *claimRepository) Create(ctx context.Context, claim *entity.Claim) (*entity.Claim, error) {
	dbClaim, err := r.store.CreateClaim(ctx, db.CreateClaimParams{
		ClaimNumber:    claim.ClaimNumber,
		PolicyID:       claim.PolicyID,
		MemberID:       claim.MemberID,
		ProviderID:     claim.ProviderID,
		PreauthID:      uuidToPgtype(claim.PreAuthID),
		Status:         claim.Status,
		TotalAmount:    claim.TotalAmount,
		DiagnosisCodes: claim.DiagnosisCodes,
		ServiceDate:    claim.ServiceDate,
		AdmissionDate:  timePtrToPgtypeTimestamptz(claim.AdmissionDate),
		DischargeDate:  timePtrToPgtypeTimestamptz(claim.DischargeDate),
		Notes:          claim.Notes,
		ClaimType:      claim.ClaimType,
		SlaBreachAt:    timePtrToPgtypeTimestamptz(claim.SLABreachAt),
		CreatedBy:      uuidToPgtype(claim.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create claim: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Claim, error) {
	dbClaim, err := r.store.GetClaimByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get claim by ID: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) GetByNumber(ctx context.Context, number string) (*entity.Claim, error) {
	dbClaim, err := r.store.GetClaimByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get claim by number: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) List(ctx context.Context, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListClaims(ctx, db.ListClaimsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list claims: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListClaimsByStatus(ctx, db.ListClaimsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list claims by status: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListClaimsByProvider(ctx, db.ListClaimsByProviderParams{
		ProviderID: providerID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list claims by provider: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListClaimsByMember(ctx, db.ListClaimsByMemberParams{
		MemberID: memberID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list claims by member: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListClaimsByPolicy(ctx, db.ListClaimsByPolicyParams{
		PolicyID: policyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list claims by policy: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountClaims(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count claims: %w", err)
	}
	return count, nil
}

func (r *claimRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	count, err := r.store.CountClaimsByStatus(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count claims by status: %w", err)
	}
	return count, nil
}

func (r *claimRepository) GetForAdjudication(ctx context.Context, limit int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.GetClaimsForAdjudication(ctx, int32(limit))
	if err != nil {
		return nil, fmt.Errorf("failed to get claims for adjudication: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) GetApprovedForRemittance(ctx context.Context, providerID uuid.UUID) ([]*entity.Claim, error) {
	dbClaims, err := r.store.GetApprovedClaimsForRemittance(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved claims for remittance: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Claim, error) {
	dbClaim, err := r.store.UpdateClaimStatus(ctx, db.UpdateClaimStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update claim status: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) UpdateAmounts(ctx context.Context, id uuid.UUID, approved, copay, memberResp int64) (*entity.Claim, error) {
	dbClaim, err := r.store.UpdateClaimAmounts(ctx, db.UpdateClaimAmountsParams{
		ID:                   id,
		ApprovedAmount:       approved,
		CoPayAmount:          copay,
		MemberResponsibility: memberResp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update claim amounts: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) Reject(ctx context.Context, id uuid.UUID, reason string) (*entity.Claim, error) {
	dbClaim, err := r.store.UpdateClaimRejection(ctx, db.UpdateClaimRejectionParams{
		ID:              id,
		RejectionReason: stringToPgtypeText(reason),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reject claim: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) GetApprovedAmountForBenefitThisYear(ctx context.Context, memberID uuid.UUID, category string) (int64, error) {
	amount, err := r.store.GetApprovedAmountForBenefitThisYear(ctx, db.GetApprovedAmountForBenefitThisYearParams{
		MemberID: memberID,
		Category: category,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get approved amount for benefit this year: %w", err)
	}
	return amount, nil
}

func (r *claimRepository) VetClaim(ctx context.Context, id uuid.UUID, vettedAmount int64, vettedBy uuid.UUID, status string) (*entity.Claim, error) {
	dbClaim, err := r.store.VetClaim(ctx, db.VetClaimParams{
		ID:           id,
		VettedAmount: int64ToPgtypeInt8(vettedAmount),
		VettedBy:     uuidToPgtype(vettedBy),
		Status:       status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to vet claim: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) MarkReadyForPayment(ctx context.Context, id uuid.UUID) (*entity.Claim, error) {
	dbClaim, err := r.store.MarkClaimReadyForPayment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to mark claim ready for payment: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) ListSLABreached(ctx context.Context, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListSLABreachedClaims(ctx, db.ListSLABreachedClaimsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list SLA breached claims: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) FindByProviderAndDate(ctx context.Context, providerID uuid.UUID, serviceDate time.Time, amount int64) (*entity.Claim, error) {
	dbClaim, err := r.store.FindClaimByProviderAndDate(ctx, db.FindClaimByProviderAndDateParams{
		ProviderID:  providerID,
		ServiceDate: serviceDate,
		TotalAmount: amount,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find claim by provider and date: %w", err)
	}
	return sqlcClaimToDomain(dbClaim), nil
}

func (r *claimRepository) ListApproachingSLA(ctx context.Context, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListApproachingSLAClaims(ctx, db.ListApproachingSLAClaimsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list approaching SLA claims: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) CountByMemberThisMonth(ctx context.Context, memberID uuid.UUID) (int64, error) {
	count, err := r.store.CountClaimsByMemberThisMonth(ctx, memberID)
	if err != nil {
		return 0, fmt.Errorf("failed to count claims by member this month: %w", err)
	}
	return count, nil
}

func (r *claimRepository) SetEscalatedTo(ctx context.Context, claimID uuid.UUID, role string) error {
	return r.store.SetClaimEscalatedTo(ctx, db.SetClaimEscalatedToParams{
		ID:          claimID,
		EscalatedTo: stringToPgtypeText(role),
	})
}

func (r *claimRepository) ListFiltered(ctx context.Context, status string, dateFrom, dateTo *time.Time, search string, limit, offset int) ([]*entity.Claim, error) {
	dbClaims, err := r.store.ListClaimsFiltered(ctx, db.ListClaimsFilteredParams{
		StatusFilter: status,
		DateFrom:     timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:       timePtrToPgtypeTimestamptz(dateTo),
		Search:       search,
		QueryLimit:   int32(limit),
		QueryOffset:  int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list claims filtered: %w", err)
	}
	claims := make([]*entity.Claim, len(dbClaims))
	for i, c := range dbClaims {
		claims[i] = sqlcClaimToDomain(c)
	}
	return claims, nil
}

func (r *claimRepository) CountFiltered(ctx context.Context, status string, dateFrom, dateTo *time.Time, search string) (int64, error) {
	count, err := r.store.CountClaimsFiltered(ctx, db.CountClaimsFilteredParams{
		StatusFilter: status,
		DateFrom:     timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:       timePtrToPgtypeTimestamptz(dateTo),
		Search:       search,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count claims filtered: %w", err)
	}
	return count, nil
}

func (r *claimRepository) CreateStatusHistory(ctx context.Context, claimID uuid.UUID, fromStatus, toStatus, action, notes string, performedBy uuid.UUID) error {
	_, err := r.store.CreateClaimStatusHistory(ctx, db.CreateClaimStatusHistoryParams{
		ClaimID:     claimID,
		FromStatus:  fromStatus,
		ToStatus:    toStatus,
		Action:      action,
		Notes:       notes,
		PerformedBy: performedBy,
	})
	if err != nil {
		return fmt.Errorf("failed to create claim status history: %w", err)
	}
	return nil
}

func (r *claimRepository) ListTimeline(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimStatusHistory, error) {
	rows, err := r.store.ListClaimTimeline(ctx, claimID)
	if err != nil {
		return nil, fmt.Errorf("failed to list claim timeline: %w", err)
	}
	entries := make([]*entity.ClaimStatusHistory, len(rows))
	for i, row := range rows {
		entries[i] = &entity.ClaimStatusHistory{
			ID:              row.ID,
			ClaimID:         row.ClaimID,
			FromStatus:      row.FromStatus,
			ToStatus:        row.ToStatus,
			Action:          row.Action,
			Notes:           row.Notes,
			PerformedBy:     row.PerformedBy,
			PerformedByName: row.PerformedByName,
			CreatedAt:       row.CreatedAt,
		}
	}
	return entries, nil
}

func sqlcClaimToDomain(c db.Claim) *entity.Claim {
	return &entity.Claim{
		ID:                   c.ID,
		ClaimNumber:          c.ClaimNumber,
		PolicyID:             c.PolicyID,
		MemberID:             c.MemberID,
		ProviderID:           c.ProviderID,
		PreAuthID:            pgtypeToUUID(c.PreauthID),
		Status:               c.Status,
		TotalAmount:          c.TotalAmount,
		ApprovedAmount:       c.ApprovedAmount,
		CoPayAmount:          c.CoPayAmount,
		MemberResponsibility: c.MemberResponsibility,
		DiagnosisCodes:       c.DiagnosisCodes,
		ServiceDate:          c.ServiceDate,
		AdmissionDate:        pgtypeTimestamptzToTimePtr(c.AdmissionDate),
		DischargeDate:        pgtypeTimestamptzToTimePtr(c.DischargeDate),
		Notes:                c.Notes,
		ClaimType:            c.ClaimType,
		VettedAmount:         pgtypeInt8ToInt64Ptr(c.VettedAmount),
		VettedBy:             pgtypeToUUID(c.VettedBy),
		VettedAt:             pgtypeTimestamptzToTimePtr(c.VettedAt),
		SLABreachAt:          pgtypeTimestamptzToTimePtr(c.SlaBreachAt),
		RejectionReason:      c.RejectionReason.String,
		EscalatedTo:          c.EscalatedTo.String,
		CreatedBy:            pgtypeToUUID(c.CreatedBy),
		CreatedAt:            c.CreatedAt,
		UpdatedAt:            c.UpdatedAt,
	}
}

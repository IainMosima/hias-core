package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type quotationVersionRepository struct {
	store db.Store
}

func NewQuotationVersionRepository(store db.Store) domainRepo.QuotationVersionRepository {
	return &quotationVersionRepository{store: store}
}

func (r *quotationVersionRepository) Create(ctx context.Context, v *entity.QuotationVersion) (*entity.QuotationVersion, error) {
	dbV, err := r.store.CreateQuotationVersion(ctx, db.CreateQuotationVersionParams{
		QuotationID:      v.QuotationID,
		VersionNumber:    int32(v.VersionNumber),
		BasePremium:      v.BasePremium,
		DiscountType:     v.DiscountType,
		DiscountValue:    v.DiscountValue,
		DiscountReason:   stringToPgtypeText(v.DiscountReason),
		LoadingType:      v.LoadingType,
		LoadingValue:     v.LoadingValue,
		LoadingReason:    stringToPgtypeText(v.LoadingReason),
		FinalPremium:     v.FinalPremium,
		MemberCount:      int32(v.MemberCount),
		ProposedMembers:  v.ProposedMembers,
		BillingFrequency: v.BillingFrequency,
		RequiresApproval: v.RequiresApproval,
		ApprovalStatus:   v.ApprovalStatus,
		PricingBreakdown: v.PricingBreakdown,
		CreatedBy:        uuidToPgtype(v.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create quotation version: %w", err)
	}
	return sqlcQuotationVersionToDomain(dbV), nil
}

func (r *quotationVersionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.QuotationVersion, error) {
	dbV, err := r.store.GetQuotationVersionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation version by ID: %w", err)
	}
	return sqlcQuotationVersionToDomain(dbV), nil
}

func (r *quotationVersionRepository) GetByQuotationAndVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int) (*entity.QuotationVersion, error) {
	dbV, err := r.store.GetQuotationVersionByNumber(ctx, db.GetQuotationVersionByNumberParams{
		QuotationID:   quotationID,
		VersionNumber: int32(versionNumber),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation version: %w", err)
	}
	return sqlcQuotationVersionToDomain(dbV), nil
}

func (r *quotationVersionRepository) ListByQuotation(ctx context.Context, quotationID uuid.UUID) ([]*entity.QuotationVersion, error) {
	dbVs, err := r.store.ListQuotationVersionsByQuotation(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list quotation versions: %w", err)
	}
	versions := make([]*entity.QuotationVersion, len(dbVs))
	for i, v := range dbVs {
		versions[i] = sqlcQuotationVersionToDomain(v)
	}
	return versions, nil
}

func (r *quotationVersionRepository) GetLatestByQuotation(ctx context.Context, quotationID uuid.UUID) (*entity.QuotationVersion, error) {
	dbV, err := r.store.GetLatestQuotationVersion(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest quotation version: %w", err)
	}
	return sqlcQuotationVersionToDomain(dbV), nil
}

func (r *quotationVersionRepository) UpdateApprovalStatus(ctx context.Context, id uuid.UUID, status string, approvedBy uuid.UUID) (*entity.QuotationVersion, error) {
	dbV, err := r.store.UpdateQuotationVersionApproval(ctx, db.UpdateQuotationVersionApprovalParams{
		ID:             id,
		ApprovalStatus: status,
		ApprovedBy:     uuidToPgtype(approvedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update approval status: %w", err)
	}
	return sqlcQuotationVersionToDomain(dbV), nil
}

func (r *quotationVersionRepository) RejectVersion(ctx context.Context, id uuid.UUID, reason string, rejectedBy uuid.UUID) (*entity.QuotationVersion, error) {
	dbV, err := r.store.RejectQuotationVersion(ctx, db.RejectQuotationVersionParams{
		ID:              id,
		RejectionReason: stringToPgtypeText(reason),
		ApprovedBy:      uuidToPgtype(rejectedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reject quotation version: %w", err)
	}
	return sqlcQuotationVersionToDomain(dbV), nil
}

func sqlcQuotationVersionToDomain(v db.QuotationVersion) *entity.QuotationVersion {
	return &entity.QuotationVersion{
		ID:               v.ID,
		QuotationID:      v.QuotationID,
		VersionNumber:    int(v.VersionNumber),
		BasePremium:      v.BasePremium,
		DiscountType:     v.DiscountType,
		DiscountValue:    v.DiscountValue,
		DiscountReason:   v.DiscountReason.String,
		LoadingType:      v.LoadingType,
		LoadingValue:     v.LoadingValue,
		LoadingReason:    v.LoadingReason.String,
		FinalPremium:     v.FinalPremium,
		MemberCount:      int(v.MemberCount),
		ProposedMembers:  v.ProposedMembers,
		BillingFrequency: v.BillingFrequency,
		RequiresApproval: v.RequiresApproval,
		ApprovalStatus:   v.ApprovalStatus,
		ApprovedBy:       pgtypeToUUID(v.ApprovedBy),
		ApprovedAt:       pgtypeTimestamptzToTimePtr(v.ApprovedAt),
		RejectionReason:  v.RejectionReason.String,
		PricingBreakdown: v.PricingBreakdown,
		CreatedBy:        pgtypeToUUID(v.CreatedBy),
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
}

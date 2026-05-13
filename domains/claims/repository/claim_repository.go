package repository

import (
	"context"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type ClaimRepository interface {
	Create(ctx context.Context, claim *entity.Claim) (*entity.Claim, error)
	GetMaxCounterForYear(ctx context.Context, year int) (int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
	GetByNumber(ctx context.Context, number string) (*entity.Claim, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Claim, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Claim, error)
	ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	Count(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status string) (int64, error)
	GetForAdjudication(ctx context.Context, limit int) ([]*entity.Claim, error)
	GetApprovedForRemittance(ctx context.Context, providerID uuid.UUID) ([]*entity.Claim, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Claim, error)
	UpdateAmounts(ctx context.Context, id uuid.UUID, approved, copay, memberResp int64) (*entity.Claim, error)
	Reject(ctx context.Context, id uuid.UUID, reason string) (*entity.Claim, error)
	GetApprovedAmountForBenefitThisYear(ctx context.Context, memberID uuid.UUID, category string) (int64, error)
	VetClaim(ctx context.Context, id uuid.UUID, vettedAmount int64, vettedBy uuid.UUID, status string) (*entity.Claim, error)
	MarkReadyForPayment(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
	ListSLABreached(ctx context.Context, limit, offset int) ([]*entity.Claim, error)
	FindByProviderAndDate(ctx context.Context, providerID uuid.UUID, serviceDate time.Time, amount int64) (*entity.Claim, error)
	ListApproachingSLA(ctx context.Context, limit, offset int) ([]*entity.Claim, error)
	CountByMemberThisMonth(ctx context.Context, memberID uuid.UUID) (int64, error)
	SetEscalatedTo(ctx context.Context, claimID uuid.UUID, role string) error
	ListFiltered(ctx context.Context, status string, dateFrom, dateTo *time.Time, search string, limit, offset int) ([]*entity.Claim, error)
	CountFiltered(ctx context.Context, status string, dateFrom, dateTo *time.Time, search string) (int64, error)
	CreateStatusHistory(ctx context.Context, claimID uuid.UUID, fromStatus, toStatus, action, notes string, performedBy uuid.UUID) error
	ListTimeline(ctx context.Context, claimID uuid.UUID) ([]*entity.ClaimStatusHistory, error)
	GetByIdempotencyKey(ctx context.Context, key string) (*entity.Claim, error)
	GetByExternalClaimID(ctx context.Context, externalID string) (*entity.Claim, error)
	CreateDraft(ctx context.Context, claim *entity.Claim) (*entity.Claim, error)
	UpdateDraft(ctx context.Context, claim *entity.Claim) (*entity.Claim, error)
	ListDrafts(ctx context.Context, createdBy uuid.UUID, limit, offset int) ([]*entity.Claim, error)
	CompleteDraft(ctx context.Context, id uuid.UUID) (*entity.Claim, error)
	DeleteDraft(ctx context.Context, id uuid.UUID) error
	UpdateClaimSource(ctx context.Context, id uuid.UUID, claimSource, idempotencyKey, externalClaimID string, sourceMetadata []byte) error
}

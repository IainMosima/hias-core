package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/preauth/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/preauth/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type preAuthRepository struct {
	store db.Store
}

func NewPreAuthRepository(store db.Store) domainRepo.PreAuthRepository {
	return &preAuthRepository{store: store}
}

func (r *preAuthRepository) Create(ctx context.Context, preauth *entity.PreAuthorization) (*entity.PreAuthorization, error) {
	dbPreAuth, err := r.store.CreatePreAuth(ctx, db.CreatePreAuthParams{
		PolicyID:       preauth.PolicyID,
		MemberID:       preauth.MemberID,
		ProviderID:     preauth.ProviderID,
		ProcedureCodes: preauth.ProcedureCodes,
		DiagnosisCodes: preauth.DiagnosisCodes,
		EstimatedCost:  preauth.EstimatedCost,
		Status:         preauth.Status,
		Notes:          preauth.Notes,
		CreatedBy:      uuidToPgtype(preauth.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pre-authorization: %w", err)
	}
	return sqlcPreAuthToDomain(dbPreAuth), nil
}

func (r *preAuthRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.PreAuthorization, error) {
	dbPreAuth, err := r.store.GetPreAuthByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pre-authorization by ID: %w", err)
	}
	return sqlcPreAuthToDomain(dbPreAuth), nil
}

func (r *preAuthRepository) GetByAuthCode(ctx context.Context, authCode string) (*entity.PreAuthorization, error) {
	dbPreAuth, err := r.store.GetPreAuthByAuthCode(ctx, pgtype.Text{String: authCode, Valid: authCode != ""})
	if err != nil {
		return nil, fmt.Errorf("failed to get pre-authorization by auth code: %w", err)
	}
	return sqlcPreAuthToDomain(dbPreAuth), nil
}

func (r *preAuthRepository) List(ctx context.Context, limit, offset int) ([]*entity.PreAuthorization, error) {
	dbPreAuths, err := r.store.ListPreAuths(ctx, db.ListPreAuthsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pre-authorizations: %w", err)
	}
	preauths := make([]*entity.PreAuthorization, len(dbPreAuths))
	for i, p := range dbPreAuths {
		preauths[i] = sqlcPreAuthToDomain(p)
	}
	return preauths, nil
}

func (r *preAuthRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.PreAuthorization, error) {
	dbPreAuths, err := r.store.ListPreAuthsByStatus(ctx, db.ListPreAuthsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pre-authorizations by status: %w", err)
	}
	preauths := make([]*entity.PreAuthorization, len(dbPreAuths))
	for i, p := range dbPreAuths {
		preauths[i] = sqlcPreAuthToDomain(p)
	}
	return preauths, nil
}

func (r *preAuthRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.PreAuthorization, error) {
	dbPreAuths, err := r.store.ListPreAuthsByProvider(ctx, db.ListPreAuthsByProviderParams{
		ProviderID: providerID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pre-authorizations by provider: %w", err)
	}
	preauths := make([]*entity.PreAuthorization, len(dbPreAuths))
	for i, p := range dbPreAuths {
		preauths[i] = sqlcPreAuthToDomain(p)
	}
	return preauths, nil
}

func (r *preAuthRepository) ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.PreAuthorization, error) {
	dbPreAuths, err := r.store.ListPreAuthsByMember(ctx, db.ListPreAuthsByMemberParams{
		MemberID: memberID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pre-authorizations by member: %w", err)
	}
	preauths := make([]*entity.PreAuthorization, len(dbPreAuths))
	for i, p := range dbPreAuths {
		preauths[i] = sqlcPreAuthToDomain(p)
	}
	return preauths, nil
}

func (r *preAuthRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountPreAuths(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count pre-authorizations: %w", err)
	}
	return count, nil
}

func (r *preAuthRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.PreAuthorization, error) {
	dbPreAuth, err := r.store.UpdatePreAuthStatus(ctx, db.UpdatePreAuthStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update pre-authorization status: %w", err)
	}
	return sqlcPreAuthToDomain(dbPreAuth), nil
}

func (r *preAuthRepository) Approve(ctx context.Context, preauth *entity.PreAuthorization) (*entity.PreAuthorization, error) {
	dbPreAuth, err := r.store.ApprovePreAuth(ctx, db.ApprovePreAuthParams{
		ID:             preauth.ID,
		AuthCode:       pgtype.Text{String: preauth.AuthCode, Valid: preauth.AuthCode != ""},
		ApprovedAmount: preauth.ApprovedAmount,
		ValidityStart:  timePtrToPgtypeTimestamptz(preauth.ValidityStart),
		ValidityEnd:    timePtrToPgtypeTimestamptz(preauth.ValidityEnd),
		ReviewedBy:     uuidToPgtype(preauth.ReviewedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to approve pre-authorization: %w", err)
	}
	return sqlcPreAuthToDomain(dbPreAuth), nil
}

func (r *preAuthRepository) Deny(ctx context.Context, id uuid.UUID, reason string, reviewedBy uuid.UUID) (*entity.PreAuthorization, error) {
	dbPreAuth, err := r.store.DenyPreAuth(ctx, db.DenyPreAuthParams{
		ID:           id,
		DenialReason: pgtype.Text{String: reason, Valid: reason != ""},
		ReviewedBy:   uuidToPgtype(reviewedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to deny pre-authorization: %w", err)
	}
	return sqlcPreAuthToDomain(dbPreAuth), nil
}

func (r *preAuthRepository) GetExpiring(ctx context.Context) ([]*entity.PreAuthorization, error) {
	dbPreAuths, err := r.store.GetExpiringPreAuths(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get expiring pre-authorizations: %w", err)
	}
	preauths := make([]*entity.PreAuthorization, len(dbPreAuths))
	for i, p := range dbPreAuths {
		preauths[i] = sqlcPreAuthToDomain(p)
	}
	return preauths, nil
}

func sqlcPreAuthToDomain(p db.Preauthorization) *entity.PreAuthorization {
	result := &entity.PreAuthorization{
		ID:             p.ID,
		PolicyID:       p.PolicyID,
		MemberID:       p.MemberID,
		ProviderID:     p.ProviderID,
		AuthCode:       p.AuthCode.String,
		ProcedureCodes: p.ProcedureCodes,
		DiagnosisCodes: p.DiagnosisCodes,
		EstimatedCost:  p.EstimatedCost,
		ApprovedAmount: p.ApprovedAmount,
		Status:         p.Status,
		Notes:          p.Notes,
		DenialReason:   p.DenialReason.String,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	if p.ReviewedBy.Valid {
		result.ReviewedBy = pgtypeToUUID(p.ReviewedBy)
	}
	if p.CreatedBy.Valid {
		result.CreatedBy = pgtypeToUUID(p.CreatedBy)
	}
	if p.ValidityStart.Valid {
		t := p.ValidityStart.Time
		result.ValidityStart = &t
	}
	if p.ValidityEnd.Valid {
		t := p.ValidityEnd.Time
		result.ValidityEnd = &t
	}
	if p.ReviewedAt.Valid {
		t := p.ReviewedAt.Time
		result.ReviewedAt = &t
	}

	return result
}

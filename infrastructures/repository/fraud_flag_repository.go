package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type fraudFlagRepository struct {
	store db.Store
}

func NewFraudFlagRepository(store db.Store) domainRepo.FraudFlagRepository {
	return &fraudFlagRepository{store: store}
}

func (r *fraudFlagRepository) Create(ctx context.Context, flag *entity.FraudFlag) (*entity.FraudFlag, error) {
	dbFlag, err := r.store.CreateFraudFlag(ctx, db.CreateFraudFlagParams{
		ClaimID:          flag.ClaimID,
		FlagType:         flag.FlagType,
		Severity:         flag.Severity,
		Details:          flag.Details,
		ReferenceClaimID: uuidToPgtype(flag.ReferenceClaimID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create fraud flag: %w", err)
	}
	return sqlcFraudFlagToDomain(dbFlag), nil
}

func (r *fraudFlagRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.FraudFlag, error) {
	dbFlag, err := r.store.GetFraudFlagByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get fraud flag by ID: %w", err)
	}
	return sqlcFraudFlagToDomain(dbFlag), nil
}

func (r *fraudFlagRepository) ListByClaim(ctx context.Context, claimID uuid.UUID) ([]*entity.FraudFlag, error) {
	dbFlags, err := r.store.ListFraudFlagsByClaim(ctx, claimID)
	if err != nil {
		return nil, fmt.Errorf("failed to list fraud flags by claim: %w", err)
	}
	flags := make([]*entity.FraudFlag, len(dbFlags))
	for i, f := range dbFlags {
		flags[i] = sqlcFraudFlagToDomain(f)
	}
	return flags, nil
}

func (r *fraudFlagRepository) ListUnresolved(ctx context.Context, limit, offset int) ([]*entity.FraudFlag, error) {
	dbFlags, err := r.store.ListUnresolvedFraudFlags(ctx, db.ListUnresolvedFraudFlagsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list unresolved fraud flags: %w", err)
	}
	flags := make([]*entity.FraudFlag, len(dbFlags))
	for i, f := range dbFlags {
		flags[i] = sqlcFraudFlagToDomain(f)
	}
	return flags, nil
}

func (r *fraudFlagRepository) Resolve(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID) (*entity.FraudFlag, error) {
	dbFlag, err := r.store.ResolveFraudFlag(ctx, db.ResolveFraudFlagParams{
		ID:         id,
		ResolvedBy: uuidToPgtype(resolvedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve fraud flag: %w", err)
	}
	return sqlcFraudFlagToDomain(dbFlag), nil
}

func (r *fraudFlagRepository) CheckDuplicate(ctx context.Context, claimNumber string, excludeID uuid.UUID) (int64, error) {
	count, err := r.store.CheckDuplicateClaim(ctx, db.CheckDuplicateClaimParams{
		ClaimNumber: claimNumber,
		ID:          excludeID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to check duplicate fraud flag: %w", err)
	}
	return count, nil
}

func (r *fraudFlagRepository) CheckFrequency(ctx context.Context, memberID, providerID uuid.UUID, procedureCode string, excludeID uuid.UUID) (int64, error) {
	count, err := r.store.CheckFrequencyClaim(ctx, db.CheckFrequencyClaimParams{
		MemberID:      memberID,
		ProviderID:    providerID,
		ProcedureCode: procedureCode,
		ID:            excludeID,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to check frequency fraud flag: %w", err)
	}
	return count, nil
}

func sqlcFraudFlagToDomain(f db.FraudFlag) *entity.FraudFlag {
	return &entity.FraudFlag{
		ID:               f.ID,
		ClaimID:          f.ClaimID,
		FlagType:         f.FlagType,
		Severity:         f.Severity,
		Details:          f.Details,
		Resolved:         f.Resolved,
		ResolvedBy:       pgtypeToUUID(f.ResolvedBy),
		ResolvedAt:       pgtypeTimestamptzToTimePtr(f.ResolvedAt),
		ReferenceClaimID: pgtypeToUUID(f.ReferenceClaimID),
		CreatedAt:        f.CreatedAt,
		UpdatedAt:        f.UpdatedAt,
	}
}

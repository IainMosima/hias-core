package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/policy/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type underwritingFlagRepository struct {
	store db.Store
}

func NewUnderwritingFlagRepository(store db.Store) domainRepo.UnderwritingFlagRepository {
	return &underwritingFlagRepository{store: store}
}

func (r *underwritingFlagRepository) Create(ctx context.Context, f *entity.UnderwritingFlag) (*entity.UnderwritingFlag, error) {
	dbF, err := r.store.CreateUnderwritingFlag(ctx, db.CreateUnderwritingFlagParams{
		AssessmentID: uuidToPgtype(f.AssessmentID),
		PolicyID:     f.PolicyID,
		MemberID:     uuidToPgtype(f.MemberID),
		FlagType:     f.FlagType,
		Severity:     f.Severity,
		Details:      f.Details,
		Status:       f.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create underwriting flag: %w", err)
	}
	return sqlcUnderwritingFlagToDomain(dbF), nil
}

func (r *underwritingFlagRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.UnderwritingFlag, error) {
	dbF, err := r.store.GetUnderwritingFlagByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get underwriting flag: %w", err)
	}
	return sqlcUnderwritingFlagToDomain(dbF), nil
}

func (r *underwritingFlagRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.UnderwritingFlag, error) {
	dbFs, err := r.store.ListUnderwritingFlagsByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list underwriting flags by policy: %w", err)
	}
	return sqlcUnderwritingFlagsToDomain(dbFs), nil
}

func (r *underwritingFlagRepository) ListByMember(ctx context.Context, memberID uuid.UUID) ([]*entity.UnderwritingFlag, error) {
	dbFs, err := r.store.ListUnderwritingFlagsByMember(ctx, uuidToPgtype(memberID))
	if err != nil {
		return nil, fmt.Errorf("failed to list underwriting flags by member: %w", err)
	}
	return sqlcUnderwritingFlagsToDomain(dbFs), nil
}

func (r *underwritingFlagRepository) ListByAssessment(ctx context.Context, assessmentID uuid.UUID) ([]*entity.UnderwritingFlag, error) {
	dbFs, err := r.store.ListUnderwritingFlagsByAssessment(ctx, uuidToPgtype(assessmentID))
	if err != nil {
		return nil, fmt.Errorf("failed to list underwriting flags by assessment: %w", err)
	}
	return sqlcUnderwritingFlagsToDomain(dbFs), nil
}

func (r *underwritingFlagRepository) ListOpen(ctx context.Context, limit, offset int32) ([]*entity.UnderwritingFlag, error) {
	dbFs, err := r.store.ListOpenUnderwritingFlags(ctx, db.ListOpenUnderwritingFlagsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list open underwriting flags: %w", err)
	}
	return sqlcUnderwritingFlagsToDomain(dbFs), nil
}

func (r *underwritingFlagRepository) CountOpen(ctx context.Context) (int64, error) {
	count, err := r.store.CountOpenUnderwritingFlags(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count open underwriting flags: %w", err)
	}
	return count, nil
}

func (r *underwritingFlagRepository) Resolve(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, resolution string) (*entity.UnderwritingFlag, error) {
	dbF, err := r.store.ResolveUnderwritingFlag(ctx, db.ResolveUnderwritingFlagParams{
		ID:         id,
		ResolvedBy: uuidToPgtype(resolvedBy),
		Resolution: stringToPgtypeText(resolution),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve underwriting flag: %w", err)
	}
	return sqlcUnderwritingFlagToDomain(dbF), nil
}

func (r *underwritingFlagRepository) Override(ctx context.Context, id uuid.UUID, resolvedBy uuid.UUID, resolution string) (*entity.UnderwritingFlag, error) {
	dbF, err := r.store.OverrideUnderwritingFlag(ctx, db.OverrideUnderwritingFlagParams{
		ID:         id,
		ResolvedBy: uuidToPgtype(resolvedBy),
		Resolution: stringToPgtypeText(resolution),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to override underwriting flag: %w", err)
	}
	return sqlcUnderwritingFlagToDomain(dbF), nil
}

func sqlcUnderwritingFlagToDomain(f db.UnderwritingFlag) *entity.UnderwritingFlag {
	return &entity.UnderwritingFlag{
		ID:           f.ID,
		AssessmentID: pgtypeToUUID(f.AssessmentID),
		PolicyID:     f.PolicyID,
		MemberID:     pgtypeToUUID(f.MemberID),
		FlagType:     f.FlagType,
		Severity:     f.Severity,
		Details:      f.Details,
		Status:       f.Status,
		ResolvedBy:   pgtypeToUUID(f.ResolvedBy),
		ResolvedAt:   pgtypeTimestamptzToTimePtr(f.ResolvedAt),
		Resolution:   f.Resolution.String,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
	}
}

func sqlcUnderwritingFlagsToDomain(flags []db.UnderwritingFlag) []*entity.UnderwritingFlag {
	result := make([]*entity.UnderwritingFlag, len(flags))
	for i, f := range flags {
		result[i] = sqlcUnderwritingFlagToDomain(f)
	}
	return result
}

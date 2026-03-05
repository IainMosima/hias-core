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

type caseRecordRepository struct {
	store db.Store
}

func NewCaseRecordRepository(store db.Store) domainRepo.CaseRecordRepository {
	return &caseRecordRepository{store: store}
}

func (r *caseRecordRepository) Create(ctx context.Context, record *entity.CaseRecord) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.CreateCaseRecord(ctx, db.CreateCaseRecordParams{
		CaseNumber:         record.CaseNumber,
		PreauthID:          record.PreAuthID,
		PolicyID:           record.PolicyID,
		MemberID:           record.MemberID,
		ProviderID:         record.ProviderID,
		Status:             record.Status,
		AdmissionDate:      timePtrToPgtypeTimestamptz(record.AdmissionDate),
		ExpectedDischarge:  timePtrToPgtypeTimestamptz(record.ExpectedDischarge),
		Diagnosis:          stringToPgtypeText(record.Diagnosis),
		TreatingDoctor:     stringToPgtypeText(record.TreatingDoctor),
		RoomType:           stringToPgtypeText(record.RoomType),
		TotalEstimatedCost: record.TotalEstimatedCost,
		Notes:              stringToPgtypeText(record.Notes),
		CreatedBy:          record.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create case record: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.GetCaseRecordByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get case record by ID: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) GetByNumber(ctx context.Context, number string) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.GetCaseRecordByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get case record by number: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) GetByPreAuth(ctx context.Context, preauthID uuid.UUID) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.GetCaseRecordByPreAuth(ctx, preauthID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case record by pre-auth ID: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.CaseRecord, error) {
	dbRecords, err := r.store.ListCaseRecordsByPolicy(ctx, db.ListCaseRecordsByPolicyParams{
		PolicyID: policyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list case records by policy: %w", err)
	}
	records := make([]*entity.CaseRecord, len(dbRecords))
	for i, rec := range dbRecords {
		records[i] = sqlcCaseRecordToDomain(rec)
	}
	return records, nil
}

func (r *caseRecordRepository) ListByMember(ctx context.Context, memberID uuid.UUID, limit, offset int) ([]*entity.CaseRecord, error) {
	dbRecords, err := r.store.ListCaseRecordsByMember(ctx, db.ListCaseRecordsByMemberParams{
		MemberID: memberID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list case records by member: %w", err)
	}
	records := make([]*entity.CaseRecord, len(dbRecords))
	for i, rec := range dbRecords {
		records[i] = sqlcCaseRecordToDomain(rec)
	}
	return records, nil
}

func (r *caseRecordRepository) ListByProvider(ctx context.Context, providerID uuid.UUID, limit, offset int) ([]*entity.CaseRecord, error) {
	dbRecords, err := r.store.ListCaseRecordsByProvider(ctx, db.ListCaseRecordsByProviderParams{
		ProviderID: providerID,
		Limit:      int32(limit),
		Offset:     int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list case records by provider: %w", err)
	}
	records := make([]*entity.CaseRecord, len(dbRecords))
	for i, rec := range dbRecords {
		records[i] = sqlcCaseRecordToDomain(rec)
	}
	return records, nil
}

func (r *caseRecordRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.CaseRecord, error) {
	dbRecords, err := r.store.ListCaseRecordsByStatus(ctx, db.ListCaseRecordsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list case records by status: %w", err)
	}
	records := make([]*entity.CaseRecord, len(dbRecords))
	for i, rec := range dbRecords {
		records[i] = sqlcCaseRecordToDomain(rec)
	}
	return records, nil
}

func (r *caseRecordRepository) Update(ctx context.Context, record *entity.CaseRecord) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.UpdateCaseRecord(ctx, db.UpdateCaseRecordParams{
		ID:                 record.ID,
		Diagnosis:          stringToPgtypeText(record.Diagnosis),
		TreatingDoctor:     stringToPgtypeText(record.TreatingDoctor),
		RoomType:           stringToPgtypeText(record.RoomType),
		TotalEstimatedCost: int64ToPgtypeInt8(record.TotalEstimatedCost),
		Notes:              stringToPgtypeText(record.Notes),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update case record: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) Admit(ctx context.Context, id uuid.UUID, admissionDate time.Time) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.AdmitCaseRecord(ctx, db.AdmitCaseRecordParams{
		ID:            id,
		AdmissionDate: timeToPgtypeTimestamptz(admissionDate),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to admit case record: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) Discharge(ctx context.Context, id uuid.UUID, dischargeDate time.Time, actualCost int64) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.DischargeCaseRecord(ctx, db.DischargeCaseRecordParams{
		ID:              id,
		ActualDischarge: timeToPgtypeTimestamptz(dischargeDate),
		TotalActualCost: actualCost,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to discharge case record: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) Close(ctx context.Context, id uuid.UUID) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.CloseCaseRecord(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to close case record: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.CaseRecord, error) {
	dbRecord, err := r.store.UpdateCaseRecordStatus(ctx, db.UpdateCaseRecordStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update case record status: %w", err)
	}
	return sqlcCaseRecordToDomain(dbRecord), nil
}

func (r *caseRecordRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	count, err := r.store.CountCaseRecordsByStatus(ctx, status)
	if err != nil {
		return 0, fmt.Errorf("failed to count case records by status: %w", err)
	}
	return count, nil
}

func sqlcCaseRecordToDomain(c db.CaseRecord) *entity.CaseRecord {
	return &entity.CaseRecord{
		ID:                 c.ID,
		CaseNumber:         c.CaseNumber,
		PreAuthID:          c.PreauthID,
		PolicyID:           c.PolicyID,
		MemberID:           c.MemberID,
		ProviderID:         c.ProviderID,
		Status:             c.Status,
		AdmissionDate:      pgtypeTimestamptzToTimePtr(c.AdmissionDate),
		ExpectedDischarge:  pgtypeTimestamptzToTimePtr(c.ExpectedDischarge),
		ActualDischarge:    pgtypeTimestamptzToTimePtr(c.ActualDischarge),
		Diagnosis:          c.Diagnosis.String,
		TreatingDoctor:     c.TreatingDoctor.String,
		RoomType:           c.RoomType.String,
		TotalEstimatedCost: c.TotalEstimatedCost,
		TotalActualCost:    c.TotalActualCost,
		Notes:              c.Notes.String,
		ClosedAt:           pgtypeTimestamptzToTimePtr(c.ClosedAt),
		CreatedBy:          c.CreatedBy,
		CreatedAt:          c.CreatedAt,
		UpdatedAt:          c.UpdatedAt,
	}
}

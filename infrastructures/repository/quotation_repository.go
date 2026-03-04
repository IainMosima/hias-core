package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/sales/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type quotationRepository struct {
	store db.Store
}

func NewQuotationRepository(store db.Store) domainRepo.QuotationRepository {
	return &quotationRepository{store: store}
}

func (r *quotationRepository) Create(ctx context.Context, q *entity.Quotation) (*entity.Quotation, error) {
	dbQ, err := r.store.CreateQuotation(ctx, db.CreateQuotationParams{
		QuotationNumber: q.QuotationNumber,
		LeadID:          q.LeadID,
		PlanID:          q.PlanID,
		QuotationType:   q.QuotationType,
		Status:          q.Status,
		CurrentVersion:  int32(q.CurrentVersion),
		ValidFrom:       timePtrToPgtypeTimestamptz(q.ValidFrom),
		ValidUntil:      timePtrToPgtypeTimestamptz(q.ValidUntil),
		ClientName:      q.ClientName,
		ClientEmail:     stringToPgtypeText(q.ClientEmail),
		ClientPhone:     stringToPgtypeText(q.ClientPhone),
		Currency:        q.Currency,
		CreatedBy:       uuidToPgtype(q.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create quotation: %w", err)
	}
	return sqlcQuotationToDomain(dbQ), nil
}

func (r *quotationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Quotation, error) {
	dbQ, err := r.store.GetQuotationByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation by ID: %w", err)
	}
	return sqlcQuotationToDomain(dbQ), nil
}

func (r *quotationRepository) GetByNumber(ctx context.Context, number string) (*entity.Quotation, error) {
	dbQ, err := r.store.GetQuotationByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotation by number: %w", err)
	}
	return sqlcQuotationToDomain(dbQ), nil
}

func (r *quotationRepository) ListByLead(ctx context.Context, leadID uuid.UUID, limit, offset int) ([]*entity.Quotation, error) {
	dbQs, err := r.store.ListQuotationsByLead(ctx, db.ListQuotationsByLeadParams{
		LeadID: leadID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list quotations by lead: %w", err)
	}
	quotations := make([]*entity.Quotation, len(dbQs))
	for i, q := range dbQs {
		quotations[i] = sqlcQuotationToDomain(q)
	}
	return quotations, nil
}

func (r *quotationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Quotation, error) {
	dbQs, err := r.store.ListQuotations(ctx, db.ListQuotationsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list quotations: %w", err)
	}
	quotations := make([]*entity.Quotation, len(dbQs))
	for i, q := range dbQs {
		quotations[i] = sqlcQuotationToDomain(q)
	}
	return quotations, nil
}

func (r *quotationRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Quotation, error) {
	dbQs, err := r.store.ListQuotationsByStatus(ctx, db.ListQuotationsByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list quotations by status: %w", err)
	}
	quotations := make([]*entity.Quotation, len(dbQs))
	for i, q := range dbQs {
		quotations[i] = sqlcQuotationToDomain(q)
	}
	return quotations, nil
}

func (r *quotationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Quotation, error) {
	dbQ, err := r.store.UpdateQuotationStatus(ctx, db.UpdateQuotationStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update quotation status: %w", err)
	}
	return sqlcQuotationToDomain(dbQ), nil
}

func (r *quotationRepository) UpdateCurrentVersion(ctx context.Context, id uuid.UUID, version int) (*entity.Quotation, error) {
	dbQ, err := r.store.UpdateQuotationCurrentVersion(ctx, db.UpdateQuotationCurrentVersionParams{
		ID:             id,
		CurrentVersion: int32(version),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update quotation current version: %w", err)
	}
	return sqlcQuotationToDomain(dbQ), nil
}

func (r *quotationRepository) SetPolicyID(ctx context.Context, id uuid.UUID, policyID uuid.UUID) (*entity.Quotation, error) {
	dbQ, err := r.store.SetQuotationPolicyID(ctx, db.SetQuotationPolicyIDParams{
		ID:       id,
		PolicyID: uuidToPgtype(policyID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to set quotation policy ID: %w", err)
	}
	return sqlcQuotationToDomain(dbQ), nil
}

func (r *quotationRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountQuotations(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count quotations: %w", err)
	}
	return count, nil
}

func (r *quotationRepository) ListExpired(ctx context.Context) ([]*entity.Quotation, error) {
	dbQs, err := r.store.ListExpiredQuotations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list expired quotations: %w", err)
	}
	quotations := make([]*entity.Quotation, len(dbQs))
	for i, q := range dbQs {
		quotations[i] = sqlcQuotationToDomain(q)
	}
	return quotations, nil
}

func sqlcQuotationToDomain(q db.Quotation) *entity.Quotation {
	return &entity.Quotation{
		ID:              q.ID,
		QuotationNumber: q.QuotationNumber,
		LeadID:          q.LeadID,
		PlanID:          q.PlanID,
		QuotationType:   q.QuotationType,
		Status:          q.Status,
		CurrentVersion:  int(q.CurrentVersion),
		PolicyID:        pgtypeToUUID(q.PolicyID),
		ValidFrom:       pgtypeTimestamptzToTimePtr(q.ValidFrom),
		ValidUntil:      pgtypeTimestamptzToTimePtr(q.ValidUntil),
		ClientName:      q.ClientName,
		ClientEmail:     q.ClientEmail.String,
		ClientPhone:     q.ClientPhone.String,
		Currency:        q.Currency,
		CreatedBy:       pgtypeToUUID(q.CreatedBy),
		CreatedAt:       q.CreatedAt,
		UpdatedAt:       q.UpdatedAt,
	}
}

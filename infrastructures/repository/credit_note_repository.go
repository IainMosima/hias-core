package repository

import (
	"context"
	"fmt"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type creditNoteRepository struct {
	store db.Store
}

func NewCreditNoteRepository(store db.Store) domainRepo.CreditNoteRepository {
	return &creditNoteRepository{store: store}
}

func (r *creditNoteRepository) Create(ctx context.Context, cn *entity.CreditNote) (*entity.CreditNote, error) {
	dbCN, err := r.store.CreateCreditNote(ctx, db.CreateCreditNoteParams{
		PolicyID:         cn.PolicyID,
		MemberID:         uuidToPgtype(cn.MemberID),
		CreditNoteNumber: cn.CreditNoteNumber,
		Amount:           cn.Amount,
		Currency:         cn.Currency,
		Reason:           cn.Reason,
		Status:           cn.Status,
		CreatedBy:        cn.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create credit note: %w", err)
	}
	return sqlcCreditNoteToDomain(dbCN), nil
}

func (r *creditNoteRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.CreditNote, error) {
	dbCN, err := r.store.GetCreditNoteByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit note: %w", err)
	}
	return sqlcCreditNoteToDomain(dbCN), nil
}

func (r *creditNoteRepository) GetByNumber(ctx context.Context, number string) (*entity.CreditNote, error) {
	dbCN, err := r.store.GetCreditNoteByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit note by number: %w", err)
	}
	return sqlcCreditNoteToDomain(dbCN), nil
}

func (r *creditNoteRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID) ([]*entity.CreditNote, error) {
	dbCNs, err := r.store.ListCreditNotesByPolicy(ctx, policyID)
	if err != nil {
		return nil, fmt.Errorf("failed to list credit notes by policy: %w", err)
	}
	return sqlcCreditNotesToDomain(dbCNs), nil
}

func (r *creditNoteRepository) ListByStatus(ctx context.Context, status string, limit, offset int32) ([]*entity.CreditNote, error) {
	dbCNs, err := r.store.ListCreditNotesByStatus(ctx, db.ListCreditNotesByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list credit notes by status: %w", err)
	}
	return sqlcCreditNotesToDomain(dbCNs), nil
}

func (r *creditNoteRepository) Approve(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) (*entity.CreditNote, error) {
	dbCN, err := r.store.ApproveCreditNote(ctx, db.ApproveCreditNoteParams{
		ID:         id,
		ApprovedBy: uuidToPgtype(approvedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to approve credit note: %w", err)
	}
	return sqlcCreditNoteToDomain(dbCN), nil
}

func (r *creditNoteRepository) Apply(ctx context.Context, id uuid.UUID, invoiceID uuid.UUID) (*entity.CreditNote, error) {
	dbCN, err := r.store.ApplyCreditNote(ctx, db.ApplyCreditNoteParams{
		ID:                 id,
		AppliedToInvoiceID: uuidToPgtype(invoiceID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to apply credit note: %w", err)
	}
	return sqlcCreditNoteToDomain(dbCN), nil
}

func (r *creditNoteRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.CreditNote, error) {
	dbCN, err := r.store.UpdateCreditNoteStatus(ctx, db.UpdateCreditNoteStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update credit note status: %w", err)
	}
	return sqlcCreditNoteToDomain(dbCN), nil
}

func sqlcCreditNoteToDomain(cn db.CreditNote) *entity.CreditNote {
	return &entity.CreditNote{
		ID:                 cn.ID,
		PolicyID:           cn.PolicyID,
		MemberID:           pgtypeToUUID(cn.MemberID),
		CreditNoteNumber:   cn.CreditNoteNumber,
		Amount:             cn.Amount,
		Currency:           cn.Currency,
		Reason:             cn.Reason,
		Status:             cn.Status,
		AppliedToInvoiceID: pgtypeToUUID(cn.AppliedToInvoiceID),
		ApprovedBy:         pgtypeToUUID(cn.ApprovedBy),
		ApprovedAt:         pgtypeTimestamptzToTimePtr(cn.ApprovedAt),
		AppliedAt:          pgtypeTimestamptzToTimePtr(cn.AppliedAt),
		CreatedBy:          cn.CreatedBy,
		CreatedAt:          cn.CreatedAt,
		UpdatedAt:          cn.UpdatedAt,
	}
}

func sqlcCreditNotesToDomain(cns []db.CreditNote) []*entity.CreditNote {
	result := make([]*entity.CreditNote, len(cns))
	for i, cn := range cns {
		result[i] = sqlcCreditNoteToDomain(cn)
	}
	return result
}

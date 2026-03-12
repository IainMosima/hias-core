package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	domainRepo "github.com/bitbiz/hias-core/domains/billing/repository"
	db "github.com/bitbiz/hias-core/infrastructures/db/sqlc"
	"github.com/google/uuid"
)

type invoiceRepository struct {
	store db.Store
}

func NewInvoiceRepository(store db.Store) domainRepo.InvoiceRepository {
	return &invoiceRepository{store: store}
}

func (r *invoiceRepository) Create(ctx context.Context, invoice *entity.Invoice) (*entity.Invoice, error) {
	dbInvoice, err := r.store.CreateInvoice(ctx, db.CreateInvoiceParams{
		PolicyID:           invoice.PolicyID,
		InvoiceNumber:      invoice.InvoiceNumber,
		Amount:             invoice.Amount,
		Currency:           invoice.Currency,
		DueDate:            invoice.DueDate,
		Status:             invoice.Status,
		BillingPeriodStart: invoice.BillingPeriodStart,
		BillingPeriodEnd:   invoice.BillingPeriodEnd,
		Notes:              invoice.Notes,
		CreatedBy:          uuidToPgtype(invoice.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}
	return sqlcInvoiceToDomain(dbInvoice), nil
}

func (r *invoiceRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	dbInvoice, err := r.store.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice by ID: %w", err)
	}
	return sqlcInvoiceToDomain(dbInvoice), nil
}

func (r *invoiceRepository) GetByNumber(ctx context.Context, number string) (*entity.Invoice, error) {
	dbInvoice, err := r.store.GetInvoiceByNumber(ctx, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice by number: %w", err)
	}
	return sqlcInvoiceToDomain(dbInvoice), nil
}

func (r *invoiceRepository) List(ctx context.Context, limit, offset int) ([]*entity.Invoice, error) {
	dbInvoices, err := r.store.ListInvoices(ctx, db.ListInvoicesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	invoices := make([]*entity.Invoice, len(dbInvoices))
	for i, inv := range dbInvoices {
		invoices[i] = sqlcInvoiceToDomain(inv)
	}
	return invoices, nil
}

func (r *invoiceRepository) ListByPolicy(ctx context.Context, policyID uuid.UUID, limit, offset int) ([]*entity.Invoice, error) {
	dbInvoices, err := r.store.ListInvoicesByPolicy(ctx, db.ListInvoicesByPolicyParams{
		PolicyID: policyID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices by policy: %w", err)
	}
	invoices := make([]*entity.Invoice, len(dbInvoices))
	for i, inv := range dbInvoices {
		invoices[i] = sqlcInvoiceToDomain(inv)
	}
	return invoices, nil
}

func (r *invoiceRepository) ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Invoice, error) {
	dbInvoices, err := r.store.ListInvoicesByStatus(ctx, db.ListInvoicesByStatusParams{
		Status: status,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices by status: %w", err)
	}
	invoices := make([]*entity.Invoice, len(dbInvoices))
	for i, inv := range dbInvoices {
		invoices[i] = sqlcInvoiceToDomain(inv)
	}
	return invoices, nil
}

func (r *invoiceRepository) ListOverdue(ctx context.Context) ([]*entity.Invoice, error) {
	dbInvoices, err := r.store.ListOverdueInvoices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list overdue invoices: %w", err)
	}
	invoices := make([]*entity.Invoice, len(dbInvoices))
	for i, inv := range dbInvoices {
		invoices[i] = sqlcInvoiceToDomain(inv)
	}
	return invoices, nil
}

func (r *invoiceRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.store.CountInvoices(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to count invoices: %w", err)
	}
	return count, nil
}

func (r *invoiceRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*entity.Invoice, error) {
	dbInvoice, err := r.store.UpdateInvoiceStatus(ctx, db.UpdateInvoiceStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice status: %w", err)
	}
	return sqlcInvoiceToDomain(dbInvoice), nil
}

func sqlcInvoiceToDomain(inv db.Invoice) *entity.Invoice {
	return &entity.Invoice{
		ID:                 inv.ID,
		PolicyID:           inv.PolicyID,
		InvoiceNumber:      inv.InvoiceNumber,
		Amount:             inv.Amount,
		Currency:           inv.Currency,
		DueDate:            inv.DueDate,
		Status:             inv.Status,
		BillingPeriodStart: inv.BillingPeriodStart,
		BillingPeriodEnd:   inv.BillingPeriodEnd,
		Notes:              inv.Notes,
		CreatedBy:          pgtypeToUUID(inv.CreatedBy),
		CreatedAt:          inv.CreatedAt,
		UpdatedAt:          inv.UpdatedAt,
	}
}

func (r *invoiceRepository) ListFiltered(ctx context.Context, dateFrom, dateTo *time.Time, limit, offset int) ([]*entity.Invoice, error) {
	dbInvoices, err := r.store.ListInvoicesFiltered(ctx, db.ListInvoicesFilteredParams{
		DateFrom:    timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:      timePtrToPgtypeTimestamptz(dateTo),
		QueryLimit:  int32(limit),
		QueryOffset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices filtered: %w", err)
	}
	invoices := make([]*entity.Invoice, len(dbInvoices))
	for i, inv := range dbInvoices {
		invoices[i] = sqlcInvoiceToDomain(inv)
	}
	return invoices, nil
}

func (r *invoiceRepository) CountFiltered(ctx context.Context, dateFrom, dateTo *time.Time) (int64, error) {
	count, err := r.store.CountInvoicesFiltered(ctx, db.CountInvoicesFilteredParams{
		DateFrom: timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:   timePtrToPgtypeTimestamptz(dateTo),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count invoices filtered: %w", err)
	}
	return count, nil
}

func (r *invoiceRepository) GetWithPolicy(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	row, err := r.store.GetInvoiceWithPolicy(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice with policy: %w", err)
	}
	inv := &entity.Invoice{
		ID: row.ID, PolicyID: row.PolicyID, InvoiceNumber: row.InvoiceNumber,
		Amount: row.Amount, Currency: row.Currency, DueDate: row.DueDate,
		Status: row.Status, BillingPeriodStart: row.BillingPeriodStart,
		BillingPeriodEnd: row.BillingPeriodEnd, Notes: row.Notes,
		CreatedBy: pgtypeToUUID(row.CreatedBy), CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		PolicyNumber: row.PolicyNumber, PolicyholderName: row.PolicyholderName,
	}
	return inv, nil
}

func (r *invoiceRepository) ListWithPolicy(ctx context.Context, limit, offset int) ([]*entity.Invoice, error) {
	rows, err := r.store.ListInvoicesWithPolicy(ctx, db.ListInvoicesWithPolicyParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices with policy: %w", err)
	}
	invoices := make([]*entity.Invoice, len(rows))
	for i, row := range rows {
		invoices[i] = &entity.Invoice{
			ID: row.ID, PolicyID: row.PolicyID, InvoiceNumber: row.InvoiceNumber,
			Amount: row.Amount, Currency: row.Currency, DueDate: row.DueDate,
			Status: row.Status, BillingPeriodStart: row.BillingPeriodStart,
			BillingPeriodEnd: row.BillingPeriodEnd, Notes: row.Notes,
			CreatedBy: pgtypeToUUID(row.CreatedBy), CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
			PolicyNumber: row.PolicyNumber, PolicyholderName: row.PolicyholderName,
		}
	}
	return invoices, nil
}

func (r *invoiceRepository) ListFilteredWithPolicy(ctx context.Context, dateFrom, dateTo *time.Time, limit, offset int) ([]*entity.Invoice, error) {
	rows, err := r.store.ListInvoicesFilteredWithPolicy(ctx, db.ListInvoicesFilteredWithPolicyParams{
		DateFrom:    timePtrToPgtypeTimestamptz(dateFrom),
		DateTo:      timePtrToPgtypeTimestamptz(dateTo),
		QueryLimit:  int32(limit),
		QueryOffset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list invoices filtered with policy: %w", err)
	}
	invoices := make([]*entity.Invoice, len(rows))
	for i, row := range rows {
		invoices[i] = &entity.Invoice{
			ID: row.ID, PolicyID: row.PolicyID, InvoiceNumber: row.InvoiceNumber,
			Amount: row.Amount, Currency: row.Currency, DueDate: row.DueDate,
			Status: row.Status, BillingPeriodStart: row.BillingPeriodStart,
			BillingPeriodEnd: row.BillingPeriodEnd, Notes: row.Notes,
			CreatedBy: pgtypeToUUID(row.CreatedBy), CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
			PolicyNumber: row.PolicyNumber, PolicyholderName: row.PolicyholderName,
		}
	}
	return invoices, nil
}

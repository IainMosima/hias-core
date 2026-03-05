package service

import (
	"context"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type CreditNoteService interface {
	CreateCreditNote(ctx context.Context, policyID, memberID uuid.UUID, amount int64, reason string, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse]
	GetCreditNote(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]billingSchema.CreditNoteResponse]
	ApproveCreditNote(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse]
	ApplyCreditNote(ctx context.Context, id uuid.UUID, invoiceID uuid.UUID) *schema.ServiceResponse[billingSchema.CreditNoteResponse]
}

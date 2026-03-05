package service

import (
	"context"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type StatementService interface {
	UploadStatement(ctx context.Context, req billingSchema.UploadStatementRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.ProviderStatementResponse]
	GetStatement(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[billingSchema.ProviderStatementResponse]
	ListByProvider(ctx context.Context, providerID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]billingSchema.ProviderStatementResponse]
	ListLineItems(ctx context.Context, statementID uuid.UUID) *schema.ServiceResponse[[]billingSchema.StatementLineItemResponse]
	ReconcileStatement(ctx context.Context, id uuid.UUID, reconciledBy uuid.UUID) *schema.ServiceResponse[billingSchema.ProviderStatementResponse]
}

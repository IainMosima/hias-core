package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	preauthSchema "github.com/bitbiz/hias-core/domains/preauth/schema"
	"github.com/google/uuid"
)

type PreAuthService interface {
	SubmitPreAuth(ctx context.Context, req preauthSchema.SubmitPreAuthRequest, createdBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse]
	GetPreAuth(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse]
	ListPreAuths(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]preauthSchema.PreAuthResponse]
	ReviewPreAuth(ctx context.Context, id uuid.UUID, req preauthSchema.ReviewPreAuthRequest, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse]
	ApprovePreAuth(ctx context.Context, id uuid.UUID, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse]
	DenyPreAuth(ctx context.Context, id uuid.UUID, reason string, reviewedBy uuid.UUID) *schema.ServiceResponse[preauthSchema.PreAuthResponse]
	ExpirePreAuths(ctx context.Context) *schema.ServiceResponse[int]
	GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64]
}

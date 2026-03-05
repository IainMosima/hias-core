package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type UnderwritingService interface {
	SubmitAssessment(ctx context.Context, req policySchema.SubmitAssessmentRequest, createdBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingResponse]
	GetAssessment(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingResponse]
	ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.UnderwritingResponse]
	ReviewAssessment(ctx context.Context, id uuid.UUID, req policySchema.ReviewAssessmentRequest, assessedBy uuid.UUID) *schema.ServiceResponse[policySchema.UnderwritingResponse]
}

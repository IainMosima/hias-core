package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type MemberService interface {
	EnrollMember(ctx context.Context, policyID uuid.UUID, req policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse]
	VerifyMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse]
	GetMemberEligibility(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[bool]
	ListMembers(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.MemberResponse]
}

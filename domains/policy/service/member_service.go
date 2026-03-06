package service

import (
	"context"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/google/uuid"
)

type MemberService interface {
	EnrollMember(ctx context.Context, policyID uuid.UUID, req policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse]
	GetMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse]
	UpdateMember(ctx context.Context, memberID uuid.UUID, req policySchema.UpdateMemberRequest) *schema.ServiceResponse[policySchema.MemberResponse]
	RemoveMember(ctx context.Context, memberID uuid.UUID, reason string) *schema.ServiceResponse[string]
	SuspendMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse]
	ReactivateMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse]
	VerifyMember(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[policySchema.MemberResponse]
	GetMemberEligibility(ctx context.Context, memberID uuid.UUID) *schema.ServiceResponse[bool]
	ListMembers(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.MemberResponse]
	BulkEnrollMembers(ctx context.Context, policyID uuid.UUID, reqs []policySchema.EnrollMemberRequest) *schema.ServiceResponse[policySchema.BulkMemberResultResponse]
	BulkRemoveMembers(ctx context.Context, policyID uuid.UUID, memberIDs []uuid.UUID, reason string) *schema.ServiceResponse[policySchema.BulkMemberResultResponse]
	ImportMembersCSV(ctx context.Context, policyID uuid.UUID, csvData []byte) *schema.ServiceResponse[policySchema.BulkMemberResultResponse]
	ListMembersFiltered(ctx context.Context, search string, page, pageSize int) *schema.ServiceResponse[[]policySchema.MemberResponse]
	CountMembersFiltered(ctx context.Context, search string) *schema.ServiceResponse[int64]
}

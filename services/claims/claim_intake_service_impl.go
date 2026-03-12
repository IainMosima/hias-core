package claims

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	claimSvc "github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	memberRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type claimIntakeService struct {
	claimService claimSvc.ClaimService
	claimRepo    claimRepo.ClaimRepository
	memberRepo   memberRepo.MemberRepository
	partnerRepo  claimRepo.APIPartnerRepository
}

func NewClaimIntakeService(
	claimService claimSvc.ClaimService,
	claimRepository claimRepo.ClaimRepository,
	memberRepository memberRepo.MemberRepository,
	partnerRepository claimRepo.APIPartnerRepository,
) claimSvc.ClaimIntakeService {
	return &claimIntakeService{
		claimService: claimService,
		claimRepo:    claimRepository,
		memberRepo:   memberRepository,
		partnerRepo:  partnerRepository,
	}
}

func (s *claimIntakeService) SubmitExternal(ctx context.Context, req claimsSchema.ExternalClaimRequest, partner *entity.APIPartner) *schema.ServiceResponse[claimsSchema.ExternalClaimResponse] {
	// 1. Idempotency check
	existingClaim, err := s.claimRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil && existingClaim != nil {
		return schema.NewServiceResponse(claimsSchema.ExternalClaimResponse{
			ClaimID:     existingClaim.ID,
			ClaimNumber: existingClaim.ClaimNumber,
			Status:      InternalToExternalStatus(existingClaim.Status),
			ReceivedAt:  existingClaim.CreatedAt,
		}, http.StatusOK, "Claim already submitted (idempotent)")
	}

	// 2. Resolve member by member_number
	member, err := s.memberRepo.GetByNumber(ctx, req.MemberNumber)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ExternalClaimResponse](
			http.StatusBadRequest, fmt.Sprintf("Member not found: %s", req.MemberNumber), err)
	}

	// 3. Resolve provider from partner's linked provider
	providerID := partner.ProviderID
	if providerID == uuid.Nil {
		return schema.NewServiceErrorResponse[claimsSchema.ExternalClaimResponse](
			http.StatusBadRequest, "Partner has no linked provider and no provider_code specified", fmt.Errorf("no provider linked"))
	}

	// 4. Resolve policy from member's policy
	policyID := member.PolicyID

	// 5. Build SubmitClaimRequest and call the existing pipeline
	claimType := req.ClaimType
	if claimType == "" {
		claimType = string(shared.ClaimTypeDirect)
	}

	submitReq := claimsSchema.SubmitClaimRequest{
		PolicyID:       policyID.String(),
		MemberID:       member.ID.String(),
		ProviderID:     providerID.String(),
		DiagnosisCodes: req.DiagnosisCodes,
		ServiceDate:    req.ServiceDate,
		AdmissionDate:  req.AdmissionDate,
		DischargeDate:  req.DischargeDate,
		Notes:          req.Notes,
		ClaimType:      claimType,
		LineItems:      req.LineItems,
	}

	// Use a system user ID for external submissions
	systemUserID := uuid.Nil
	resp := s.claimService.SubmitClaim(ctx, submitReq, systemUserID)
	if resp.Error != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ExternalClaimResponse](
			resp.StatusCode, resp.Message, resp.Error)
	}

	// 6. Update claim with source tracking
	sourceMetadata, _ := json.Marshal(req.Metadata)
	if req.Metadata == nil {
		sourceMetadata = []byte("{}")
	}

	claimSource := string(shared.ClaimSourceProviderPortal)
	if partner.PartnerType == string(shared.PartnerTypePartnerNetwork) || partner.PartnerType == string(shared.PartnerTypeTPA) {
		claimSource = string(shared.ClaimSourcePartnerAPI)
	}

	_ = s.claimRepo.UpdateClaimSource(ctx, resp.Data.ID, claimSource, req.IdempotencyKey, req.ExternalClaimID, sourceMetadata)

	return schema.NewServiceResponse(claimsSchema.ExternalClaimResponse{
		ClaimID:     resp.Data.ID,
		ClaimNumber: resp.Data.ClaimNumber,
		Status:      InternalToExternalStatus(resp.Data.Status),
		ReceivedAt:  resp.Data.CreatedAt,
	}, http.StatusCreated, "Claim submitted successfully")
}

func (s *claimIntakeService) GetExternalStatus(ctx context.Context, claimID uuid.UUID, partnerID uuid.UUID) *schema.ServiceResponse[claimsSchema.ExternalClaimStatusResponse] {
	claim, err := s.claimRepo.GetByID(ctx, claimID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ExternalClaimStatusResponse](
			http.StatusNotFound, "Claim not found", err)
	}

	// Verify the claim belongs to the partner's provider
	partner, err := s.partnerRepo.GetByID(ctx, partnerID)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ExternalClaimStatusResponse](
			http.StatusForbidden, "Partner not found", err)
	}

	if partner.ProviderID != uuid.Nil && claim.ProviderID != partner.ProviderID {
		return schema.NewServiceErrorResponse[claimsSchema.ExternalClaimStatusResponse](
			http.StatusForbidden, "Claim does not belong to your provider", fmt.Errorf("provider mismatch"))
	}

	return schema.NewServiceResponse(claimsSchema.ExternalClaimStatusResponse{
		ClaimID:        claim.ID,
		ClaimNumber:    claim.ClaimNumber,
		Status:         InternalToExternalStatus(claim.Status),
		TotalAmount:    claim.TotalAmount,
		ApprovedAmount: claim.ApprovedAmount,
		UpdatedAt:      claim.UpdatedAt,
	}, http.StatusOK, "Claim status retrieved")
}

func (s *claimIntakeService) CreateDraft(ctx context.Context, req claimsSchema.DraftClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	// Parse optional UUIDs
	var policyID, memberID, providerID, preAuthID uuid.UUID
	if req.PolicyID != "" {
		policyID, _ = uuid.Parse(req.PolicyID)
	}
	if req.MemberID != "" {
		memberID, _ = uuid.Parse(req.MemberID)
	}
	if req.ProviderID != "" {
		providerID, _ = uuid.Parse(req.ProviderID)
	}
	if req.PreAuthID != "" {
		preAuthID, _ = uuid.Parse(req.PreAuthID)
	}

	// Calculate total from line items
	var totalAmount int64
	for _, item := range req.LineItems {
		totalAmount += item.UnitPrice * int64(item.Quantity)
	}

	// Encode diagnosis codes
	diagCodes, _ := json.Marshal(req.DiagnosisCodes)
	if req.DiagnosisCodes == nil {
		diagCodes = []byte("[]")
	}

	serviceDate := time.Now()
	if req.ServiceDate != nil {
		serviceDate = *req.ServiceDate
	}

	claimType := req.ClaimType
	if claimType == "" {
		claimType = string(shared.ClaimTypeDirect)
	}

	// Generate a draft claim number
	claimNumber := fmt.Sprintf("DRF-%s-%06d", time.Now().Format("2006"), time.Now().UnixNano()%1000000)

	draft := &entity.Claim{
		ClaimNumber:    claimNumber,
		PolicyID:       policyID,
		MemberID:       memberID,
		ProviderID:     providerID,
		PreAuthID:      preAuthID,
		Status:         string(shared.ClaimStatusReceived),
		TotalAmount:    totalAmount,
		DiagnosisCodes: diagCodes,
		ServiceDate:    serviceDate,
		Notes:          req.Notes,
		ClaimType:      claimType,
		CreatedBy:      createdBy,
	}

	created, err := s.claimRepo.CreateDraft(ctx, draft)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusInternalServerError, "Failed to create draft claim", err)
	}

	resp := claimsSchema.ToClaimResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Draft claim created")
}

func (s *claimIntakeService) UpdateDraft(ctx context.Context, id uuid.UUID, req claimsSchema.DraftClaimRequest, updatedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	// Verify the draft exists
	existing, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusNotFound, "Draft claim not found", err)
	}
	if !existing.IsDraft {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest, "Claim is not a draft", fmt.Errorf("claim is not a draft"))
	}

	// Parse optional UUIDs
	var policyID, memberID, providerID, preAuthID uuid.UUID
	if req.PolicyID != "" {
		policyID, _ = uuid.Parse(req.PolicyID)
	}
	if req.MemberID != "" {
		memberID, _ = uuid.Parse(req.MemberID)
	}
	if req.ProviderID != "" {
		providerID, _ = uuid.Parse(req.ProviderID)
	}
	if req.PreAuthID != "" {
		preAuthID, _ = uuid.Parse(req.PreAuthID)
	}

	var totalAmount int64
	for _, item := range req.LineItems {
		totalAmount += item.UnitPrice * int64(item.Quantity)
	}

	diagCodes, _ := json.Marshal(req.DiagnosisCodes)
	if req.DiagnosisCodes == nil {
		diagCodes = []byte("[]")
	}

	serviceDate := existing.ServiceDate
	if req.ServiceDate != nil {
		serviceDate = *req.ServiceDate
	}

	claimType := req.ClaimType
	if claimType == "" {
		claimType = existing.ClaimType
	}

	draft := &entity.Claim{
		ID:             id,
		PolicyID:       policyID,
		MemberID:       memberID,
		ProviderID:     providerID,
		PreAuthID:      preAuthID,
		DiagnosisCodes: diagCodes,
		ServiceDate:    serviceDate,
		Notes:          req.Notes,
		ClaimType:      claimType,
		TotalAmount:    totalAmount,
	}

	updated, err := s.claimRepo.UpdateDraft(ctx, draft)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusInternalServerError, "Failed to update draft claim", err)
	}

	resp := claimsSchema.ToClaimResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Draft claim updated")
}

func (s *claimIntakeService) ListDrafts(ctx context.Context, createdBy uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]claimsSchema.ClaimResponse] {
	offset := (page - 1) * pageSize
	drafts, err := s.claimRepo.ListDrafts(ctx, createdBy, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]claimsSchema.ClaimResponse](
			http.StatusInternalServerError, "Failed to list draft claims", err)
	}

	responses := make([]claimsSchema.ClaimResponse, len(drafts))
	for i, d := range drafts {
		responses[i] = claimsSchema.ToClaimResponse(d)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Draft claims retrieved")
}

func (s *claimIntakeService) SubmitDraft(ctx context.Context, id uuid.UUID, submittedBy uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	// Get the draft
	draft, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusNotFound, "Draft claim not found", err)
	}
	if !draft.IsDraft {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest, "Claim is not a draft", fmt.Errorf("claim is not a draft"))
	}

	// Validate required fields for submission
	if draft.PolicyID == uuid.Nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest, "Policy ID is required to submit", fmt.Errorf("missing policy_id"))
	}
	if draft.MemberID == uuid.Nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest, "Member ID is required to submit", fmt.Errorf("missing member_id"))
	}
	if draft.ProviderID == uuid.Nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest, "Provider ID is required to submit", fmt.Errorf("missing provider_id"))
	}

	// Parse line items from diagnosis codes (we need to build a SubmitClaimRequest)
	var diagCodes []string
	_ = json.Unmarshal(draft.DiagnosisCodes, &diagCodes)

	submitReq := claimsSchema.SubmitClaimRequest{
		PolicyID:       draft.PolicyID.String(),
		MemberID:       draft.MemberID.String(),
		ProviderID:     draft.ProviderID.String(),
		PreAuthID:      draft.PreAuthID.String(),
		DiagnosisCodes: diagCodes,
		ServiceDate:    draft.ServiceDate,
		Notes:          draft.Notes,
		ClaimType:      draft.ClaimType,
		LineItems:      []claimsSchema.LineItemRequest{}, // Line items would need separate storage for drafts
	}

	// If the draft has no separate line items, create a single line item from total
	if draft.TotalAmount > 0 && len(submitReq.LineItems) == 0 {
		submitReq.LineItems = []claimsSchema.LineItemRequest{
			{
				ProcedureCode: "DRAFT",
				ProcedureName: "Draft claim submission",
				Quantity:      1,
				UnitPrice:     draft.TotalAmount,
			},
		}
	}

	// Delete the draft first
	_ = s.claimRepo.DeleteDraft(ctx, id)

	// Submit through the standard pipeline
	resp := s.claimService.SubmitClaim(ctx, submitReq, submittedBy)
	if resp.Error != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			resp.StatusCode, resp.Message, resp.Error)
	}

	return schema.NewServiceResponse(resp.Data, http.StatusOK, "Draft claim submitted successfully")
}

func (s *claimIntakeService) DeleteDraft(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[claimsSchema.ClaimResponse] {
	existing, err := s.claimRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusNotFound, "Draft claim not found", err)
	}
	if !existing.IsDraft {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusBadRequest, "Claim is not a draft", fmt.Errorf("claim is not a draft"))
	}

	if err := s.claimRepo.DeleteDraft(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[claimsSchema.ClaimResponse](
			http.StatusInternalServerError, "Failed to delete draft claim", err)
	}

	resp := claimsSchema.ToClaimResponse(existing)
	return schema.NewServiceResponse(resp, http.StatusOK, "Draft claim deleted")
}

// InternalToExternalStatus maps internal claim statuses to simplified external-facing statuses.
func InternalToExternalStatus(internal string) string {
	switch internal {
	case string(shared.ClaimStatusReceived), string(shared.ClaimStatusValidated):
		return string(shared.ExternalStatusReceived)
	case string(shared.ClaimStatusAdjudicated), string(shared.ClaimStatusVetted), string(shared.ClaimStatusPartiallyVetted):
		return string(shared.ExternalStatusProcessing)
	case string(shared.ClaimStatusManualReview), string(shared.ClaimStatusEscalated):
		return string(shared.ExternalStatusUnderReview)
	case string(shared.ClaimStatusApproved), string(shared.ClaimStatusReadyForPayment):
		return string(shared.ExternalStatusApproved)
	case string(shared.ClaimStatusRejected):
		return string(shared.ExternalStatusRejected)
	case string(shared.ClaimStatusPaid), string(shared.ClaimStatusPartPaid):
		return string(shared.ExternalStatusSettled)
	default:
		return string(shared.ExternalStatusProcessing)
	}
}

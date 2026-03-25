package sales

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	billingService "github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	notifService "github.com/bitbiz/hias-core/domains/notification/service"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	policyService "github.com/bitbiz/hias-core/domains/policy/service"
	productService "github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/domains/sales/entity"
	salesRepo "github.com/bitbiz/hias-core/domains/sales/repository"
	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/bitbiz/hias-core/domains/sales/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type quotationServiceImpl struct {
	quotationRepo  salesRepo.QuotationRepository
	versionRepo    salesRepo.QuotationVersionRepository
	documentRepo   salesRepo.QuotationDocumentRepository
	approvalRepo   salesRepo.ApprovalLimitRepository
	leadRepo       salesRepo.LeadRepository
	auditSvc       auditService.AuditService
	premiumRuleSvc productService.PremiumRuleService
	notifSvc       notifService.NotificationService
	policySvc      policyService.PolicyService
	memberSvc      policyService.MemberService
	installmentSvc billingService.InstallmentService
}

func NewQuotationService(
	quotationRepo salesRepo.QuotationRepository,
	versionRepo salesRepo.QuotationVersionRepository,
	documentRepo salesRepo.QuotationDocumentRepository,
	approvalRepo salesRepo.ApprovalLimitRepository,
	leadRepo salesRepo.LeadRepository,
	auditSvc auditService.AuditService,
	premiumRuleSvc productService.PremiumRuleService,
	notifSvc notifService.NotificationService,
	policySvc policyService.PolicyService,
	memberSvc policyService.MemberService,
	installmentSvc billingService.InstallmentService,
) service.QuotationService {
	return &quotationServiceImpl{
		quotationRepo:  quotationRepo,
		versionRepo:    versionRepo,
		documentRepo:   documentRepo,
		approvalRepo:   approvalRepo,
		leadRepo:       leadRepo,
		auditSvc:       auditSvc,
		premiumRuleSvc: premiumRuleSvc,
		notifSvc:       notifSvc,
		policySvc:      policySvc,
		memberSvc:      memberSvc,
		installmentSvc: installmentSvc,
	}
}

func (s *quotationServiceImpl) CreateQuotation(ctx context.Context, req salesSchema.CreateQuotationRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationDetailResponse] {
	leadID, _ := uuid.Parse(req.LeadID)
	planID, _ := uuid.Parse(req.PlanID)

	// Validate lead exists
	lead, err := s.leadRepo.GetByID(ctx, leadID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDetailResponse](http.StatusNotFound, "Lead not found", err)
	}

	// Validate lead status
	validLeadStatuses := map[string]bool{
		string(shared.LeadStatusNew):          true,
		string(shared.LeadStatusContacted):    true,
		string(shared.LeadStatusQualified):    true,
		string(shared.LeadStatusProposalSent): true,
		string(shared.LeadStatusNegotiation):  true,
	}
	if !validLeadStatuses[lead.Status] {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDetailResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot create quotation for lead in %s status", lead.Status),
			nil,
		)
	}

	now := time.Now()
	validUntil := now.AddDate(0, 0, shared.QuotationValidityDays)

	// Set default discount/loading types
	discountType := req.DiscountType
	if discountType == "" {
		discountType = string(shared.DiscountTypePercentage)
	}
	loadingType := req.LoadingType
	if loadingType == "" {
		loadingType = string(shared.LoadingTypePercentage)
	}

	proposedMembers := req.ProposedMembers
	if proposedMembers == nil {
		proposedMembers = json.RawMessage("[]")
	}

	// Calculate base premium via underwriting rules (with age-band support)
	premiumResp := s.premiumRuleSvc.CalculatePremiumWithMembers(ctx, planID, req.MemberCount, proposedMembers)
	if premiumResp.Error != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDetailResponse](
			http.StatusInternalServerError, "Failed to calculate premium from underwriting rules", premiumResp.Error)
	}
	basePremium := premiumResp.Data

	// Calculate final premium
	finalPremium := calculateFinalPremium(basePremium, discountType, req.DiscountValue, loadingType, req.LoadingValue)

	// Check if approval is required based on actual limits
	requiresApproval := false
	if req.DiscountValue > 0 || req.LoadingValue > 0 {
		defaultLimit, _ := s.approvalRepo.GetByRole(ctx, string(shared.UserRoleSalesAgent))
		if defaultLimit == nil {
			requiresApproval = true
		} else {
			if discountType == string(shared.DiscountTypePercentage) && req.DiscountValue > defaultLimit.MaxDiscountPercentage {
				requiresApproval = true
			} else if discountType == string(shared.DiscountTypeFixed) && req.DiscountValue > defaultLimit.MaxDiscountAmount {
				requiresApproval = true
			}
			if loadingType == string(shared.LoadingTypePercentage) && req.LoadingValue > defaultLimit.MaxLoadingPercentage {
				requiresApproval = true
			} else if loadingType == string(shared.LoadingTypeFixed) && req.LoadingValue > defaultLimit.MaxLoadingAmount {
				requiresApproval = true
			}
		}
	}
	approvalStatus := string(shared.ApprovalStatusNone)
	if requiresApproval {
		approvalStatus = string(shared.ApprovalStatusPending)
	}

	// Build pricing breakdown
	breakdown, _ := json.Marshal(map[string]interface{}{
		"base_premium":   basePremium,
		"discount_type":  discountType,
		"discount_value": req.DiscountValue,
		"loading_type":   loadingType,
		"loading_value":  req.LoadingValue,
		"final_premium":  finalPremium,
		"member_count":   req.MemberCount,
	})

	// Create quotation
	quotation := &entity.Quotation{
		QuotationNumber: utils.GenerateQuotationNumber(),
		LeadID:          leadID,
		PlanID:          planID,
		QuotationType:   req.QuotationType,
		Status:          string(shared.QuotationStatusDraft),
		CurrentVersion:  1,
		ValidFrom:       &now,
		ValidUntil:      &validUntil,
		ClientName:      req.ClientName,
		ClientEmail:     req.ClientEmail,
		ClientPhone:     req.ClientPhone,
		Currency:        string(shared.CurrencyKES),
		CreatedBy:       createdBy,
	}

	created, err := s.quotationRepo.Create(ctx, quotation)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDetailResponse](http.StatusInternalServerError, "Failed to create quotation", err)
	}

	// Create version 1
	version := &entity.QuotationVersion{
		QuotationID:      created.ID,
		VersionNumber:    1,
		BasePremium:      basePremium,
		DiscountType:     discountType,
		DiscountValue:    req.DiscountValue,
		DiscountReason:   req.DiscountReason,
		LoadingType:      loadingType,
		LoadingValue:     req.LoadingValue,
		LoadingReason:    req.LoadingReason,
		FinalPremium:     finalPremium,
		MemberCount:      req.MemberCount,
		ProposedMembers:  proposedMembers,
		BillingFrequency: req.BillingFrequency,
		RequiresApproval: requiresApproval,
		ApprovalStatus:   approvalStatus,
		PricingBreakdown: breakdown,
		CreatedBy:        createdBy,
	}

	createdVersion, err := s.versionRepo.Create(ctx, version)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDetailResponse](http.StatusInternalServerError, "Failed to create quotation version", err)
	}

	// Update lead status to PROPOSAL_SENT if it's in an earlier stage
	if lead.Status == string(shared.LeadStatusNew) || lead.Status == string(shared.LeadStatusContacted) || lead.Status == string(shared.LeadStatusQualified) {
		if _, updateErr := s.leadRepo.UpdateStatus(ctx, leadID, string(shared.LeadStatusProposalSent)); updateErr != nil {
			log.Printf("Failed to update lead status: %v", updateErr)
		}
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeQuotation), created.ID, string(shared.AuditActionCreate))

	response := salesSchema.QuotationDetailResponse{
		QuotationResponse: salesSchema.ToQuotationResponse(created),
		Versions:          []salesSchema.QuotationVersionResponse{salesSchema.ToQuotationVersionResponse(createdVersion)},
		Documents:         []salesSchema.QuotationDocumentResponse{},
	}

	return schema.NewServiceResponse(response, http.StatusCreated, "Quotation created")
}

func (s *quotationServiceImpl) GetQuotation(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationDetailResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDetailResponse](http.StatusNotFound, "Quotation not found", err)
	}

	versions, _ := s.versionRepo.ListByQuotation(ctx, id)
	documents, _ := s.documentRepo.ListByQuotation(ctx, id)

	versionResponses := make([]salesSchema.QuotationVersionResponse, len(versions))
	for i, v := range versions {
		versionResponses[i] = salesSchema.ToQuotationVersionResponse(v)
	}

	docResponses := make([]salesSchema.QuotationDocumentResponse, len(documents))
	for i, d := range documents {
		docResponses[i] = salesSchema.ToQuotationDocumentResponse(d)
	}

	response := salesSchema.QuotationDetailResponse{
		QuotationResponse: salesSchema.ToQuotationResponse(quotation),
		Versions:          versionResponses,
		Documents:         docResponses,
	}

	return schema.NewServiceResponse(response, http.StatusOK, "Quotation retrieved")
}

func (s *quotationServiceImpl) ListQuotations(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]salesSchema.QuotationResponse] {
	offset := (page - 1) * pageSize
	quotations, err := s.quotationRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to list quotations", err)
	}
	responses := make([]salesSchema.QuotationResponse, len(quotations))
	for i, q := range quotations {
		responses[i] = salesSchema.ToQuotationResponse(q)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Quotations retrieved")
}

func (s *quotationServiceImpl) ListQuotationsByLead(ctx context.Context, leadID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]salesSchema.QuotationResponse] {
	offset := (page - 1) * pageSize
	quotations, err := s.quotationRepo.ListByLead(ctx, leadID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to list quotations by lead", err)
	}
	responses := make([]salesSchema.QuotationResponse, len(quotations))
	for i, q := range quotations {
		responses[i] = salesSchema.ToQuotationResponse(q)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Quotations retrieved")
}

func (s *quotationServiceImpl) CreateVersion(ctx context.Context, quotationID uuid.UUID, req salesSchema.CreateQuotationVersionRequest, createdBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationVersionResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, quotationID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusNotFound, "Quotation not found", err)
	}

	// Can only create new versions for DRAFT or ISSUED
	if quotation.Status != string(shared.QuotationStatusDraft) && quotation.Status != string(shared.QuotationStatusIssued) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot create version for quotation in %s status", quotation.Status),
			nil,
		)
	}

	// Get latest version number
	latest, err := s.versionRepo.GetLatestByQuotation(ctx, quotationID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusInternalServerError, "Failed to get latest version", err)
	}

	newVersionNumber := latest.VersionNumber + 1

	discountType := req.DiscountType
	if discountType == "" {
		discountType = string(shared.DiscountTypePercentage)
	}
	loadingType := req.LoadingType
	if loadingType == "" {
		loadingType = string(shared.LoadingTypePercentage)
	}

	proposedMembers := req.ProposedMembers
	if proposedMembers == nil {
		proposedMembers = json.RawMessage("[]")
	}

	// Calculate base premium via underwriting rules (with age-band support)
	premiumResp := s.premiumRuleSvc.CalculatePremiumWithMembers(ctx, quotation.PlanID, req.MemberCount, proposedMembers)
	if premiumResp.Error != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](
			http.StatusInternalServerError, "Failed to calculate premium from underwriting rules", premiumResp.Error)
	}
	basePremium := premiumResp.Data
	finalPremium := calculateFinalPremium(basePremium, discountType, req.DiscountValue, loadingType, req.LoadingValue)

	requiresApproval := false
	if req.DiscountValue > 0 || req.LoadingValue > 0 {
		defaultLimit, _ := s.approvalRepo.GetByRole(ctx, string(shared.UserRoleSalesAgent))
		if defaultLimit == nil {
			requiresApproval = true
		} else {
			if discountType == string(shared.DiscountTypePercentage) && req.DiscountValue > defaultLimit.MaxDiscountPercentage {
				requiresApproval = true
			} else if discountType == string(shared.DiscountTypeFixed) && req.DiscountValue > defaultLimit.MaxDiscountAmount {
				requiresApproval = true
			}
			if loadingType == string(shared.LoadingTypePercentage) && req.LoadingValue > defaultLimit.MaxLoadingPercentage {
				requiresApproval = true
			} else if loadingType == string(shared.LoadingTypeFixed) && req.LoadingValue > defaultLimit.MaxLoadingAmount {
				requiresApproval = true
			}
		}
	}
	approvalStatus := string(shared.ApprovalStatusNone)
	if requiresApproval {
		approvalStatus = string(shared.ApprovalStatusPending)
	}

	breakdown, _ := json.Marshal(map[string]interface{}{
		"base_premium":   basePremium,
		"discount_type":  discountType,
		"discount_value": req.DiscountValue,
		"loading_type":   loadingType,
		"loading_value":  req.LoadingValue,
		"final_premium":  finalPremium,
		"member_count":   req.MemberCount,
	})

	version := &entity.QuotationVersion{
		QuotationID:      quotationID,
		VersionNumber:    newVersionNumber,
		BasePremium:      basePremium,
		DiscountType:     discountType,
		DiscountValue:    req.DiscountValue,
		DiscountReason:   req.DiscountReason,
		LoadingType:      loadingType,
		LoadingValue:     req.LoadingValue,
		LoadingReason:    req.LoadingReason,
		FinalPremium:     finalPremium,
		MemberCount:      req.MemberCount,
		ProposedMembers:  proposedMembers,
		BillingFrequency: req.BillingFrequency,
		RequiresApproval: requiresApproval,
		ApprovalStatus:   approvalStatus,
		PricingBreakdown: breakdown,
		CreatedBy:        createdBy,
	}

	created, err := s.versionRepo.Create(ctx, version)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusInternalServerError, "Failed to create version", err)
	}

	// Update quotation's current_version
	if _, updateErr := s.quotationRepo.UpdateCurrentVersion(ctx, quotationID, newVersionNumber); updateErr != nil {
		log.Printf("Failed to update quotation current version: %v", updateErr)
	}

	s.logAudit(ctx, createdBy, string(shared.AuditEntityTypeQuotationVersion), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(salesSchema.ToQuotationVersionResponse(created), http.StatusCreated, "Version created")
}

func (s *quotationServiceImpl) GetVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int) *schema.ServiceResponse[salesSchema.QuotationVersionResponse] {
	version, err := s.versionRepo.GetByQuotationAndVersion(ctx, quotationID, versionNumber)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusNotFound, "Version not found", err)
	}
	return schema.NewServiceResponse(salesSchema.ToQuotationVersionResponse(version), http.StatusOK, "Version retrieved")
}

func (s *quotationServiceImpl) ListVersions(ctx context.Context, quotationID uuid.UUID) *schema.ServiceResponse[[]salesSchema.QuotationVersionResponse] {
	versions, err := s.versionRepo.ListByQuotation(ctx, quotationID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.QuotationVersionResponse](http.StatusInternalServerError, "Failed to list versions", err)
	}
	responses := make([]salesSchema.QuotationVersionResponse, len(versions))
	for i, v := range versions {
		responses[i] = salesSchema.ToQuotationVersionResponse(v)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Versions retrieved")
}

func (s *quotationServiceImpl) CompareVersions(ctx context.Context, quotationID uuid.UUID, versionA, versionB int) *schema.ServiceResponse[salesSchema.VersionComparisonResponse] {
	va, err := s.versionRepo.GetByQuotationAndVersion(ctx, quotationID, versionA)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.VersionComparisonResponse](http.StatusNotFound, fmt.Sprintf("Version %d not found", versionA), err)
	}
	vb, err := s.versionRepo.GetByQuotationAndVersion(ctx, quotationID, versionB)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.VersionComparisonResponse](http.StatusNotFound, fmt.Sprintf("Version %d not found", versionB), err)
	}

	response := salesSchema.VersionComparisonResponse{
		VersionA: salesSchema.ToQuotationVersionResponse(va),
		VersionB: salesSchema.ToQuotationVersionResponse(vb),
		PricingDiff: salesSchema.PricingDiff{
			BasePremiumDiff:  vb.BasePremium - va.BasePremium,
			DiscountDiff:     vb.DiscountValue - va.DiscountValue,
			LoadingDiff:      vb.LoadingValue - va.LoadingValue,
			FinalPremiumDiff: vb.FinalPremium - va.FinalPremium,
			MemberCountDiff:  vb.MemberCount - va.MemberCount,
		},
	}

	return schema.NewServiceResponse(response, http.StatusOK, "Version comparison")
}

func (s *quotationServiceImpl) SubmitForApproval(ctx context.Context, quotationID uuid.UUID, versionNumber int, submittedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationVersionResponse] {
	version, err := s.versionRepo.GetByQuotationAndVersion(ctx, quotationID, versionNumber)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusNotFound, "Version not found", err)
	}

	if version.ApprovalStatus != string(shared.ApprovalStatusNone) && version.ApprovalStatus != string(shared.ApprovalStatusRejected) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Version already in %s status", version.ApprovalStatus),
			nil,
		)
	}

	updated, err := s.versionRepo.UpdateApprovalStatus(ctx, version.ID, string(shared.ApprovalStatusPending), uuid.Nil)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusInternalServerError, "Failed to submit for approval", err)
	}

	s.logAudit(ctx, submittedBy, string(shared.AuditEntityTypeQuotationVersion), version.ID, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationVersionResponse(updated), http.StatusOK, "Submitted for approval")
}

func (s *quotationServiceImpl) ApproveVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int, req salesSchema.ApproveVersionRequest, approvedBy uuid.UUID, approverRole string) *schema.ServiceResponse[salesSchema.QuotationVersionResponse] {
	version, err := s.versionRepo.GetByQuotationAndVersion(ctx, quotationID, versionNumber)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusNotFound, "Version not found", err)
	}

	if version.ApprovalStatus != string(shared.ApprovalStatusPending) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot approve version in %s status", version.ApprovalStatus),
			nil,
		)
	}

	// Fetch approver's approval limits
	approverLimit, limitErr := s.approvalRepo.GetByRole(ctx, approverRole)
	if limitErr != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](
			http.StatusForbidden, "No approval limits configured for role: "+approverRole, limitErr)
	}

	// Validate discount within limits
	if version.DiscountValue > 0 {
		if version.DiscountType == string(shared.DiscountTypePercentage) {
			if approverLimit.MaxDiscountPercentage > 0 && version.DiscountValue > approverLimit.MaxDiscountPercentage {
				msg := fmt.Sprintf("Discount of %d bps exceeds your limit of %d bps", version.DiscountValue, approverLimit.MaxDiscountPercentage)
				if approverLimit.EscalationRole != "" {
					msg += fmt.Sprintf(". Requires escalation to %s", approverLimit.EscalationRole)
				}
				return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusForbidden, msg, nil)
			}
		} else if version.DiscountType == string(shared.DiscountTypeFixed) {
			if approverLimit.MaxDiscountAmount > 0 && version.DiscountValue > approverLimit.MaxDiscountAmount {
				msg := fmt.Sprintf("Discount of %d exceeds your limit of %d", version.DiscountValue, approverLimit.MaxDiscountAmount)
				if approverLimit.EscalationRole != "" {
					msg += fmt.Sprintf(". Requires escalation to %s", approverLimit.EscalationRole)
				}
				return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusForbidden, msg, nil)
			}
		}
	}

	// Validate loading within limits
	if version.LoadingValue > 0 {
		if version.LoadingType == string(shared.LoadingTypePercentage) {
			if approverLimit.MaxLoadingPercentage > 0 && version.LoadingValue > approverLimit.MaxLoadingPercentage {
				msg := fmt.Sprintf("Loading of %d bps exceeds your limit of %d bps", version.LoadingValue, approverLimit.MaxLoadingPercentage)
				if approverLimit.EscalationRole != "" {
					msg += fmt.Sprintf(". Requires escalation to %s", approverLimit.EscalationRole)
				}
				return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusForbidden, msg, nil)
			}
		} else if version.LoadingType == string(shared.LoadingTypeFixed) {
			if approverLimit.MaxLoadingAmount > 0 && version.LoadingValue > approverLimit.MaxLoadingAmount {
				msg := fmt.Sprintf("Loading of %d exceeds your limit of %d", version.LoadingValue, approverLimit.MaxLoadingAmount)
				if approverLimit.EscalationRole != "" {
					msg += fmt.Sprintf(". Requires escalation to %s", approverLimit.EscalationRole)
				}
				return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusForbidden, msg, nil)
			}
		}
	}

	updated, err := s.versionRepo.UpdateApprovalStatus(ctx, version.ID, string(shared.ApprovalStatusApproved), approvedBy)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusInternalServerError, "Failed to approve version", err)
	}

	s.logAudit(ctx, approvedBy, string(shared.AuditEntityTypeQuotationVersion), version.ID, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationVersionResponse(updated), http.StatusOK, "Version approved")
}

func (s *quotationServiceImpl) RejectVersion(ctx context.Context, quotationID uuid.UUID, versionNumber int, req salesSchema.RejectVersionRequest, rejectedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationVersionResponse] {
	version, err := s.versionRepo.GetByQuotationAndVersion(ctx, quotationID, versionNumber)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusNotFound, "Version not found", err)
	}

	if version.ApprovalStatus != string(shared.ApprovalStatusPending) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot reject version in %s status", version.ApprovalStatus),
			nil,
		)
	}

	updated, err := s.versionRepo.RejectVersion(ctx, version.ID, req.Reason, rejectedBy)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationVersionResponse](http.StatusInternalServerError, "Failed to reject version", err)
	}

	s.logAudit(ctx, rejectedBy, string(shared.AuditEntityTypeQuotationVersion), version.ID, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationVersionResponse(updated), http.StatusOK, "Version rejected")
}

func (s *quotationServiceImpl) IssueQuotation(ctx context.Context, id uuid.UUID, issuedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusNotFound, "Quotation not found", err)
	}

	if quotation.Status != string(shared.QuotationStatusDraft) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot issue quotation in %s status", quotation.Status),
			nil,
		)
	}

	// Check current version approval
	currentVersion, err := s.versionRepo.GetByQuotationAndVersion(ctx, id, quotation.CurrentVersion)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to get current version", err)
	}

	if currentVersion.RequiresApproval && currentVersion.ApprovalStatus != string(shared.ApprovalStatusApproved) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](
			http.StatusBadRequest,
			"Current version requires approval before issuing",
			nil,
		)
	}

	updated, err := s.quotationRepo.UpdateStatus(ctx, id, string(shared.QuotationStatusIssued))
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to issue quotation", err)
	}

	s.logAudit(ctx, issuedBy, string(shared.AuditEntityTypeQuotation), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationResponse(updated), http.StatusOK, "Quotation issued")
}

func (s *quotationServiceImpl) AcceptQuotation(ctx context.Context, id uuid.UUID, acceptedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusNotFound, "Quotation not found", err)
	}

	if quotation.Status != string(shared.QuotationStatusPendingDecision) && quotation.Status != string(shared.QuotationStatusIssued) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot accept quotation in %s status", quotation.Status),
			nil,
		)
	}

	updated, err := s.quotationRepo.UpdateStatus(ctx, id, string(shared.QuotationStatusAccepted))
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to accept quotation", err)
	}

	s.logAudit(ctx, acceptedBy, string(shared.AuditEntityTypeQuotation), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationResponse(updated), http.StatusOK, "Quotation accepted")
}

func (s *quotationServiceImpl) DeclineQuotation(ctx context.Context, id uuid.UUID, declinedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusNotFound, "Quotation not found", err)
	}

	if quotation.Status != string(shared.QuotationStatusPendingDecision) && quotation.Status != string(shared.QuotationStatusIssued) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot decline quotation in %s status", quotation.Status),
			nil,
		)
	}

	updated, err := s.quotationRepo.UpdateStatus(ctx, id, string(shared.QuotationStatusDeclined))
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to decline quotation", err)
	}

	s.logAudit(ctx, declinedBy, string(shared.AuditEntityTypeQuotation), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationResponse(updated), http.StatusOK, "Quotation declined")
}

func (s *quotationServiceImpl) ExpireQuotations(ctx context.Context) *schema.ServiceResponse[int] {
	expired, err := s.quotationRepo.ListExpired(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int](http.StatusInternalServerError, "Failed to list expired quotations", err)
	}

	count := 0
	for _, q := range expired {
		if _, updateErr := s.quotationRepo.UpdateStatus(ctx, q.ID, string(shared.QuotationStatusExpired)); updateErr != nil {
			log.Printf("Failed to expire quotation %s: %v", q.QuotationNumber, updateErr)
		} else {
			count++
			s.logAudit(ctx, uuid.Nil, string(shared.AuditEntityTypeQuotation), q.ID, string(shared.AuditActionStateChange))
		}
	}

	return schema.NewServiceResponse(count, http.StatusOK, fmt.Sprintf("Expired %d quotations", count))
}

func (s *quotationServiceImpl) SendToClient(ctx context.Context, id uuid.UUID, req salesSchema.SendQuotationRequest, sentBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusNotFound, "Quotation not found", err)
	}

	if quotation.Status != string(shared.QuotationStatusIssued) {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot send quotation in %s status; must be ISSUED", quotation.Status),
			nil,
		)
	}

	// Send notification via NotificationService
	if s.notifSvc != nil {
		subject := fmt.Sprintf("Quotation %s", quotation.QuotationNumber)
		body := fmt.Sprintf("Dear %s, your insurance quotation %s is ready for your review. %s",
			quotation.ClientName, quotation.QuotationNumber, req.Message)
		notifResp := s.notifSvc.Send(ctx, sentBy, req.Channel, string(shared.NotificationTypeQuotation), subject, body)
		if notifResp.Error != nil {
			log.Printf("Warning: Failed to send notification: %v", notifResp.Error)
		}
	}

	updated, err := s.quotationRepo.UpdateStatus(ctx, id, string(shared.QuotationStatusPendingDecision))
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationResponse](http.StatusInternalServerError, "Failed to update quotation status", err)
	}

	s.logAudit(ctx, sentBy, string(shared.AuditEntityTypeQuotation), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ToQuotationResponse(updated), http.StatusOK, "Quotation sent to client")
}

func (s *quotationServiceImpl) ConvertToPolicy(ctx context.Context, id uuid.UUID, req salesSchema.ConvertToPolicyRequest, convertedBy uuid.UUID) *schema.ServiceResponse[salesSchema.ConversionResultResponse] {
	quotation, err := s.quotationRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ConversionResultResponse](http.StatusNotFound, "Quotation not found", err)
	}

	if quotation.Status != string(shared.QuotationStatusAccepted) {
		return schema.NewServiceErrorResponse[salesSchema.ConversionResultResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot convert quotation in %s status; must be %s", quotation.Status, string(shared.QuotationStatusAccepted)),
			nil,
		)
	}

	// Get current version for premium/member data
	currentVersion, err := s.versionRepo.GetByQuotationAndVersion(ctx, id, quotation.CurrentVersion)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ConversionResultResponse](http.StatusInternalServerError, "Failed to get current quotation version", err)
	}

	// Parse start date
	startDate, err := utils.ParseFlexibleDate(req.StartDate)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ConversionResultResponse](http.StatusBadRequest, "Invalid start date format, use YYYY-MM-DD or ISO 8601", err)
	}
	endDate := startDate.AddDate(1, 0, 0)

	// 1. Create Policy via PolicyService
	policyReq := policySchema.CreatePolicyRequest{
		PlanID:            quotation.PlanID.String(),
		PolicyholderName:  quotation.ClientName,
		PolicyholderEmail: quotation.ClientEmail,
		PolicyholderPhone: quotation.ClientPhone,
		StartDate:         startDate.Format("2006-01-02"),
		EndDate:           endDate.Format("2006-01-02"),
	}
	policyResp := s.policySvc.CreatePolicy(ctx, policyReq, convertedBy)
	if policyResp.Error != nil {
		return schema.NewServiceErrorResponse[salesSchema.ConversionResultResponse](
			http.StatusInternalServerError, "Failed to create policy: "+policyResp.Message, policyResp.Error)
	}
	policyID := policyResp.Data.ID
	policyNumber := policyResp.Data.PolicyNumber

	// 2. Enroll proposed members via MemberService
	var members []struct {
		Name         string `json:"name"`
		DateOfBirth  string `json:"date_of_birth"`
		Gender       string `json:"gender"`
		Relationship string `json:"relationship"`
		Phone        string `json:"phone"`
		Email        string `json:"email"`
	}
	if jsonErr := json.Unmarshal(currentVersion.ProposedMembers, &members); jsonErr == nil {
		for _, m := range members {
			memberReq := policySchema.EnrollMemberRequest{
				Name:         m.Name,
				DateOfBirth:  m.DateOfBirth,
				Gender:       m.Gender,
				Relationship: m.Relationship,
				Phone:        m.Phone,
				Email:        m.Email,
			}
			memberResp := s.memberSvc.EnrollMember(ctx, policyID, memberReq)
			if memberResp.Error != nil {
				log.Printf("Warning: Failed to enroll member %s: %v", m.Name, memberResp.Error)
			}
		}
	}

	// 3. Create installment schedule via InstallmentService
	scheduleReq := billingSchema.CreateInstallmentScheduleRequest{
		PolicyID:  policyID.String(),
		Frequency: currentVersion.BillingFrequency,
		StartDate: startDate,
	}
	scheduleResp := s.installmentSvc.CreateSchedule(ctx, scheduleReq, convertedBy)
	if scheduleResp.Error != nil {
		log.Printf("Warning: Failed to create installment schedule: %v", scheduleResp.Error)
	}

	// 4. Update quotation to CONVERTED with real policy ID
	updated, err := s.quotationRepo.SetPolicyID(ctx, id, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.ConversionResultResponse](http.StatusInternalServerError, "Failed to update quotation", err)
	}

	// 5. Update lead to WON
	if _, leadErr := s.leadRepo.UpdateStatus(ctx, quotation.LeadID, string(shared.LeadStatusWon)); leadErr != nil {
		log.Printf("Failed to update lead status to WON: %v", leadErr)
	}

	s.logAudit(ctx, convertedBy, string(shared.AuditEntityTypeQuotation), id, string(shared.AuditActionStateChange))

	return schema.NewServiceResponse(salesSchema.ConversionResultResponse{
		QuotationID:     updated.ID,
		PolicyID:        policyID,
		QuotationNumber: updated.QuotationNumber,
		PolicyNumber:    policyNumber,
		Message:         "Quotation converted to policy successfully",
	}, http.StatusOK, "Quotation converted to policy")
}

func (s *quotationServiceImpl) UploadDocument(ctx context.Context, quotationID uuid.UUID, meta salesSchema.UploadDocumentMeta, s3Key string, uploadedBy uuid.UUID) *schema.ServiceResponse[salesSchema.QuotationDocumentResponse] {
	// Verify quotation exists
	_, err := s.quotationRepo.GetByID(ctx, quotationID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDocumentResponse](http.StatusNotFound, "Quotation not found", err)
	}

	canEditRoles := meta.CanEditRoles
	if canEditRoles == nil {
		canEditRoles = json.RawMessage(fmt.Sprintf(`["%s"]`, shared.UserRoleAdmin))
	}
	canDeleteRoles := meta.CanDeleteRoles
	if canDeleteRoles == nil {
		canDeleteRoles = json.RawMessage(fmt.Sprintf(`["%s"]`, shared.UserRoleAdmin))
	}

	doc := &entity.QuotationDocument{
		QuotationID:    quotationID,
		VersionNumber:  meta.VersionNumber,
		FileName:       meta.FileName,
		FileType:       meta.FileType,
		FileSize:       meta.FileSize,
		S3Key:          s3Key,
		UploadedBy:     uploadedBy,
		CanEditRoles:   canEditRoles,
		CanDeleteRoles: canDeleteRoles,
	}

	created, err := s.documentRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDocumentResponse](http.StatusInternalServerError, "Failed to upload document", err)
	}

	s.logAudit(ctx, uploadedBy, string(shared.AuditEntityTypeQuotationDocument), created.ID, string(shared.AuditActionCreate))

	return schema.NewServiceResponse(salesSchema.ToQuotationDocumentResponse(created), http.StatusCreated, "Document uploaded")
}

func (s *quotationServiceImpl) ListDocuments(ctx context.Context, quotationID uuid.UUID) *schema.ServiceResponse[[]salesSchema.QuotationDocumentResponse] {
	docs, err := s.documentRepo.ListByQuotation(ctx, quotationID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]salesSchema.QuotationDocumentResponse](http.StatusInternalServerError, "Failed to list documents", err)
	}
	responses := make([]salesSchema.QuotationDocumentResponse, len(docs))
	for i, d := range docs {
		responses[i] = salesSchema.ToQuotationDocumentResponse(d)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Documents retrieved")
}

func (s *quotationServiceImpl) UpdateDocument(ctx context.Context, docID uuid.UUID, req salesSchema.UpdateDocumentMeta, updatedBy uuid.UUID, userRole string) *schema.ServiceResponse[salesSchema.QuotationDocumentResponse] {
	doc, err := s.documentRepo.GetByID(ctx, docID)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDocumentResponse](http.StatusNotFound, "Document not found", err)
	}

	// Check role permissions
	var editRoles []string
	if err := json.Unmarshal(doc.CanEditRoles, &editRoles); err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDocumentResponse](http.StatusForbidden, "Insufficient role to edit this document", nil)
	}
	hasPermission := false
	for _, role := range editRoles {
		if role == userRole {
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDocumentResponse](http.StatusForbidden, "Insufficient role to edit this document", nil)
	}

	// Apply updates
	if req.FileName != "" {
		doc.FileName = req.FileName
	}
	if req.CanEditRoles != nil {
		doc.CanEditRoles = req.CanEditRoles
	}
	if req.CanDeleteRoles != nil {
		doc.CanDeleteRoles = req.CanDeleteRoles
	}

	updated, err := s.documentRepo.Update(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[salesSchema.QuotationDocumentResponse](http.StatusInternalServerError, "Failed to update document", err)
	}

	s.logAudit(ctx, updatedBy, string(shared.AuditEntityTypeQuotationDocument), docID, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(salesSchema.ToQuotationDocumentResponse(updated), http.StatusOK, "Document updated")
}

func (s *quotationServiceImpl) DeleteDocument(ctx context.Context, docID uuid.UUID, deletedBy uuid.UUID, userRole string) *schema.ServiceResponse[bool] {
	doc, err := s.documentRepo.GetByID(ctx, docID)
	if err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusNotFound, "Document not found", err)
	}

	// Check role permissions
	var deleteRoles []string
	if err := json.Unmarshal(doc.CanDeleteRoles, &deleteRoles); err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusForbidden, "Insufficient role to delete this document", nil)
	}
	hasDeletePermission := false
	for _, role := range deleteRoles {
		if role == userRole {
			hasDeletePermission = true
			break
		}
	}
	if !hasDeletePermission {
		return schema.NewServiceErrorResponse[bool](http.StatusForbidden, "Insufficient role to delete this document", nil)
	}

	if err := s.documentRepo.SoftDelete(ctx, docID); err != nil {
		return schema.NewServiceErrorResponse[bool](http.StatusInternalServerError, "Failed to delete document", err)
	}

	s.logAudit(ctx, deletedBy, string(shared.AuditEntityTypeQuotationDocument), docID, string(shared.AuditActionDelete))

	return schema.NewServiceResponse(true, http.StatusOK, "Document deleted")
}

func (s *quotationServiceImpl) GetTotalCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.quotationRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get count", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Count retrieved")
}

func (s *quotationServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}

// calculateFinalPremium computes final premium: basePremium - discount + loading
func calculateFinalPremium(basePremium int64, discountType string, discountValue int64, loadingType string, loadingValue int64) int64 {
	var discountAmount int64
	if discountType == string(shared.DiscountTypePercentage) {
		discountAmount = basePremium * discountValue / 10000 // discountValue in basis points (1000 = 10%)
	} else {
		discountAmount = discountValue
	}

	var loadingAmount int64
	if loadingType == string(shared.LoadingTypePercentage) {
		loadingAmount = basePremium * loadingValue / 10000
	} else {
		loadingAmount = loadingValue
	}

	final := basePremium - discountAmount + loadingAmount
	if final < 0 {
		final = 0
	}
	return final
}

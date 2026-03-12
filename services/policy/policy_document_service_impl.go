package policy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	claimRepo "github.com/bitbiz/hias-core/domains/claims/repository"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	notificationService "github.com/bitbiz/hias-core/domains/notification/service"
	"github.com/bitbiz/hias-core/domains/policy/entity"
	"github.com/bitbiz/hias-core/domains/policy/repository"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	preauthRepo "github.com/bitbiz/hias-core/domains/preauth/repository"
	productRepo "github.com/bitbiz/hias-core/domains/product/repository"
	providerRepo "github.com/bitbiz/hias-core/domains/provider/repository"
	"github.com/bitbiz/hias-core/infrastructures/documents"
	"github.com/bitbiz/hias-core/shared"
	awsSvc "github.com/bitbiz/hias-core/shared/aws"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type policyDocumentServiceImpl struct {
	docRepo         repository.PolicyDocumentRepository
	policyRepo      repository.PolicyRepository
	memberRepo      repository.MemberRepository
	endorsementRepo repository.EndorsementRepository
	renewalRepo     repository.PolicyRenewalRepository
	claimRepo       claimRepo.ClaimRepository
	planRepo        productRepo.PlanRepository
	benefitRepo     productRepo.BenefitRepository
	preauthRepo     preauthRepo.PreAuthRepository
	providerRepo    providerRepo.ProviderRepository
	pdfGenerator    documents.PDFGenerator
	s3Svc           awsSvc.S3Service
	auditSvc        auditService.AuditService
	notificationSvc notificationService.NotificationService
}

func NewPolicyDocumentService(
	docRepo repository.PolicyDocumentRepository,
	policyRepo repository.PolicyRepository,
	memberRepo repository.MemberRepository,
	planRepo productRepo.PlanRepository,
	benefitRepo productRepo.BenefitRepository,
	renewalRepo repository.PolicyRenewalRepository,
	preauthRepo preauthRepo.PreAuthRepository,
	providerRepo providerRepo.ProviderRepository,
	pdfGenerator documents.PDFGenerator,
	s3Svc awsSvc.S3Service,
	auditSvc auditService.AuditService,
	notificationSvc notificationService.NotificationService,
	endorsementRepo repository.EndorsementRepository,
	claimRepo claimRepo.ClaimRepository,
) service.PolicyDocumentService {
	return &policyDocumentServiceImpl{
		docRepo:         docRepo,
		policyRepo:      policyRepo,
		memberRepo:      memberRepo,
		planRepo:        planRepo,
		benefitRepo:     benefitRepo,
		renewalRepo:     renewalRepo,
		preauthRepo:     preauthRepo,
		providerRepo:    providerRepo,
		pdfGenerator:    pdfGenerator,
		s3Svc:           s3Svc,
		auditSvc:        auditSvc,
		notificationSvc: notificationSvc,
		endorsementRepo: endorsementRepo,
		claimRepo:       claimRepo,
	}
}

// --- Existing methods (unchanged signatures) ---

func (s *policyDocumentServiceImpl) GenerateWelcomeLetter(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	return s.GenerateDocument(ctx, policySchema.GenerateDocumentRequest{
		EntityType:     string(shared.DocumentEntityTypePolicy),
		EntityID:       policyID.String(),
		DocumentType:   string(shared.PolicyDocumentTypeWelcomeLetter),
		GenerationMode: string(shared.GenerationModeManual),
		GeneratedBy:    generatedBy,
	})
}

func (s *policyDocumentServiceImpl) GenerateMemberCard(ctx context.Context, memberID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	return s.GenerateDocument(ctx, policySchema.GenerateDocumentRequest{
		EntityType:     string(shared.DocumentEntityTypeMember),
		EntityID:       memberID.String(),
		DocumentType:   string(shared.PolicyDocumentTypeMemberCard),
		GenerationMode: string(shared.GenerationModeManual),
		GeneratedBy:    generatedBy,
	})
}

func (s *policyDocumentServiceImpl) GeneratePolicySchedule(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	return s.GenerateDocument(ctx, policySchema.GenerateDocumentRequest{
		EntityType:     string(shared.DocumentEntityTypePolicy),
		EntityID:       policyID.String(),
		DocumentType:   string(shared.PolicyDocumentTypePolicySchedule),
		GenerationMode: string(shared.GenerationModeManual),
		GeneratedBy:    generatedBy,
	})
}

func (s *policyDocumentServiceImpl) GenerateRenewalNotice(ctx context.Context, renewalID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	return s.GenerateDocument(ctx, policySchema.GenerateDocumentRequest{
		EntityType:     string(shared.DocumentEntityTypeRenewal),
		EntityID:       renewalID.String(),
		DocumentType:   string(shared.PolicyDocumentTypeRenewalNotice),
		GenerationMode: string(shared.GenerationModeManual),
		GeneratedBy:    generatedBy,
	})
}

func (s *policyDocumentServiceImpl) GenerateLOU(ctx context.Context, preauthID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	return s.GenerateDocument(ctx, policySchema.GenerateDocumentRequest{
		EntityType:     string(shared.DocumentEntityTypePreauth),
		EntityID:       preauthID.String(),
		DocumentType:   string(shared.PolicyDocumentTypeLOU),
		GenerationMode: string(shared.GenerationModeManual),
		GeneratedBy:    generatedBy,
	})
}

func (s *policyDocumentServiceImpl) GenerateDeclineLetter(ctx context.Context, policyID uuid.UUID, memberName, claimNumber, rejectionReason string, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	// For decline letters, we still need the policy-based generation since
	// the claim handler already resolves the claim and passes policyID.
	// This preserves backwards compatibility.
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Policy not found", err)
	}

	pdfBytes, err := s.pdfGenerator.GenerateDeclineLetter(pol, memberName, claimNumber, rejectionReason)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate decline letter PDF", err)
	}

	version, _ := s.docRepo.GetNextVersion(ctx, string(shared.DocumentEntityTypePolicy), policyID, string(shared.PolicyDocumentTypeDeclineLetter))
	if version == 0 {
		version = 1
	}

	s3Key := fmt.Sprintf("policies/%s/documents/decline_letter_%s.pdf", policyID, uuid.New().String())
	fileName := fmt.Sprintf("Decline_Letter_%s_%s.pdf", claimNumber, pol.PolicyNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:       policyID,
		DocumentType:   string(shared.PolicyDocumentTypeDeclineLetter),
		FileName:       fileName,
		FileSize:       int64(len(pdfBytes)),
		S3Key:          s3Key,
		GeneratedBy:    generatedBy,
		Version:        version,
		Status:         string(shared.GeneratedDocStatusGenerated),
		GenerationMode: string(shared.GenerationModeManual),
		EntityType:     string(shared.DocumentEntityTypePolicy),
		EntityID:       policyID,
	}

	created, err := s.docRepo.CreateV2(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	s.logAudit(ctx, generatedBy, string(shared.AuditEntityTypePolicyDocument), created.ID, string(shared.AuditActionCreate))
	s.sendDocumentNotification(ctx, pol.CreatedBy, "Decline Letter", fileName)

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Decline letter generated")
}

func (s *policyDocumentServiceImpl) ListByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]policySchema.PolicyDocumentResponse] {
	docs, err := s.docRepo.ListByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to list documents", err)
	}

	responses := make([]policySchema.PolicyDocumentResponse, len(docs))
	for i, d := range docs {
		responses[i] = policySchema.ToPolicyDocumentResponse(d)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Documents retrieved")
}

func (s *policyDocumentServiceImpl) GetDocument(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	doc, err := s.docRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Document not found", err)
	}
	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(doc), http.StatusOK, "Document retrieved")
}

func (s *policyDocumentServiceImpl) DeleteDocument(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	if err := s.docRepo.Delete(ctx, id); err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete document", err)
	}
	return schema.NewServiceResponse("Document deleted", http.StatusOK, "Document deleted")
}

func (s *policyDocumentServiceImpl) BulkGenerateMemberCards(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[[]policySchema.PolicyDocumentResponse] {
	members, err := s.memberRepo.ListActiveByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to list members", err)
	}

	var responses []policySchema.PolicyDocumentResponse
	for _, m := range members {
		resp := s.GenerateMemberCard(ctx, m.ID, generatedBy)
		if resp.Error != nil {
			log.Printf("Failed to generate card for member %s: %v", m.MemberNumber, resp.Error)
			continue
		}
		responses = append(responses, resp.Data)
	}

	return schema.NewServiceResponse(responses, http.StatusCreated, fmt.Sprintf("Generated %d member cards", len(responses)))
}

// --- V1 Unified Document Generation ---

func (s *policyDocumentServiceImpl) GenerateDocument(ctx context.Context, req policySchema.GenerateDocumentRequest) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusBadRequest, "Invalid entity ID", err)
	}

	// Readiness check
	readiness := s.CanGenerateDocument(ctx, req.EntityType, entityID, req.DocumentType)
	if readiness.Error != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](readiness.StatusCode, readiness.Message, readiness.Error)
	}
	if !readiness.Data.CanGenerate {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusBadRequest, fmt.Sprintf("Cannot generate document: %s", strings.Join(readiness.Data.Errors, "; ")), nil)
	}

	// Get next version
	version, err := s.docRepo.GetNextVersion(ctx, req.EntityType, entityID, req.DocumentType)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get next version", err)
	}

	generationMode := req.GenerationMode
	if generationMode == "" {
		generationMode = string(shared.GenerationModeManual)
	}

	// Resolve policyID and memberID for backwards compat
	policyID, memberID := s.resolveIDs(ctx, req.EntityType, entityID)

	// Create PENDING record
	pendingDoc := &entity.PolicyDocument{
		PolicyID:       policyID,
		MemberID:       memberID,
		DocumentType:   req.DocumentType,
		FileName:       "pending",
		S3Key:          "",
		GeneratedBy:    req.GeneratedBy,
		Version:        version,
		Status:         string(shared.GeneratedDocStatusPending),
		GenerationMode: generationMode,
		EntityType:     req.EntityType,
		EntityID:       entityID,
	}

	created, err := s.docRepo.CreateV2(ctx, pendingDoc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to create document record", err)
	}

	// Dispatch generation
	pdfBytes, fileName, genErr := s.dispatchGeneration(ctx, req.EntityType, entityID, req.DocumentType)
	if genErr != nil {
		// Mark as FAILED
		s.docRepo.UpdateStatus(ctx, created.ID, string(shared.GeneratedDocStatusFailed), genErr.Error())
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Document generation failed", genErr)
	}

	// Upload to S3
	s3Key := fmt.Sprintf("policies/%s/documents/%s_v%d_%s.pdf", policyID, strings.ToLower(req.DocumentType), version, uuid.New().String())
	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		s.docRepo.UpdateStatus(ctx, created.ID, string(shared.GeneratedDocStatusFailed), err.Error())
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	// Atomically update the PENDING record to GENERATED with file info
	created, err = s.docRepo.UpdateGenerated(ctx, created.ID, fileName, int64(len(pdfBytes)), s3Key, string(shared.GeneratedDocStatusGenerated))
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	// Supersede previous version if exists
	if version > 1 {
		prev, prevErr := s.docRepo.GetLatestByEntity(ctx, req.EntityType, entityID, req.DocumentType)
		if prevErr == nil && prev.ID != created.ID {
			s.docRepo.Supersede(ctx, prev.ID, created.ID)
		}
	}

	// Audit + notification
	s.logAudit(ctx, req.GeneratedBy, string(shared.AuditEntityTypePolicyDocument), created.ID, string(shared.AuditActionCreate))
	if policyID != uuid.Nil {
		pol, polErr := s.policyRepo.GetByID(ctx, policyID)
		if polErr == nil {
			s.sendDocumentNotification(ctx, pol.CreatedBy, req.DocumentType, fileName)
		}
	}

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Document generated")
}

func (s *policyDocumentServiceImpl) CanGenerateDocument(ctx context.Context, entityType string, entityID uuid.UUID, docType string) *schema.ServiceResponse[policySchema.DocumentReadinessResponse] {
	var errors []string

	switch docType {
	case string(shared.PolicyDocumentTypeWelcomeLetter):
		errors = s.validateWelcomeLetter(ctx, entityType, entityID)
	case string(shared.PolicyDocumentTypePolicySchedule):
		errors = s.validatePolicySchedule(ctx, entityType, entityID)
	case string(shared.PolicyDocumentTypeMemberCard):
		errors = s.validateMemberCard(ctx, entityType, entityID)
	case string(shared.PolicyDocumentTypeEndorsement):
		errors = s.validateEndorsementLetter(ctx, entityType, entityID)
	case string(shared.PolicyDocumentTypeRenewalNotice):
		errors = s.validateRenewalNotice(ctx, entityType, entityID)
	case string(shared.PolicyDocumentTypeLOU):
		errors = s.validateLOU(ctx, entityType, entityID)
	case string(shared.PolicyDocumentTypeDeclineLetter):
		errors = s.validateDeclineLetter(ctx, entityType, entityID)
	default:
		errors = []string{fmt.Sprintf("Unknown document type: %s", docType)}
	}

	resp := policySchema.DocumentReadinessResponse{
		CanGenerate: len(errors) == 0,
		Errors:      errors,
	}
	return schema.NewServiceResponse(resp, http.StatusOK, "Readiness check complete")
}

func (s *policyDocumentServiceImpl) GetDocumentAvailability(ctx context.Context, entityType string, entityID uuid.UUID) *schema.ServiceResponse[[]policySchema.DocumentAvailabilityItem] {
	docTypes := s.applicableDocTypes(entityType)
	items := make([]policySchema.DocumentAvailabilityItem, 0, len(docTypes))

	for _, dt := range docTypes {
		item := policySchema.DocumentAvailabilityItem{
			DocumentType: dt,
		}

		// Check readiness
		readiness := s.CanGenerateDocument(ctx, entityType, entityID, dt)
		if readiness.Error == nil {
			item.CanGenerate = readiness.Data.CanGenerate
			item.Errors = readiness.Data.Errors
		}

		// Check latest existing doc
		latest, err := s.docRepo.GetLatestByEntity(ctx, entityType, entityID, dt)
		if err == nil && latest != nil {
			item.Exists = true
			item.LatestStatus = latest.Status
			item.LatestVersion = latest.Version
			item.LatestFileURL = latest.S3Key
			t := latest.CreatedAt
			item.GeneratedAt = &t
		}

		items = append(items, item)
	}

	return schema.NewServiceResponse(items, http.StatusOK, "Document availability retrieved")
}

// --- Readiness validators ---

func (s *policyDocumentServiceImpl) validateWelcomeLetter(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	var errors []string
	if entityType != string(shared.DocumentEntityTypePolicy) {
		return []string{"Welcome letter requires entity_type=policy"}
	}
	pol, err := s.policyRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Policy not found"}
	}
	if pol.Status != string(shared.PolicyStatusActive) {
		errors = append(errors, fmt.Sprintf("Policy must be ACTIVE, currently %s", pol.Status))
	}
	members, err := s.memberRepo.ListActiveByPolicy(ctx, entityID)
	if err != nil || len(members) == 0 {
		errors = append(errors, "Policy must have at least one active member")
	}
	return errors
}

func (s *policyDocumentServiceImpl) validatePolicySchedule(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	var errors []string
	if entityType != string(shared.DocumentEntityTypePolicy) {
		return []string{"Policy schedule requires entity_type=policy"}
	}
	pol, err := s.policyRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Policy not found"}
	}
	if pol.Status != string(shared.PolicyStatusDraft) && pol.Status != string(shared.PolicyStatusActive) {
		errors = append(errors, fmt.Sprintf("Policy must be DRAFT or ACTIVE, currently %s", pol.Status))
	}
	if _, err := s.planRepo.GetByID(ctx, pol.PlanID); err != nil {
		errors = append(errors, "Plan not found")
	}
	members, err := s.memberRepo.ListActiveByPolicy(ctx, entityID)
	if err != nil || len(members) == 0 {
		errors = append(errors, "Policy must have at least one member")
	}
	if _, err := s.benefitRepo.ListByPlan(ctx, pol.PlanID); err != nil {
		errors = append(errors, "Failed to load benefits")
	}
	return errors
}

func (s *policyDocumentServiceImpl) validateMemberCard(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	var errors []string
	if entityType != string(shared.DocumentEntityTypeMember) {
		return []string{"Member card requires entity_type=member"}
	}
	member, err := s.memberRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Member not found"}
	}
	if member.Status != string(shared.MemberStatusActive) {
		errors = append(errors, fmt.Sprintf("Member must be ACTIVE, currently %s", member.Status))
	}
	pol, err := s.policyRepo.GetByID(ctx, member.PolicyID)
	if err != nil {
		errors = append(errors, "Policy not found")
	} else if pol.Status != string(shared.PolicyStatusActive) {
		errors = append(errors, fmt.Sprintf("Policy must be ACTIVE, currently %s", pol.Status))
	}
	return errors
}

func (s *policyDocumentServiceImpl) validateEndorsementLetter(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	if entityType != string(shared.DocumentEntityTypeEndorsement) {
		return []string{"Endorsement letter requires entity_type=endorsement"}
	}
	if s.endorsementRepo == nil {
		return []string{"Endorsement repository not configured"}
	}
	endorsement, err := s.endorsementRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Endorsement not found"}
	}
	if endorsement.Status != string(shared.EndorsementStatusApproved) && endorsement.Status != string(shared.EndorsementStatusApplied) {
		return []string{fmt.Sprintf("Endorsement must be APPROVED or APPLIED, currently %s", endorsement.Status)}
	}
	if _, err := s.policyRepo.GetByID(ctx, endorsement.PolicyID); err != nil {
		return []string{"Linked policy not found"}
	}
	return nil
}

func (s *policyDocumentServiceImpl) validateRenewalNotice(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	if entityType != string(shared.DocumentEntityTypeRenewal) {
		return []string{"Renewal notice requires entity_type=renewal"}
	}
	renewal, err := s.renewalRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Renewal not found"}
	}
	if _, err := s.policyRepo.GetByID(ctx, renewal.PolicyID); err != nil {
		return []string{"Linked policy not found"}
	}
	return nil
}

func (s *policyDocumentServiceImpl) validateLOU(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	if entityType != string(shared.DocumentEntityTypePreauth) {
		return []string{"LOU requires entity_type=preauth"}
	}
	preauth, err := s.preauthRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Pre-authorization not found"}
	}
	if preauth.Status != string(shared.PreAuthStatusApproved) {
		return []string{fmt.Sprintf("Pre-authorization must be APPROVED, currently %s", preauth.Status)}
	}
	if _, err := s.policyRepo.GetByID(ctx, preauth.PolicyID); err != nil {
		return []string{"Policy not found"}
	}
	if _, err := s.memberRepo.GetByID(ctx, preauth.MemberID); err != nil {
		return []string{"Member not found"}
	}
	if _, err := s.providerRepo.GetByID(ctx, preauth.ProviderID); err != nil {
		return []string{"Provider not found"}
	}
	return nil
}

func (s *policyDocumentServiceImpl) validateDeclineLetter(ctx context.Context, entityType string, entityID uuid.UUID) []string {
	if entityType != string(shared.DocumentEntityTypeClaim) {
		return []string{"Decline letter requires entity_type=claim"}
	}
	if s.claimRepo == nil {
		return []string{"Claim repository not configured"}
	}
	claim, err := s.claimRepo.GetByID(ctx, entityID)
	if err != nil {
		return []string{"Claim not found"}
	}
	if claim.Status != string(shared.ClaimStatusRejected) {
		return []string{fmt.Sprintf("Claim must be REJECTED, currently %s", claim.Status)}
	}
	if claim.RejectionReason == "" {
		return []string{"Claim has no rejection reason"}
	}
	if _, err := s.policyRepo.GetByID(ctx, claim.PolicyID); err != nil {
		return []string{"Policy not found"}
	}
	return nil
}

// --- Dispatch generation (calls existing PDF generators) ---

func (s *policyDocumentServiceImpl) dispatchGeneration(ctx context.Context, entityType string, entityID uuid.UUID, docType string) ([]byte, string, error) {
	switch docType {
	case string(shared.PolicyDocumentTypeWelcomeLetter):
		return s.generateWelcomeLetterPDF(ctx, entityID)
	case string(shared.PolicyDocumentTypePolicySchedule):
		return s.generatePolicySchedulePDF(ctx, entityID)
	case string(shared.PolicyDocumentTypeMemberCard):
		return s.generateMemberCardPDF(ctx, entityID)
	case string(shared.PolicyDocumentTypeEndorsement):
		return s.generateEndorsementLetterPDF(ctx, entityID)
	case string(shared.PolicyDocumentTypeRenewalNotice):
		return s.generateRenewalNoticePDF(ctx, entityID)
	case string(shared.PolicyDocumentTypeLOU):
		return s.generateLOUPDF(ctx, entityID)
	default:
		return nil, "", fmt.Errorf("unsupported document type: %s", docType)
	}
}

func (s *policyDocumentServiceImpl) generateWelcomeLetterPDF(ctx context.Context, policyID uuid.UUID) ([]byte, string, error) {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return nil, "", err
	}
	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return nil, "", err
	}
	members, err := s.memberRepo.ListActiveByPolicy(ctx, policyID)
	if err != nil {
		return nil, "", err
	}
	pdfBytes, err := s.pdfGenerator.GenerateWelcomeLetter(pol, members, plan.Name)
	if err != nil {
		return nil, "", err
	}
	return pdfBytes, fmt.Sprintf("Welcome_Letter_%s.pdf", pol.PolicyNumber), nil
}

func (s *policyDocumentServiceImpl) generatePolicySchedulePDF(ctx context.Context, policyID uuid.UUID) ([]byte, string, error) {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return nil, "", err
	}
	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return nil, "", err
	}
	members, err := s.memberRepo.ListActiveByPolicy(ctx, policyID)
	if err != nil {
		return nil, "", err
	}
	benefits, err := s.benefitRepo.ListByPlan(ctx, pol.PlanID)
	if err != nil {
		return nil, "", err
	}
	pdfBytes, err := s.pdfGenerator.GeneratePolicySchedule(pol, members, plan, benefits)
	if err != nil {
		return nil, "", err
	}
	return pdfBytes, fmt.Sprintf("Policy_Schedule_%s.pdf", pol.PolicyNumber), nil
}

func (s *policyDocumentServiceImpl) generateMemberCardPDF(ctx context.Context, memberID uuid.UUID) ([]byte, string, error) {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return nil, "", err
	}
	pol, err := s.policyRepo.GetByID(ctx, member.PolicyID)
	if err != nil {
		return nil, "", err
	}
	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return nil, "", err
	}
	pdfBytes, err := s.pdfGenerator.GenerateMemberCard(member, pol, plan.Name)
	if err != nil {
		return nil, "", err
	}
	return pdfBytes, fmt.Sprintf("Member_Card_%s.pdf", member.MemberNumber), nil
}

func (s *policyDocumentServiceImpl) generateEndorsementLetterPDF(ctx context.Context, endorsementID uuid.UUID) ([]byte, string, error) {
	endorsement, err := s.endorsementRepo.GetByID(ctx, endorsementID)
	if err != nil {
		return nil, "", err
	}
	pol, err := s.policyRepo.GetByID(ctx, endorsement.PolicyID)
	if err != nil {
		return nil, "", err
	}
	// Reuse welcome letter generator for endorsement letters (existing pattern)
	members, err := s.memberRepo.ListActiveByPolicy(ctx, endorsement.PolicyID)
	if err != nil {
		return nil, "", err
	}
	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return nil, "", err
	}
	pdfBytes, err := s.pdfGenerator.GenerateWelcomeLetter(pol, members, plan.Name)
	if err != nil {
		return nil, "", err
	}
	return pdfBytes, fmt.Sprintf("Endorsement_Letter_%s.pdf", pol.PolicyNumber), nil
}

func (s *policyDocumentServiceImpl) generateRenewalNoticePDF(ctx context.Context, renewalID uuid.UUID) ([]byte, string, error) {
	renewal, err := s.renewalRepo.GetByID(ctx, renewalID)
	if err != nil {
		return nil, "", err
	}
	pol, err := s.policyRepo.GetByID(ctx, renewal.PolicyID)
	if err != nil {
		return nil, "", err
	}
	pdfBytes, err := s.pdfGenerator.GenerateRenewalNotice(pol, renewal)
	if err != nil {
		return nil, "", err
	}
	return pdfBytes, fmt.Sprintf("Renewal_Notice_%s.pdf", pol.PolicyNumber), nil
}

func (s *policyDocumentServiceImpl) generateLOUPDF(ctx context.Context, preauthID uuid.UUID) ([]byte, string, error) {
	preauth, err := s.preauthRepo.GetByID(ctx, preauthID)
	if err != nil {
		return nil, "", err
	}
	pol, err := s.policyRepo.GetByID(ctx, preauth.PolicyID)
	if err != nil {
		return nil, "", err
	}
	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return nil, "", err
	}
	member, err := s.memberRepo.GetByID(ctx, preauth.MemberID)
	if err != nil {
		return nil, "", err
	}
	provider, err := s.providerRepo.GetByID(ctx, preauth.ProviderID)
	if err != nil {
		return nil, "", err
	}
	pdfBytes, err := s.pdfGenerator.GenerateLOU(preauth, pol, member.Name, provider.Name, plan.Name)
	if err != nil {
		return nil, "", err
	}
	return pdfBytes, fmt.Sprintf("LOU_%s_%s.pdf", preauth.AuthCode, pol.PolicyNumber), nil
}

// --- Upload Flow ---

func (s *policyDocumentServiceImpl) RequestUploadURL(ctx context.Context, policyID uuid.UUID, req policySchema.UploadPolicyDocumentURLRequest, uploadedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentUploadURLResponse] {
	// Validate policy exists
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusNotFound, "Policy not found", err)
	}

	// Validate file inputs
	if err := utils.ValidateFileName(req.FileName); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusBadRequest, err.Error(), err)
	}
	if err := utils.ValidateFileSize(req.FileSize, 50*1024*1024); err != nil { // 50MB max
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusBadRequest, err.Error(), err)
	}
	if err := utils.ValidateMimeType(req.MimeType, nil); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusBadRequest, err.Error(), err)
	}

	// Defaults
	entityType := req.EntityType
	if entityType == "" {
		entityType = string(shared.DocumentEntityTypePolicy)
	}
	entityIDStr := req.EntityID
	if entityIDStr == "" {
		entityIDStr = policyID.String()
	}
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusBadRequest, "Invalid entity_id", err)
	}

	// Get next version
	version, err := s.docRepo.GetNextVersion(ctx, entityType, entityID, req.DocumentType)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusInternalServerError, "Failed to get next version", err)
	}

	// Generate S3 key
	s3Key := utils.GenerateS3Key(entityType, entityIDStr, req.DocumentType, req.FileName)

	// Create DB record with PENDING_UPLOAD status
	doc := &entity.PolicyDocument{
		PolicyID:       pol.ID,
		DocumentType:   req.DocumentType,
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		MimeType:       req.MimeType,
		S3Key:          s3Key,
		GeneratedBy:    uploadedBy,
		Version:        version,
		Status:         string(shared.GeneratedDocStatusPendingUpload),
		GenerationMode: string(shared.GenerationModeUpload),
		EntityType:     entityType,
		EntityID:       entityID,
	}

	created, err := s.docRepo.CreateV2(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusInternalServerError, "Failed to create document record", err)
	}

	// Get presigned PUT URL
	if s.s3Svc == nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusInternalServerError, "S3 service not configured", nil)
	}
	var expiresIn int64 = 900
	uploadURL, err := s.s3Svc.GetPresignedPutURL(ctx, s3Key, req.MimeType, expiresIn)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentUploadURLResponse](http.StatusInternalServerError, "Failed to get presigned URL", err)
	}

	resp := policySchema.PolicyDocumentUploadURLResponse{
		DocumentID: created.ID,
		UploadURL:  uploadURL,
		S3Key:      s3Key,
		ExpiresIn:  expiresIn,
	}

	return schema.NewServiceResponse(resp, http.StatusCreated, "Upload URL generated")
}

func (s *policyDocumentServiceImpl) ConfirmUpload(ctx context.Context, documentID uuid.UUID, uploadedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	// Get document
	doc, err := s.docRepo.GetByID(ctx, documentID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Document not found", err)
	}

	// Validate status
	if doc.Status != string(shared.GeneratedDocStatusPendingUpload) {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusBadRequest, fmt.Sprintf("Document status must be PENDING_UPLOAD, currently %s", doc.Status), nil)
	}

	// Verify file exists in S3
	if s.s3Svc == nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "S3 service not configured", nil)
	}
	fileSize, err := s.s3Svc.HeadObject(ctx, doc.S3Key)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusBadRequest, "File not found in S3. Upload may not be complete.", err)
	}

	// Update document: status → GENERATED, update file_size from HeadObject
	updated, err := s.docRepo.ConfirmUpload(ctx, documentID, fileSize)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to confirm upload", err)
	}

	// Audit log
	s.logAudit(ctx, uploadedBy, string(shared.AuditEntityTypePolicyDocument), documentID, string(shared.AuditActionUpdate))

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(updated), http.StatusOK, "Upload confirmed")
}

// --- Helpers ---

func (s *policyDocumentServiceImpl) resolveIDs(ctx context.Context, entityType string, entityID uuid.UUID) (policyID, memberID uuid.UUID) {
	switch entityType {
	case string(shared.DocumentEntityTypePolicy):
		return entityID, uuid.Nil
	case string(shared.DocumentEntityTypeMember):
		member, err := s.memberRepo.GetByID(ctx, entityID)
		if err != nil {
			return uuid.Nil, entityID
		}
		return member.PolicyID, entityID
	case string(shared.DocumentEntityTypeEndorsement):
		if s.endorsementRepo != nil {
			endorsement, err := s.endorsementRepo.GetByID(ctx, entityID)
			if err == nil {
				return endorsement.PolicyID, uuid.Nil
			}
		}
		return uuid.Nil, uuid.Nil
	case string(shared.DocumentEntityTypeRenewal):
		renewal, err := s.renewalRepo.GetByID(ctx, entityID)
		if err == nil {
			return renewal.PolicyID, uuid.Nil
		}
		return uuid.Nil, uuid.Nil
	case string(shared.DocumentEntityTypePreauth):
		preauth, err := s.preauthRepo.GetByID(ctx, entityID)
		if err == nil {
			return preauth.PolicyID, uuid.Nil
		}
		return uuid.Nil, uuid.Nil
	default:
		return uuid.Nil, uuid.Nil
	}
}

func (s *policyDocumentServiceImpl) applicableDocTypes(entityType string) []string {
	switch entityType {
	case string(shared.DocumentEntityTypePolicy):
		return []string{string(shared.PolicyDocumentTypeWelcomeLetter), string(shared.PolicyDocumentTypePolicySchedule)}
	case string(shared.DocumentEntityTypeMember):
		return []string{string(shared.PolicyDocumentTypeMemberCard)}
	case string(shared.DocumentEntityTypeEndorsement):
		return []string{string(shared.PolicyDocumentTypeEndorsement)}
	case string(shared.DocumentEntityTypeRenewal):
		return []string{string(shared.PolicyDocumentTypeRenewalNotice)}
	case string(shared.DocumentEntityTypePreauth):
		return []string{string(shared.PolicyDocumentTypeLOU)}
	case string(shared.DocumentEntityTypeClaim):
		return []string{string(shared.PolicyDocumentTypeDeclineLetter)}
	default:
		return nil
	}
}

func (s *policyDocumentServiceImpl) uploadToS3(ctx context.Context, key string, pdfBytes []byte) error {
	if s.s3Svc == nil {
		return nil // graceful degradation when S3 is not configured
	}
	_, err := s.s3Svc.Upload(ctx, awsSvc.UploadRequest{
		Key:         key,
		Body:        bytes.NewReader(pdfBytes),
		ContentType: "application/pdf",
	})
	return err
}

func (s *policyDocumentServiceImpl) sendDocumentNotification(ctx context.Context, recipientID uuid.UUID, docType, fileName string) {
	if s.notificationSvc == nil {
		return
	}
	subject := fmt.Sprintf("%s Ready", docType)
	body := fmt.Sprintf("Your %s (%s) has been generated and is available for download.", docType, fileName)
	resp := s.notificationSvc.Send(ctx, recipientID, string(shared.NotificationChannelInApp), string(shared.NotificationTypePolicy), subject, body)
	if resp.Error != nil {
		log.Printf("Warning: failed to send document notification: %v", resp.Error)
	}
}

func (s *policyDocumentServiceImpl) logAudit(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string) {
	if s.auditSvc != nil {
		resp := s.auditSvc.LogEvent(ctx, userID, entityType, entityID, action, nil, nil, "", "")
		if resp.Error != nil {
			log.Printf("Failed to log audit: %v", resp.Error)
		}
	}
}

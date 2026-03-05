package policy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
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
	"github.com/google/uuid"
)

type policyDocumentServiceImpl struct {
	docRepo         repository.PolicyDocumentRepository
	policyRepo      repository.PolicyRepository
	memberRepo      repository.MemberRepository
	planRepo        productRepo.PlanRepository
	benefitRepo     productRepo.BenefitRepository
	renewalRepo     repository.PolicyRenewalRepository
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
	}
}

func (s *policyDocumentServiceImpl) GenerateWelcomeLetter(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Policy not found", err)
	}

	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get plan", err)
	}

	members, err := s.memberRepo.ListActiveByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to list members", err)
	}

	pdfBytes, err := s.pdfGenerator.GenerateWelcomeLetter(pol, members, plan.Name)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate PDF", err)
	}

	s3Key := fmt.Sprintf("policies/%s/documents/welcome_letter_%s.pdf", policyID, uuid.New().String())
	fileName := fmt.Sprintf("Welcome_Letter_%s.pdf", pol.PolicyNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:     policyID,
		DocumentType: string(shared.PolicyDocumentTypeWelcomeLetter),
		FileName:     fileName,
		FileSize:     int64(len(pdfBytes)),
		S3Key:        s3Key,
		GeneratedBy:  generatedBy,
	}

	created, err := s.docRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	s.logAudit(ctx, generatedBy, string(shared.AuditEntityTypePolicyDocument), created.ID, string(shared.AuditActionCreate))
	s.sendDocumentNotification(ctx, pol.CreatedBy, "Welcome Letter", fileName)

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Welcome letter generated")
}

func (s *policyDocumentServiceImpl) GenerateMemberCard(ctx context.Context, memberID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Member not found", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, member.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get policy", err)
	}

	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get plan", err)
	}

	pdfBytes, err := s.pdfGenerator.GenerateMemberCard(member, pol, plan.Name)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate PDF", err)
	}

	s3Key := fmt.Sprintf("policies/%s/documents/member_card_%s.pdf", member.PolicyID, uuid.New().String())
	fileName := fmt.Sprintf("Member_Card_%s.pdf", member.MemberNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:     member.PolicyID,
		MemberID:     memberID,
		DocumentType: string(shared.PolicyDocumentTypeMemberCard),
		FileName:     fileName,
		FileSize:     int64(len(pdfBytes)),
		S3Key:        s3Key,
		GeneratedBy:  generatedBy,
	}

	created, err := s.docRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Member card generated")
}

func (s *policyDocumentServiceImpl) GeneratePolicySchedule(ctx context.Context, policyID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Policy not found", err)
	}

	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get plan", err)
	}

	members, err := s.memberRepo.ListActiveByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to list members", err)
	}

	benefits, err := s.benefitRepo.ListByPlan(ctx, pol.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to list benefits", err)
	}

	pdfBytes, err := s.pdfGenerator.GeneratePolicySchedule(pol, members, plan, benefits)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate PDF", err)
	}

	s3Key := fmt.Sprintf("policies/%s/documents/policy_schedule_%s.pdf", policyID, uuid.New().String())
	fileName := fmt.Sprintf("Policy_Schedule_%s.pdf", pol.PolicyNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:     policyID,
		DocumentType: string(shared.PolicyDocumentTypePolicySchedule),
		FileName:     fileName,
		FileSize:     int64(len(pdfBytes)),
		S3Key:        s3Key,
		GeneratedBy:  generatedBy,
	}

	created, err := s.docRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	s.sendDocumentNotification(ctx, pol.CreatedBy, "Policy Schedule", fileName)

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Policy schedule generated")
}

func (s *policyDocumentServiceImpl) GenerateRenewalNotice(ctx context.Context, renewalID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	renewal, err := s.renewalRepo.GetByID(ctx, renewalID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Renewal not found", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, renewal.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get policy", err)
	}

	pdfBytes, err := s.pdfGenerator.GenerateRenewalNotice(pol, renewal)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate PDF", err)
	}

	s3Key := fmt.Sprintf("policies/%s/documents/renewal_notice_%s.pdf", pol.ID, uuid.New().String())
	fileName := fmt.Sprintf("Renewal_Notice_%s.pdf", pol.PolicyNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:     pol.ID,
		DocumentType: string(shared.PolicyDocumentTypeRenewalNotice),
		FileName:     fileName,
		FileSize:     int64(len(pdfBytes)),
		S3Key:        s3Key,
		GeneratedBy:  generatedBy,
	}

	created, err := s.docRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Renewal notice generated")
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

func (s *policyDocumentServiceImpl) GenerateLOU(ctx context.Context, preauthID uuid.UUID, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	preauth, err := s.preauthRepo.GetByID(ctx, preauthID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Pre-authorization not found", err)
	}

	if preauth.Status != string(shared.PreAuthStatusApproved) {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusBadRequest, "Pre-authorization must be approved to generate LOU", nil)
	}

	// Check for existing LOU for same preauth
	existingDocs, _ := s.docRepo.ListByPolicy(ctx, preauth.PolicyID)
	expectedPrefix := fmt.Sprintf("LOU_%s_", preauth.AuthCode)
	for _, doc := range existingDocs {
		if doc.DocumentType == string(shared.PolicyDocumentTypeLOU) && strings.HasPrefix(doc.FileName, expectedPrefix) {
			return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(doc), http.StatusOK,
				"Existing LOU found for this pre-authorization (generated on "+doc.CreatedAt.Format("2006-01-02")+")")
		}
	}

	pol, err := s.policyRepo.GetByID(ctx, preauth.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get policy", err)
	}

	plan, err := s.planRepo.GetByID(ctx, pol.PlanID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get plan", err)
	}

	member, err := s.memberRepo.GetByID(ctx, preauth.MemberID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get member", err)
	}

	provider, err := s.providerRepo.GetByID(ctx, preauth.ProviderID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to get provider", err)
	}

	pdfBytes, err := s.pdfGenerator.GenerateLOU(preauth, pol, member.Name, provider.Name, plan.Name)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate LOU PDF", err)
	}

	s3Key := fmt.Sprintf("policies/%s/documents/lou_%s.pdf", pol.ID, uuid.New().String())
	fileName := fmt.Sprintf("LOU_%s_%s.pdf", preauth.AuthCode, pol.PolicyNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload LOU to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:     pol.ID,
		DocumentType: string(shared.PolicyDocumentTypeLOU),
		FileName:     fileName,
		FileSize:     int64(len(pdfBytes)),
		S3Key:        s3Key,
		GeneratedBy:  generatedBy,
	}

	created, err := s.docRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save LOU document record", err)
	}

	s.logAudit(ctx, generatedBy, string(shared.AuditEntityTypePolicyDocument), created.ID, string(shared.AuditActionCreate))
	s.sendDocumentNotification(ctx, pol.CreatedBy, "Letter of Undertaking", fileName)

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Letter of Undertaking generated")
}

func (s *policyDocumentServiceImpl) GenerateDeclineLetter(ctx context.Context, policyID uuid.UUID, memberName, claimNumber, rejectionReason string, generatedBy uuid.UUID) *schema.ServiceResponse[policySchema.PolicyDocumentResponse] {
	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusNotFound, "Policy not found", err)
	}

	pdfBytes, err := s.pdfGenerator.GenerateDeclineLetter(pol, memberName, claimNumber, rejectionReason)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to generate decline letter PDF", err)
	}

	s3Key := fmt.Sprintf("policies/%s/documents/decline_letter_%s.pdf", policyID, uuid.New().String())
	fileName := fmt.Sprintf("Decline_Letter_%s_%s.pdf", claimNumber, pol.PolicyNumber)

	if err := s.uploadToS3(ctx, s3Key, pdfBytes); err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to upload document to S3", err)
	}

	doc := &entity.PolicyDocument{
		PolicyID:     policyID,
		DocumentType: string(shared.PolicyDocumentTypeDeclineLetter),
		FileName:     fileName,
		FileSize:     int64(len(pdfBytes)),
		S3Key:        s3Key,
		GeneratedBy:  generatedBy,
	}

	created, err := s.docRepo.Create(ctx, doc)
	if err != nil {
		return schema.NewServiceErrorResponse[policySchema.PolicyDocumentResponse](http.StatusInternalServerError, "Failed to save document record", err)
	}

	s.logAudit(ctx, generatedBy, string(shared.AuditEntityTypePolicyDocument), created.ID, string(shared.AuditActionCreate))
	s.sendDocumentNotification(ctx, pol.CreatedBy, "Decline Letter", fileName)

	return schema.NewServiceResponse(policySchema.ToPolicyDocumentResponse(created), http.StatusCreated, "Decline letter generated")
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

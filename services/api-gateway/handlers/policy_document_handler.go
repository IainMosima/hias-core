package handlers

import (
	"net/http"

	claimService "github.com/bitbiz/hias-core/domains/claims/service"
	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PolicyDocumentHandler struct {
	policyDocSvc service.PolicyDocumentService
	claimSvc     claimService.ClaimService
}

func NewPolicyDocumentHandler(policyDocSvc service.PolicyDocumentService, claimSvc claimService.ClaimService) *PolicyDocumentHandler {
	return &PolicyDocumentHandler{policyDocSvc: policyDocSvc, claimSvc: claimSvc}
}

// GenerateWelcomeLetter godoc
// @Summary      Generate a welcome letter
// @Description  Generate a welcome letter document for a policy
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents/welcome-letter [post]
func (h *PolicyDocumentHandler) GenerateWelcomeLetter(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyDocSvc.GenerateWelcomeLetter(ctx.Request.Context(), policyID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GenerateMemberCard godoc
// @Summary      Generate a member card
// @Description  Generate a member card document for a member
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/card [post]
func (h *PolicyDocumentHandler) GenerateMemberCard(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	resp := h.policyDocSvc.GenerateMemberCard(ctx.Request.Context(), memberID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GeneratePolicySchedule godoc
// @Summary      Generate a policy schedule
// @Description  Generate a policy schedule document for a policy
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents/policy-schedule [post]
func (h *PolicyDocumentHandler) GeneratePolicySchedule(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyDocSvc.GeneratePolicySchedule(ctx.Request.Context(), policyID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *PolicyDocumentHandler) GenerateRenewalNotice(ctx *gin.Context) {
	renewalID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid renewal ID")
		return
	}

	resp := h.policyDocSvc.GenerateRenewalNotice(ctx.Request.Context(), renewalID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListDocuments godoc
// @Summary      List policy documents
// @Description  List all documents associated with a policy
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents [get]
func (h *PolicyDocumentHandler) ListDocuments(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyDocSvc.ListByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetDocument godoc
// @Summary      Get a policy document
// @Description  Get a policy document by ID
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policy-documents/{id} [get]
func (h *PolicyDocumentHandler) GetDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.policyDocSvc.GetDocument(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteDocument godoc
// @Summary      Delete a policy document
// @Description  Delete a policy document by ID
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policy-documents/{id} [delete]
func (h *PolicyDocumentHandler) DeleteDocument(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.policyDocSvc.DeleteDocument(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeletePolicyDocument godoc
// @Summary      Delete a policy document (nested)
// @Description  Delete a policy document by ID from the nested policy path
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id    path string true "Policy ID"
// @Param        docId path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents/{docId} [delete]
func (h *PolicyDocumentHandler) DeletePolicyDocument(ctx *gin.Context) {
	docID, err := uuid.Parse(ctx.Param("docId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.policyDocSvc.DeleteDocument(ctx.Request.Context(), docID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkGenerateMemberCards godoc
// @Summary      Bulk generate member cards
// @Description  Generate member cards for all members in a policy
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents/member-cards [post]
func (h *PolicyDocumentHandler) BulkGenerateMemberCards(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.policyDocSvc.BulkGenerateMemberCards(ctx.Request.Context(), policyID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GenerateLOU godoc
// @Summary      Generate a Letter of Undertaking
// @Description  Generate a Letter of Undertaking (LOU) for a pre-authorization
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Pre-authorization ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/preauths/{id}/lou [post]
func (h *PolicyDocumentHandler) GenerateLOU(ctx *gin.Context) {
	preauthID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid pre-authorization ID")
		return
	}

	resp := h.policyDocSvc.GenerateLOU(ctx.Request.Context(), preauthID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GenerateDeclineLetter godoc
// @Summary      Generate a decline letter
// @Description  Generate a decline letter for a rejected claim
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id path string true "Claim ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/claims/{id}/decline-letter [post]
func (h *PolicyDocumentHandler) GenerateDeclineLetter(ctx *gin.Context) {
	claimID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	// Fetch claim details
	claimResp := h.claimSvc.GetClaim(ctx.Request.Context(), claimID)
	if claimResp.Error != nil {
		utils.RespondError(ctx, claimResp.StatusCode, claimResp.Message)
		return
	}
	claim := claimResp.Data

	if claim.RejectionReason == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "Claim has no rejection reason; must be rejected before generating decline letter")
		return
	}

	resp := h.policyDocSvc.GenerateDeclineLetter(ctx.Request.Context(), claim.PolicyID, "", claim.ClaimNumber, claim.RejectionReason, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// --- Upload Flow Endpoints ---

// RequestUploadURL godoc
// @Summary      Request a presigned upload URL for a policy document
// @Description  Creates a pending document record and returns a presigned PUT URL
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id   path string true "Policy ID"
// @Param        body body policySchema.UploadPolicyDocumentURLRequest true "Upload request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents/upload-url [post]
func (h *PolicyDocumentHandler) RequestUploadURL(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	var req policySchema.UploadPolicyDocumentURLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp := h.policyDocSvc.RequestUploadURL(ctx.Request.Context(), policyID, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ConfirmUpload godoc
// @Summary      Confirm a policy document upload
// @Description  Verifies the file exists in S3 and updates the document status to GENERATED
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        id    path string true "Policy ID"
// @Param        docId path string true "Document ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/documents/{docId}/confirm-upload [post]
func (h *PolicyDocumentHandler) ConfirmUpload(ctx *gin.Context) {
	docID, err := uuid.Parse(ctx.Param("docId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid document ID")
		return
	}

	resp := h.policyDocSvc.ConfirmUpload(ctx.Request.Context(), docID, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// --- V1 Unified Document Generation Endpoints ---

// GenerateDocument godoc
// @Summary      Generate a document
// @Description  Generate any document type for a given entity
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        body body policySchema.GenerateDocumentRequest true "Generate document request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/generate [post]
func (h *PolicyDocumentHandler) GenerateDocument(ctx *gin.Context) {
	var req policySchema.GenerateDocumentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	req.GenerationMode = "MANUAL"
	req.GeneratedBy = getUserID(ctx)

	resp := h.policyDocSvc.GenerateDocument(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CanGenerate godoc
// @Summary      Check if a document can be generated
// @Description  Check readiness for generating a specific document type
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        entity_type query string true "Entity type"
// @Param        entity_id query string true "Entity ID"
// @Param        document_type query string true "Document type"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/can-generate [get]
func (h *PolicyDocumentHandler) CanGenerate(ctx *gin.Context) {
	entityType := ctx.Query("entity_type")
	entityIDStr := ctx.Query("entity_id")
	docType := ctx.Query("document_type")

	if entityType == "" || entityIDStr == "" || docType == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "entity_type, entity_id, and document_type are required")
		return
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid entity_id")
		return
	}

	resp := h.policyDocSvc.CanGenerateDocument(ctx.Request.Context(), entityType, entityID, docType)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetAvailability godoc
// @Summary      Get document availability
// @Description  Get availability of all applicable document types for an entity
// @Tags         PolicyDocuments
// @Accept       json
// @Produce      json
// @Param        entity_type query string true "Entity type"
// @Param        entity_id query string true "Entity ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/documents/availability [get]
func (h *PolicyDocumentHandler) GetAvailability(ctx *gin.Context) {
	entityType := ctx.Query("entity_type")
	entityIDStr := ctx.Query("entity_id")

	if entityType == "" || entityIDStr == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "entity_type and entity_id are required")
		return
	}

	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid entity_id")
		return
	}

	resp := h.policyDocSvc.GetDocumentAvailability(ctx.Request.Context(), entityType, entityID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

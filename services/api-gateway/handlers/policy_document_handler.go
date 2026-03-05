package handlers

import (
	"net/http"

	claimService "github.com/bitbiz/hias-core/domains/claims/service"
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

package handlers

import (
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnderwritingHandler struct {
	underwritingSvc service.UnderwritingService
}

func NewUnderwritingHandler(underwritingSvc service.UnderwritingService) *UnderwritingHandler {
	return &UnderwritingHandler{underwritingSvc: underwritingSvc}
}

// SubmitAssessment godoc
// @Summary      Submit an underwriting assessment
// @Description  Submit a new underwriting assessment for a policy
// @Tags         Underwriting
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.SubmitAssessmentRequest true "Assessment details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/underwriting [post]
func (h *UnderwritingHandler) SubmitAssessment(ctx *gin.Context) {
	var req policySchema.SubmitAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Use policy ID from URL path
	req.PolicyID = ctx.Param("id")

	resp := h.underwritingSvc.SubmitAssessment(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetAssessment godoc
// @Summary      Get an underwriting assessment
// @Description  Get underwriting assessment details by ID
// @Tags         Underwriting
// @Accept       json
// @Produce      json
// @Param        id path string true "Assessment ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting/{id} [get]
func (h *UnderwritingHandler) GetAssessment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid assessment ID")
		return
	}

	resp := h.underwritingSvc.GetAssessment(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListAssessments godoc
// @Summary      List underwriting assessments
// @Description  List all underwriting assessments for a policy
// @Tags         Underwriting
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/underwriting [get]
func (h *UnderwritingHandler) ListAssessments(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.underwritingSvc.ListByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ReviewAssessment godoc
// @Summary      Review an underwriting assessment
// @Description  Review and update an underwriting assessment decision
// @Tags         Underwriting
// @Accept       json
// @Produce      json
// @Param        id path string true "Assessment ID"
// @Param        request body policySchema.ReviewAssessmentRequest true "Review details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/underwriting/{id}/review [put]
func (h *UnderwritingHandler) ReviewAssessment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid assessment ID")
		return
	}

	var req policySchema.ReviewAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.underwritingSvc.ReviewAssessment(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

package handlers

import (
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EndorsementHandler struct {
	endorsementSvc service.EndorsementService
}

func NewEndorsementHandler(endorsementSvc service.EndorsementService) *EndorsementHandler {
	return &EndorsementHandler{endorsementSvc: endorsementSvc}
}

// CreateEndorsement godoc
// @Summary      Create a new endorsement
// @Description  Create a new endorsement for a policy
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.CreateEndorsementRequest true "Endorsement details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/endorsements [post]
func (h *EndorsementHandler) CreateEndorsement(ctx *gin.Context) {
	var req policySchema.CreateEndorsementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Use policy ID from URL path
	req.PolicyID = ctx.Param("id")

	resp := h.endorsementSvc.CreateEndorsement(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetEndorsement godoc
// @Summary      Get an endorsement
// @Description  Get endorsement details by ID
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Endorsement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/endorsements/{id} [get]
func (h *EndorsementHandler) GetEndorsement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid endorsement ID")
		return
	}

	resp := h.endorsementSvc.GetEndorsement(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListEndorsements godoc
// @Summary      List endorsements for a policy
// @Description  List all endorsements associated with a policy
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/endorsements [get]
func (h *EndorsementHandler) ListEndorsements(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.endorsementSvc.ListByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApproveEndorsement godoc
// @Summary      Approve an endorsement
// @Description  Approve a pending endorsement by ID
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Endorsement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/endorsements/{id}/approve [put]
func (h *EndorsementHandler) ApproveEndorsement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid endorsement ID")
		return
	}

	resp := h.endorsementSvc.ApproveEndorsement(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RejectEndorsement godoc
// @Summary      Reject an endorsement
// @Description  Reject a pending endorsement by ID with a reason
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Endorsement ID"
// @Param        request body policySchema.RejectEndorsementRequest true "Rejection reason"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/endorsements/{id}/reject [put]
func (h *EndorsementHandler) RejectEndorsement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid endorsement ID")
		return
	}

	var req policySchema.RejectEndorsementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.endorsementSvc.RejectEndorsement(ctx.Request.Context(), id, req.Reason, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApplyEndorsement godoc
// @Summary      Apply an endorsement
// @Description  Apply an approved endorsement to the policy
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Endorsement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/endorsements/{id}/apply [put]
func (h *EndorsementHandler) ApplyEndorsement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid endorsement ID")
		return
	}

	resp := h.endorsementSvc.ApplyEndorsement(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CancelEndorsement godoc
// @Summary      Cancel an endorsement
// @Description  Cancel a pending or approved endorsement
// @Tags         Endorsements
// @Accept       json
// @Produce      json
// @Param        id path string true "Endorsement ID"
// @Param        request body policySchema.CancelEndorsementRequest true "Cancellation reason"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/endorsements/{id}/cancel [put]
func (h *EndorsementHandler) CancelEndorsement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid endorsement ID")
		return
	}

	var req policySchema.CancelEndorsementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.endorsementSvc.CancelEndorsement(ctx.Request.Context(), id, req.Reason, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

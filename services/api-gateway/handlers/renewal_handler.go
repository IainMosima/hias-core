package handlers

import (
	"net/http"

	policySchema "github.com/bitbiz/hias-core/domains/policy/schema"
	"github.com/bitbiz/hias-core/domains/policy/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RenewalHandler struct {
	renewalSvc service.RenewalService
}

func NewRenewalHandler(renewalSvc service.RenewalService) *RenewalHandler {
	return &RenewalHandler{renewalSvc: renewalSvc}
}

// InitiateRenewal godoc
// @Summary      Initiate a policy renewal
// @Description  Initiate a renewal process for a policy
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        request body policySchema.InitiateRenewalRequest true "Renewal details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/renewals [post]
func (h *RenewalHandler) InitiateRenewal(ctx *gin.Context) {
	var req policySchema.InitiateRenewalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Use policy ID from URL path
	req.PolicyID = ctx.Param("id")

	resp := h.renewalSvc.InitiateRenewal(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetRenewal godoc
// @Summary      Get a renewal
// @Description  Get renewal details by ID
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        id path string true "Renewal ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/renewals/{id} [get]
func (h *RenewalHandler) GetRenewal(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid renewal ID")
		return
	}

	resp := h.renewalSvc.GetRenewal(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRenewals godoc
// @Summary      List renewals for a policy
// @Description  List all renewals associated with a policy
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/renewals [get]
func (h *RenewalHandler) ListRenewals(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	resp := h.renewalSvc.ListByPolicy(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApproveRenewal godoc
// @Summary      Approve a renewal
// @Description  Approve a pending renewal by ID
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        id path string true "Renewal ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/renewals/{id}/approve [put]
func (h *RenewalHandler) ApproveRenewal(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid renewal ID")
		return
	}

	resp := h.renewalSvc.ApproveRenewal(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RejectRenewal godoc
// @Summary      Reject a renewal
// @Description  Reject a pending renewal by ID with a reason
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        id path string true "Renewal ID"
// @Param        request body policySchema.RejectRenewalRequest true "Rejection reason"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/renewals/{id}/reject [put]
func (h *RenewalHandler) RejectRenewal(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid renewal ID")
		return
	}

	var req policySchema.RejectRenewalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.renewalSvc.RejectRenewal(ctx.Request.Context(), id, req.Reason)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CompleteRenewal godoc
// @Summary      Complete a renewal
// @Description  Complete an approved renewal process
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        id path string true "Renewal ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/renewals/{id}/complete [post]
func (h *RenewalHandler) CompleteRenewal(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid renewal ID")
		return
	}

	resp := h.renewalSvc.CompleteRenewal(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ExpireRenewals godoc
// @Summary      Expire pending renewals
// @Description  Expire all pending renewals that have passed their deadline
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/renewals/expire [post]
func (h *RenewalHandler) ExpireRenewals(ctx *gin.Context) {
	resp := h.renewalSvc.ExpirePendingRenewals(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkInitiateRenewals godoc
// @Summary      Bulk initiate renewals
// @Description  Initiate renewal processes for multiple policies at once
// @Tags         Renewals
// @Accept       json
// @Produce      json
// @Param        request body policySchema.BulkRenewalRequest true "List of policy IDs"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/renewals/bulk [post]
func (h *RenewalHandler) BulkInitiateRenewals(ctx *gin.Context) {
	var req policySchema.BulkRenewalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var ids []uuid.UUID
	for _, idStr := range req.PolicyIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID: "+idStr)
			return
		}
		ids = append(ids, id)
	}

	resp := h.renewalSvc.BulkInitiateRenewals(ctx.Request.Context(), ids, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

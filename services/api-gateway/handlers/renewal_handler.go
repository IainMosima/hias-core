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

func (h *RenewalHandler) ExpireRenewals(ctx *gin.Context) {
	resp := h.renewalSvc.ExpirePendingRenewals(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

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

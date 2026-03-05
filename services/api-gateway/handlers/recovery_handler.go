package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecoveryHandler struct {
	recoverySvc service.RecoveryService
}

func NewRecoveryHandler(recoverySvc service.RecoveryService) *RecoveryHandler {
	return &RecoveryHandler{recoverySvc: recoverySvc}
}

func (h *RecoveryHandler) CreateRecovery(ctx *gin.Context) {
	var req reinsuranceSchema.CreateRecoveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.CreateRecovery(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) GetRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	resp := h.recoverySvc.GetRecovery(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) ListRecoveries(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	claimIDStr := ctx.Query("claim")
	treatyIDStr := ctx.Query("treaty")

	if claimIDStr != "" {
		claimID, err := uuid.Parse(claimIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
			return
		}
		resp := h.recoverySvc.ListByClaim(ctx.Request.Context(), claimID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.recoverySvc.GetRecoveryCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if treatyIDStr != "" {
		treatyID, err := uuid.Parse(treatyIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
			return
		}
		resp := h.recoverySvc.ListByTreaty(ctx.Request.Context(), treatyID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.recoverySvc.GetRecoveryCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	utils.RespondError(ctx, http.StatusBadRequest, "Query parameter 'claim' or 'treaty' is required")
}

func (h *RecoveryHandler) ListOutstanding(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.recoverySvc.ListOutstanding(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.recoverySvc.GetRecoveryCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

func (h *RecoveryHandler) ApplyRecoveryForClaim(ctx *gin.Context) {
	claimID, err := uuid.Parse(ctx.Param("claimId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req reinsuranceSchema.ApplyRecoveryForClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.ApplyRecoveryForClaim(ctx.Request.Context(), claimID, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) AcknowledgeRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.AcknowledgeRecovery(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) RequestInfo(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.RequestInfo(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) ApproveRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.ApproveRecovery(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) RecordPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecordPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.RecordPayment(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) WriteOffRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.WriteOffRecovery(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) AgedAnalysis(ctx *gin.Context) {
	resp := h.recoverySvc.GetAgedAnalysis(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RecoveryHandler) GetWorkflowEvents(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	resp := h.recoverySvc.GetWorkflowEvents(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

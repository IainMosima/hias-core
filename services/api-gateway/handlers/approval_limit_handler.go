package handlers

import (
	"net/http"

	salesSchema "github.com/bitbiz/hias-core/domains/sales/schema"
	"github.com/bitbiz/hias-core/domains/sales/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalLimitHandler struct {
	approvalSvc service.ApprovalLimitService
}

func NewApprovalLimitHandler(approvalSvc service.ApprovalLimitService) *ApprovalLimitHandler {
	return &ApprovalLimitHandler{approvalSvc: approvalSvc}
}

func (h *ApprovalLimitHandler) ListLimits(ctx *gin.Context) {
	resp := h.approvalSvc.GetLimits(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ApprovalLimitHandler) CreateLimit(ctx *gin.Context) {
	var req salesSchema.CreateApprovalLimitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.approvalSvc.CreateLimit(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ApprovalLimitHandler) UpdateLimit(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid approval limit ID")
		return
	}

	var req salesSchema.UpdateApprovalLimitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.approvalSvc.UpdateLimit(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

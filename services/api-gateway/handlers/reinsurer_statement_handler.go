package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReinsurerStatementHandler struct {
	statementSvc service.ReinsurerStatementService
}

func NewReinsurerStatementHandler(statementSvc service.ReinsurerStatementService) *ReinsurerStatementHandler {
	return &ReinsurerStatementHandler{statementSvc: statementSvc}
}

func (h *ReinsurerStatementHandler) GenerateStatement(ctx *gin.Context) {
	var req reinsuranceSchema.GenerateStatementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.statementSvc.GenerateStatement(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReinsurerStatementHandler) GetStatement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid statement ID")
		return
	}

	resp := h.statementSvc.GetStatement(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReinsurerStatementHandler) ListByTreaty(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	treatyIDStr := ctx.Query("treaty")
	if treatyIDStr == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "treaty query parameter is required")
		return
	}

	treatyID, err := uuid.Parse(treatyIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.statementSvc.ListByTreaty(ctx.Request.Context(), treatyID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReinsurerStatementHandler) IssueStatement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid statement ID")
		return
	}

	resp := h.statementSvc.IssueStatement(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReinsurerStatementHandler) AcknowledgeStatement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid statement ID")
		return
	}

	resp := h.statementSvc.AcknowledgeStatement(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReinsurerStatementHandler) SettleStatement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid statement ID")
		return
	}

	resp := h.statementSvc.SettleStatement(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReinsurerStatementHandler) CalculateProfitCommission(ctx *gin.Context) {
	var req reinsuranceSchema.CalculateProfitCommissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.statementSvc.CalculateProfitCommission(ctx.Request.Context(), req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

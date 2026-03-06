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

// GenerateStatement godoc
// @Summary      Generate a reinsurer statement
// @Description  Generate a new reinsurer statement for a treaty period
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.GenerateStatementRequest true "Generate statement request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reinsurer-statements [post]
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

// GetStatement godoc
// @Summary      Get a reinsurer statement
// @Description  Retrieve a single reinsurer statement by ID
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reinsurer-statements/{id} [get]
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

// ListByTreaty godoc
// @Summary      List statements by treaty
// @Description  Retrieve a paginated list of reinsurer statements for a specific treaty
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        treaty query string true "Treaty ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/statements [get]
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

// IssueStatement godoc
// @Summary      Issue a reinsurer statement
// @Description  Issue a draft reinsurer statement by ID
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reinsurer-statements/{id}/issue [put]
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

// AcknowledgeStatement godoc
// @Summary      Acknowledge a reinsurer statement
// @Description  Acknowledge receipt of a reinsurer statement by ID
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reinsurer-statements/{id}/acknowledge [put]
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

// SettleStatement godoc
// @Summary      Settle a reinsurer statement
// @Description  Mark a reinsurer statement as settled by ID
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reinsurer-statements/{id}/settle [put]
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

// CalculateProfitCommission godoc
// @Summary      Calculate profit commission
// @Description  Calculate profit commission for a treaty period
// @Tags         ReinsurerStatements
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.CalculateProfitCommissionRequest true "Calculate profit commission request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reinsurer-statements/profit-commission [post]
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

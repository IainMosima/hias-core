package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StatementHandler struct {
	statementSvc service.StatementService
}

func NewStatementHandler(statementSvc service.StatementService) *StatementHandler {
	return &StatementHandler{statementSvc: statementSvc}
}

// UploadStatement godoc
// @Summary      Upload a provider statement
// @Description  Upload a new provider statement for reconciliation
// @Tags         ProviderStatements
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.UploadStatementRequest true "Statement details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/statements [post]
func (h *StatementHandler) UploadStatement(ctx *gin.Context) {
	var req billingSchema.UploadStatementRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.statementSvc.UploadStatement(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetStatement godoc
// @Summary      Get a provider statement
// @Description  Get provider statement details by ID
// @Tags         ProviderStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/provider-statements/{id} [get]
func (h *StatementHandler) GetStatement(ctx *gin.Context) {
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

// ListByProvider godoc
// @Summary      List statements by provider
// @Description  List all statements associated with a provider
// @Tags         ProviderStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/statements [get]
func (h *StatementHandler) ListByProvider(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)
	resp := h.statementSvc.ListByProvider(ctx.Request.Context(), providerID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListLineItems godoc
// @Summary      List statement line items
// @Description  List all line items for a provider statement
// @Tags         ProviderStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/provider-statements/{id}/line-items [get]
func (h *StatementHandler) ListLineItems(ctx *gin.Context) {
	statementID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid statement ID")
		return
	}

	resp := h.statementSvc.ListLineItems(ctx.Request.Context(), statementID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ReconcileStatement godoc
// @Summary      Reconcile a provider statement
// @Description  Reconcile a provider statement against claims
// @Tags         ProviderStatements
// @Accept       json
// @Produce      json
// @Param        id path string true "Statement ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/provider-statements/{id}/reconcile [post]
func (h *StatementHandler) ReconcileStatement(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid statement ID")
		return
	}

	resp := h.statementSvc.ReconcileStatement(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

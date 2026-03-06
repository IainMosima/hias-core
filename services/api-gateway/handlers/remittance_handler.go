package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RemittanceHandler struct {
	remittanceSvc service.RemittanceService
}

func NewRemittanceHandler(remittanceSvc service.RemittanceService) *RemittanceHandler {
	return &RemittanceHandler{remittanceSvc: remittanceSvc}
}

// CreateRemittance godoc
// @Summary      Create a remittance
// @Description  Create a new remittance for a provider
// @Tags         Remittances
// @Accept       json
// @Produce      json
// @Param        request body object true "Remittance creation payload with provider_id"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/remittances [post]
func (h *RemittanceHandler) CreateRemittance(ctx *gin.Context) {
	var req struct {
		ProviderID string `json:"provider_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	providerID, _ := uuid.Parse(req.ProviderID)
	resp := h.remittanceSvc.CreateRemittance(ctx.Request.Context(), providerID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetRemittance godoc
// @Summary      Get a remittance by ID
// @Description  Retrieve a single remittance by its unique identifier
// @Tags         Remittances
// @Accept       json
// @Produce      json
// @Param        id path string true "Remittance ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/remittances/{id} [get]
func (h *RemittanceHandler) GetRemittance(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid remittance ID")
		return
	}

	resp := h.remittanceSvc.GetRemittance(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRemittances godoc
// @Summary      List all remittances
// @Description  Retrieve a paginated list of remittances
// @Tags         Remittances
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/remittances [get]
func (h *RemittanceHandler) ListRemittances(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.remittanceSvc.ListRemittances(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RemittanceHandler) SendAdvice(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid remittance ID")
		return
	}

	resp := h.remittanceSvc.SendRemittanceAdvice(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ExportPaymentFile godoc
// @Summary      Export payment file for a remittance
// @Description  Export the payment file associated with the specified remittance
// @Tags         Remittances
// @Accept       json
// @Produce      json
// @Param        id path string true "Remittance ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/remittances/{id}/export [get]
func (h *RemittanceHandler) ExportPaymentFile(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid remittance ID")
		return
	}

	resp := h.remittanceSvc.ExportPaymentFile(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RemittanceHandler struct {
	remittanceService service.RemittanceService
}

func NewRemittanceHandler(remittanceService service.RemittanceService) *RemittanceHandler {
	return &RemittanceHandler{remittanceService: remittanceService}
}

func (h *RemittanceHandler) CreateRemittance(ctx *gin.Context) {
	var req struct {
		ProviderID string `json:"provider_id" binding:"required,uuid"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	providerID, _ := uuid.Parse(req.ProviderID)

	resp := h.remittanceService.CreateRemittance(ctx.Request.Context(), providerID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RemittanceHandler) GetRemittance(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid remittance ID")
		return
	}

	resp := h.remittanceService.GetRemittance(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *RemittanceHandler) ListRemittances(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.remittanceService.ListRemittances(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}

	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, 0)
}

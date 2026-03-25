package handlers

import (
	"net/http"

	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/domains/provider/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RateCardHandler struct {
	rateCardSvc service.RateCardService
}

func NewRateCardHandler(rateCardSvc service.RateCardService) *RateCardHandler {
	return &RateCardHandler{rateCardSvc: rateCardSvc}
}

// CreateRateCard godoc
// @Summary      Create a rate card
// @Description  Create a new rate card for a specific healthcare provider
// @Tags         RateCards
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        request body providerSchema.CreateRateCardRequest true "Rate card creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/rate-cards [post]
func (h *RateCardHandler) CreateRateCard(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var req providerSchema.CreateRateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.rateCardSvc.CreateRateCard(ctx.Request.Context(), providerID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRateCards godoc
// @Summary      List rate cards
// @Description  Retrieve all rate cards for a specific healthcare provider
// @Tags         RateCards
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/rate-cards [get]
func (h *RateCardHandler) ListRateCards(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	resp := h.rateCardSvc.ListRateCards(ctx.Request.Context(), providerID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// BulkCreateRateCards godoc
// @Summary      Bulk create rate cards
// @Description  Create multiple rate cards at once for a specific healthcare provider
// @Tags         RateCards
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        request body providerSchema.BulkCreateRateCardRequest true "Bulk rate card creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/rate-cards/bulk [post]
func (h *RateCardHandler) BulkCreateRateCards(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var req providerSchema.BulkCreateRateCardRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.rateCardSvc.BulkCreateRateCards(ctx.Request.Context(), providerID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

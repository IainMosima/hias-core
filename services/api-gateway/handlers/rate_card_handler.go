package handlers

import (
	"fmt"
	"net/http"

	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/bitbiz/hias-core/domains/provider/repository"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RateCardHandler struct {
	rateCardRepo repository.RateCardRepository
}

func NewRateCardHandler(rateCardRepo repository.RateCardRepository) *RateCardHandler {
	return &RateCardHandler{rateCardRepo: rateCardRepo}
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

	rateCard := &entity.RateCard{
		ProviderID:    providerID,
		ProcedureCode: req.ProcedureCode,
		ProcedureName: req.ProcedureName,
		RateAmount:    req.RateAmount,
		EffectiveDate: req.EffectiveDate,
	}

	created, err := h.rateCardRepo.Create(ctx.Request.Context(), rateCard)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to create rate card")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, "Rate card created", providerSchema.ToRateCardResponse(created))
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

	rateCards, err := h.rateCardRepo.ListByProvider(ctx.Request.Context(), providerID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to list rate cards")
		return
	}

	responses := make([]providerSchema.RateCardResponse, len(rateCards))
	for i, r := range rateCards {
		responses[i] = providerSchema.ToRateCardResponse(r)
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Rate cards retrieved", responses)
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

	var responses []providerSchema.RateCardResponse
	for _, rc := range req.RateCards {
		rateCard := &entity.RateCard{
			ProviderID:    providerID,
			ProcedureCode: rc.ProcedureCode,
			ProcedureName: rc.ProcedureName,
			RateAmount:    rc.RateAmount,
			EffectiveDate: rc.EffectiveDate,
			AgeFrom:       rc.AgeFrom,
			AgeTo:         rc.AgeTo,
			Gender:        rc.Gender,
			Relationship:  rc.Relationship,
		}
		created, createErr := h.rateCardRepo.Create(ctx.Request.Context(), rateCard)
		if createErr != nil {
			utils.RespondError(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to create rate card for %s: %v", rc.ProcedureCode, createErr))
			return
		}
		responses = append(responses, providerSchema.ToRateCardResponse(created))
	}

	utils.RespondSuccess(ctx, http.StatusCreated, fmt.Sprintf("%d rate cards created", len(responses)), responses)
}

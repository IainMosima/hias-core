package handlers

import (
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

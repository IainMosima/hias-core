package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/bitbiz/hias-core/domains/provider/repository"
	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContractHandler struct {
	contractRepo repository.ContractRepository
}

func NewContractHandler(contractRepo repository.ContractRepository) *ContractHandler {
	return &ContractHandler{contractRepo: contractRepo}
}

func (h *ContractHandler) CreateContract(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var req providerSchema.CreateContractRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	contract := &entity.Contract{
		ProviderID: providerID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Terms:      req.Terms,
		Status:     string(shared.ContractStatusActive),
	}

	created, err := h.contractRepo.Create(ctx.Request.Context(), contract)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to create contract")
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, "Contract created", providerSchema.ToContractResponse(created))
}

func (h *ContractHandler) ListContracts(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	contracts, err := h.contractRepo.ListByProvider(ctx.Request.Context(), providerID)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Failed to list contracts")
		return
	}

	responses := make([]providerSchema.ContractResponse, len(contracts))
	for i, c := range contracts {
		responses[i] = providerSchema.ToContractResponse(c)
	}

	utils.RespondSuccess(ctx, http.StatusOK, "Contracts retrieved", responses)
}

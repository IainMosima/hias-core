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

// CreateContract godoc
// @Summary      Create a provider contract
// @Description  Create a new contract for a specific healthcare provider
// @Tags         Contracts
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        request body providerSchema.CreateContractRequest true "Contract creation request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/contracts [post]
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

// ListContracts godoc
// @Summary      List provider contracts
// @Description  Retrieve all contracts for a specific healthcare provider
// @Tags         Contracts
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/contracts [get]
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

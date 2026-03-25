package handlers

import (
	"net/http"

	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/domains/provider/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContractHandler struct {
	contractSvc service.ContractService
}

func NewContractHandler(contractSvc service.ContractService) *ContractHandler {
	return &ContractHandler{contractSvc: contractSvc}
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

	resp := h.contractSvc.CreateContract(ctx.Request.Context(), providerID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
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

	resp := h.contractSvc.ListContracts(ctx.Request.Context(), providerID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProviderNetworkHandler struct {
	networkSvc service.ProviderNetworkService
}

func NewProviderNetworkHandler(networkSvc service.ProviderNetworkService) *ProviderNetworkHandler {
	return &ProviderNetworkHandler{networkSvc: networkSvc}
}

// CreateProviderNetwork godoc
// @Summary      Create a provider network for a plan
// @Description  Create a new provider network associated with the specified plan
// @Tags         ProviderNetworks
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body schema.CreateProviderNetworkRequest true "Provider network creation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/provider-networks [post]
func (h *ProviderNetworkHandler) CreateProviderNetwork(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req schema.CreateProviderNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.networkSvc.CreateProviderNetwork(ctx.Request.Context(), planID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListProviderNetworks godoc
// @Summary      List provider networks for a plan
// @Description  Retrieve all provider networks associated with the specified plan
// @Tags         ProviderNetworks
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/provider-networks [get]
func (h *ProviderNetworkHandler) ListProviderNetworks(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	resp := h.networkSvc.ListProviderNetworksByPlan(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateProviderNetworkStatus godoc
// @Summary      Update provider network status
// @Description  Update the status of an existing provider network by its unique identifier
// @Tags         ProviderNetworks
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider Network ID"
// @Param        request body schema.UpdateProviderNetworkStatusRequest true "Status update payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/provider-networks/{id}/status [put]
func (h *ProviderNetworkHandler) UpdateProviderNetworkStatus(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider network ID")
		return
	}

	var req schema.UpdateProviderNetworkStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.networkSvc.UpdateProviderNetworkStatus(ctx.Request.Context(), id, req.Status)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteProviderNetwork godoc
// @Summary      Delete a provider network
// @Description  Delete a provider network by its unique identifier
// @Tags         ProviderNetworks
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider Network ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/provider-networks/{id} [delete]
func (h *ProviderNetworkHandler) DeleteProviderNetwork(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider network ID")
		return
	}

	resp := h.networkSvc.DeleteProviderNetwork(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

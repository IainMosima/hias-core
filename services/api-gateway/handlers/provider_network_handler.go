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

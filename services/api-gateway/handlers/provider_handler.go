package handlers

import (
	"net/http"

	providerSchema "github.com/bitbiz/hias-core/domains/provider/schema"
	"github.com/bitbiz/hias-core/domains/provider/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProviderHandler struct {
	providerSvc service.ProviderService
}

func NewProviderHandler(providerSvc service.ProviderService) *ProviderHandler {
	return &ProviderHandler{providerSvc: providerSvc}
}

// RegisterProvider godoc
// @Summary      Register a new provider
// @Description  Registers a new healthcare provider in the system
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        request body providerSchema.RegisterProviderRequest true "Register provider request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers [post]
func (h *ProviderHandler) RegisterProvider(ctx *gin.Context) {
	var req providerSchema.RegisterProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.providerSvc.RegisterProvider(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetProvider godoc
// @Summary      Get a provider by ID
// @Description  Retrieves the details of a specific healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id} [get]
func (h *ProviderHandler) GetProvider(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	resp := h.providerSvc.GetProvider(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListProviders godoc
// @Summary      List all providers
// @Description  Returns a paginated list of all healthcare providers with optional search
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        search query string false "Search by name, license number, or email"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers [get]
func (h *ProviderHandler) ListProviders(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	search := ctx.Query("search")

	if search != "" {
		resp := h.providerSvc.ListProvidersFiltered(ctx.Request.Context(), search, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.providerSvc.CountProvidersFiltered(ctx.Request.Context(), search)
		if countResp.Error != nil {
			utils.RespondError(ctx, countResp.StatusCode, countResp.Message)
			return
		}
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	resp := h.providerSvc.ListProviders(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.providerSvc.GetTotalCount(ctx.Request.Context())
	if countResp.Error != nil {
		utils.RespondError(ctx, countResp.StatusCode, countResp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// UpdateProvider godoc
// @Summary      Update a provider
// @Description  Updates the details of an existing healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        request body providerSchema.UpdateProviderRequest true "Update provider request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id} [put]
func (h *ProviderHandler) UpdateProvider(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var req providerSchema.UpdateProviderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.providerSvc.UpdateProvider(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CredentialProvider godoc
// @Summary      Credential a provider
// @Description  Marks the specified provider as credentialed
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/credential [put]
func (h *ProviderHandler) CredentialProvider(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	resp := h.providerSvc.CredentialProvider(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ActivateProvider godoc
// @Summary      Activate a provider
// @Description  Activates the specified healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/activate [put]
func (h *ProviderHandler) ActivateProvider(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	resp := h.providerSvc.ActivateProvider(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// SuspendProvider godoc
// @Summary      Suspend a provider
// @Description  Suspends the specified healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/suspend [put]
func (h *ProviderHandler) SuspendProvider(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	resp := h.providerSvc.SuspendProvider(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// TerminateProvider godoc
// @Summary      Terminate a provider
// @Description  Terminates the specified healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/terminate [put]
func (h *ProviderHandler) TerminateProvider(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	resp := h.providerSvc.TerminateProvider(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateTier godoc
// @Summary      Update provider tier
// @Description  Updates the tier classification of a healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        request body object true "Tier update request with tier field"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/tier [put]
func (h *ProviderHandler) UpdateTier(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var req struct {
		Tier string `json:"tier" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.providerSvc.UpdateTier(ctx.Request.Context(), id, req.Tier, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByTier godoc
// @Summary      List providers by tier
// @Description  Returns a paginated list of providers filtered by tier
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        tier query string true "Provider tier"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/by-tier [get]
func (h *ProviderHandler) ListByTier(ctx *gin.Context) {
	tier := ctx.Query("tier")
	if tier == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "tier query parameter required")
		return
	}

	pagination := utils.GetPaginationParams(ctx)
	resp := h.providerSvc.ListByTier(ctx.Request.Context(), tier, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateAccreditation godoc
// @Summary      Update provider accreditation
// @Description  Updates the accreditation details of a healthcare provider
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        request body providerSchema.UpdateAccreditationRequest true "Update accreditation request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/accreditation [put]
func (h *ProviderHandler) UpdateAccreditation(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	var req providerSchema.UpdateAccreditationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.providerSvc.UpdateAccreditation(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByAccreditationStatus godoc
// @Summary      List providers by accreditation status
// @Description  Returns a paginated list of providers filtered by accreditation status
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        status query string true "Accreditation status"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/by-accreditation [get]
func (h *ProviderHandler) ListByAccreditationStatus(ctx *gin.Context) {
	status := ctx.Query("status")
	if status == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "status query parameter required")
		return
	}

	pagination := utils.GetPaginationParams(ctx)
	resp := h.providerSvc.ListByAccreditationStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListExpiringAccreditations godoc
// @Summary      List providers with expiring accreditations
// @Description  Returns a paginated list of providers whose accreditations are expiring within 30 days
// @Tags         Providers
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/expiring-accreditations [get]
func (h *ProviderHandler) ListExpiringAccreditations(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	resp := h.providerSvc.ListExpiringAccreditations(ctx.Request.Context(), 30, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

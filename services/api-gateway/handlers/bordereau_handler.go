package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BordereauHandler struct {
	bordereauSvc service.BordereauService
}

func NewBordereauHandler(bordereauSvc service.BordereauService) *BordereauHandler {
	return &BordereauHandler{bordereauSvc: bordereauSvc}
}

// GeneratePremiumBordereau godoc
// @Summary      Generate a premium bordereau
// @Description  Generate a premium bordereau report for a treaty
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.GenerateBordereauRequest true "Generate bordereau request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/bordereaux/premium [post]
func (h *BordereauHandler) GeneratePremiumBordereau(ctx *gin.Context) {
	var req reinsuranceSchema.GenerateBordereauRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.bordereauSvc.GeneratePremiumBordereau(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GenerateClaimBordereau godoc
// @Summary      Generate a claim bordereau
// @Description  Generate a claim bordereau report for a treaty
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.GenerateBordereauRequest true "Generate bordereau request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/bordereaux/claim [post]
func (h *BordereauHandler) GenerateClaimBordereau(ctx *gin.Context) {
	var req reinsuranceSchema.GenerateBordereauRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.bordereauSvc.GenerateClaimBordereau(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetBordereau godoc
// @Summary      Get a bordereau
// @Description  Retrieve a single bordereau by ID
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        id path string true "Bordereau ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/bordereaux/{id} [get]
func (h *BordereauHandler) GetBordereau(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid bordereau ID")
		return
	}

	resp := h.bordereauSvc.GetBordereau(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListByTreaty godoc
// @Summary      List bordereaux by treaty
// @Description  Retrieve a paginated list of bordereaux for a specific treaty
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        treaty query string true "Treaty ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/bordereaux [get]
func (h *BordereauHandler) ListByTreaty(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	treatyIDStr := ctx.Query("treaty")
	if treatyIDStr == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "Query parameter 'treaty' is required")
		return
	}

	treatyID, err := uuid.Parse(treatyIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.bordereauSvc.ListByTreaty(ctx.Request.Context(), treatyID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, 0)
}

// FinalizeBordereau godoc
// @Summary      Finalize a bordereau
// @Description  Finalize a draft bordereau by ID
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        id path string true "Bordereau ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/bordereaux/{id}/finalize [put]
func (h *BordereauHandler) FinalizeBordereau(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid bordereau ID")
		return
	}

	resp := h.bordereauSvc.FinalizeBordereau(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// MarkSent godoc
// @Summary      Mark a bordereau as sent
// @Description  Mark a finalized bordereau as sent to the reinsurer
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        id path string true "Bordereau ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/bordereaux/{id}/mark-sent [put]
func (h *BordereauHandler) MarkSent(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid bordereau ID")
		return
	}

	resp := h.bordereauSvc.MarkSent(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListItems godoc
// @Summary      List bordereau items
// @Description  Retrieve all line items for a specific bordereau
// @Tags         Bordereaux
// @Accept       json
// @Produce      json
// @Param        id path string true "Bordereau ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/bordereaux/{id}/items [get]
func (h *BordereauHandler) ListItems(ctx *gin.Context) {
	bordereauID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid bordereau ID")
		return
	}

	resp := h.bordereauSvc.ListItems(ctx.Request.Context(), bordereauID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

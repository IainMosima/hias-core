package handlers

import (
	"net/http"

	"github.com/bitbiz/hias-core/domains/product/schema"
	"github.com/bitbiz/hias-core/domains/product/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BenefitHandler struct {
	benefitSvc service.BenefitService
}

func NewBenefitHandler(benefitSvc service.BenefitService) *BenefitHandler {
	return &BenefitHandler{benefitSvc: benefitSvc}
}

// CreateBenefit godoc
// @Summary      Create a benefit for a plan
// @Description  Create a new benefit associated with the specified plan
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Param        request body schema.CreateBenefitRequest true "Benefit creation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/benefits [post]
func (h *BenefitHandler) CreateBenefit(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	var req schema.CreateBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.benefitSvc.CreateBenefit(ctx.Request.Context(), planID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListBenefits godoc
// @Summary      List benefits for a plan
// @Description  Retrieve all benefits associated with the specified plan
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Plan ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/plans/{id}/benefits [get]
func (h *BenefitHandler) ListBenefits(ctx *gin.Context) {
	planID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid plan ID")
		return
	}

	resp := h.benefitSvc.ListBenefitsByPlan(ctx.Request.Context(), planID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetBenefit godoc
// @Summary      Get a benefit by ID
// @Description  Retrieve a single benefit by its unique identifier
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Benefit ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/benefits/{id} [get]
func (h *BenefitHandler) GetBenefit(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid benefit ID")
		return
	}

	resp := h.benefitSvc.GetBenefit(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateBenefit godoc
// @Summary      Update a benefit
// @Description  Update an existing benefit by its unique identifier
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Benefit ID"
// @Param        request body schema.UpdateBenefitRequest true "Benefit update payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/benefits/{id} [put]
func (h *BenefitHandler) UpdateBenefit(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid benefit ID")
		return
	}

	var req schema.UpdateBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.benefitSvc.UpdateBenefit(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// DeleteBenefit godoc
// @Summary      Delete a benefit
// @Description  Delete a benefit by its unique identifier
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Benefit ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/benefits/{id} [delete]
func (h *BenefitHandler) DeleteBenefit(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid benefit ID")
		return
	}

	resp := h.benefitSvc.DeleteBenefit(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// CreateSubBenefit godoc
// @Summary      Create a sub-benefit
// @Description  Create a new sub-benefit under the specified parent benefit
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Parent Benefit ID"
// @Param        request body schema.CreateBenefitRequest true "Sub-benefit creation payload"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/benefits/{id}/sub-benefits [post]
func (h *BenefitHandler) CreateSubBenefit(ctx *gin.Context) {
	parentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid benefit ID")
		return
	}

	var req schema.CreateBenefitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.benefitSvc.CreateSubBenefit(ctx.Request.Context(), parentID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListSubBenefits godoc
// @Summary      List sub-benefits
// @Description  Retrieve all sub-benefits under the specified parent benefit
// @Tags         Benefits
// @Accept       json
// @Produce      json
// @Param        id path string true "Parent Benefit ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/benefits/{id}/sub-benefits [get]
func (h *BenefitHandler) ListSubBenefits(ctx *gin.Context) {
	parentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid benefit ID")
		return
	}

	resp := h.benefitSvc.ListSubBenefits(ctx.Request.Context(), parentID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

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

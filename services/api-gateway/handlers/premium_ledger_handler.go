package handlers

import (
	"net/http"

	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PremiumLedgerHandler struct {
	svc service.PremiumLedgerService
}

func NewPremiumLedgerHandler(svc service.PremiumLedgerService) *PremiumLedgerHandler {
	return &PremiumLedgerHandler{svc: svc}
}

// CreateEntry godoc
// @Summary      Create a premium ledger entry
// @Description  Record a new premium ledger entry for a policy
// @Tags         PremiumLedger
// @Accept       json
// @Produce      json
// @Param        request body billingSchema.CreatePremiumLedgerRequest true "Create premium ledger entry request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/premium-ledger [post]
func (h *PremiumLedgerHandler) CreateEntry(ctx *gin.Context) {
	var req billingSchema.CreatePremiumLedgerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}
	resp := h.svc.RecordEntry(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetRegister godoc
// @Summary      Get premium register for a policy
// @Description  Retrieve the paginated premium register for a specific policy
// @Tags         PremiumLedger
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/premium-register [get]
func (h *PremiumLedgerHandler) GetRegister(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	pagination := utils.GetPaginationParams(ctx)
	resp := h.svc.GetRegister(ctx.Request.Context(), policyID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetBalance godoc
// @Summary      Get premium balance for a policy
// @Description  Retrieve the current premium balance for a specific policy
// @Tags         PremiumLedger
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/premium-balance [get]
func (h *PremiumLedgerHandler) GetBalance(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}
	resp := h.svc.GetBalance(ctx.Request.Context(), policyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

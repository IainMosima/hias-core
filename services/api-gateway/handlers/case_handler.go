package handlers

import (
	"net/http"

	claimsSchema "github.com/bitbiz/hias-core/domains/claims/schema"
	"github.com/bitbiz/hias-core/domains/claims/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CaseHandler struct {
	caseSvc service.CaseService
}

func NewCaseHandler(caseSvc service.CaseService) *CaseHandler {
	return &CaseHandler{caseSvc: caseSvc}
}

func (h *CaseHandler) CreateCase(ctx *gin.Context) {
	var req claimsSchema.CreateCaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.caseSvc.CreateCase(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) GetCase(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid case ID")
		return
	}

	resp := h.caseSvc.GetCase(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) ListCases(ctx *gin.Context) {
	status := ctx.Query("status")
	if status == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "status query parameter required")
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	resp := h.caseSvc.ListByStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, int64(len(resp.Data)))
}

func (h *CaseHandler) ListByPolicy(ctx *gin.Context) {
	policyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid policy ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	resp := h.caseSvc.ListByPolicy(ctx.Request.Context(), policyID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, int64(len(resp.Data)))
}

func (h *CaseHandler) ListByMember(ctx *gin.Context) {
	memberID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid member ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	resp := h.caseSvc.ListByMember(ctx.Request.Context(), memberID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, int64(len(resp.Data)))
}

func (h *CaseHandler) ListByProvider(ctx *gin.Context) {
	providerID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid provider ID")
		return
	}

	pagination := utils.GetPaginationParams(ctx)

	resp := h.caseSvc.ListByProvider(ctx.Request.Context(), providerID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, int64(len(resp.Data)))
}

func (h *CaseHandler) AdmitCase(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid case ID")
		return
	}

	var req claimsSchema.AdmitCaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.caseSvc.AdmitCase(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) UpdateCase(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid case ID")
		return
	}

	var req claimsSchema.UpdateCaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.caseSvc.UpdateCase(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) StartTreatment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid case ID")
		return
	}

	resp := h.caseSvc.StartTreatment(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) DischargeCase(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid case ID")
		return
	}

	var req claimsSchema.DischargeCaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.caseSvc.DischargeCase(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) CloseCase(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid case ID")
		return
	}

	resp := h.caseSvc.CloseCase(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *CaseHandler) CountByStatus(ctx *gin.Context) {
	status := ctx.Query("status")
	if status == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "status query parameter required")
		return
	}

	resp := h.caseSvc.CountByStatus(ctx.Request.Context(), status)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

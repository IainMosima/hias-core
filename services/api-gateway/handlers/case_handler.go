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

// CreateCase godoc
// @Summary      Create a new case
// @Description  Create a new case for claims management
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        request body claimsSchema.CreateCaseRequest true "Case details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases [post]
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

// GetCase godoc
// @Summary      Get a case
// @Description  Get case details by ID
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Case ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/{id} [get]
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

// ListCases godoc
// @Summary      List cases
// @Description  List cases filtered by status with pagination
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        status query string true "Case status"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases [get]
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

// ListByPolicy godoc
// @Summary      List cases by policy
// @Description  List all cases associated with a policy
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Policy ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/policies/{id}/cases [get]
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

// ListByMember godoc
// @Summary      List cases by member
// @Description  List all cases associated with a member
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Member ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/members/{id}/cases [get]
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

// ListByProvider godoc
// @Summary      List cases by provider
// @Description  List all cases associated with a provider
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Provider ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/providers/{id}/cases [get]
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

// AdmitCase godoc
// @Summary      Admit a case
// @Description  Admit a case for inpatient treatment
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Case ID"
// @Param        request body claimsSchema.AdmitCaseRequest true "Admission details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/{id}/admit [put]
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

// UpdateCase godoc
// @Summary      Update a case
// @Description  Update case details by ID
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Case ID"
// @Param        request body claimsSchema.UpdateCaseRequest true "Updated case details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/{id} [put]
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

// StartTreatment godoc
// @Summary      Start treatment for a case
// @Description  Transition a case to treatment-in-progress status
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Case ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/{id}/start-treatment [put]
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

// DischargeCase godoc
// @Summary      Discharge a case
// @Description  Discharge a patient from a case
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Case ID"
// @Param        request body claimsSchema.DischargeCaseRequest true "Discharge details"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/{id}/discharge [put]
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

// CloseCase godoc
// @Summary      Close a case
// @Description  Close a case after all processing is complete
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        id path string true "Case ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/{id}/close [put]
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

// CountByStatus godoc
// @Summary      Count cases by status
// @Description  Get the count of cases filtered by status
// @Tags         Cases
// @Accept       json
// @Produce      json
// @Param        status query string true "Case status"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/cases/count [get]
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

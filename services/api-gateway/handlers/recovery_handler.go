package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecoveryHandler struct {
	recoverySvc service.RecoveryService
}

func NewRecoveryHandler(recoverySvc service.RecoveryService) *RecoveryHandler {
	return &RecoveryHandler{recoverySvc: recoverySvc}
}

// CreateRecovery godoc
// @Summary      Create a recovery
// @Description  Create a new reinsurance recovery record
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.CreateRecoveryRequest true "Create recovery request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries [post]
func (h *RecoveryHandler) CreateRecovery(ctx *gin.Context) {
	var req reinsuranceSchema.CreateRecoveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.CreateRecovery(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetRecovery godoc
// @Summary      Get a recovery
// @Description  Retrieve a single recovery by ID
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id} [get]
func (h *RecoveryHandler) GetRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	resp := h.recoverySvc.GetRecovery(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListRecoveries godoc
// @Summary      List recoveries
// @Description  Retrieve a paginated list of recoveries filtered by claim or treaty
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        claim query string false "Claim ID"
// @Param        treaty query string false "Treaty ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/recoveries [get]
func (h *RecoveryHandler) ListRecoveries(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	claimIDStr := ctx.Query("claim")
	treatyIDStr := ctx.Query("treaty")

	if claimIDStr != "" {
		claimID, err := uuid.Parse(claimIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
			return
		}
		resp := h.recoverySvc.ListByClaim(ctx.Request.Context(), claimID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.recoverySvc.GetRecoveryCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if treatyIDStr != "" {
		treatyID, err := uuid.Parse(treatyIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
			return
		}
		resp := h.recoverySvc.ListByTreaty(ctx.Request.Context(), treatyID, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.recoverySvc.GetRecoveryCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	utils.RespondError(ctx, http.StatusBadRequest, "Query parameter 'claim' or 'treaty' is required")
}

// ListOutstanding godoc
// @Summary      List outstanding recoveries
// @Description  Retrieve a paginated list of outstanding recoveries
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/outstanding [get]
func (h *RecoveryHandler) ListOutstanding(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	resp := h.recoverySvc.ListOutstanding(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.recoverySvc.GetRecoveryCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// ApplyRecoveryForClaim godoc
// @Summary      Apply recovery for a claim
// @Description  Apply reinsurance recovery for a specific claim
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        claimId path string true "Claim ID"
// @Param        request body reinsuranceSchema.ApplyRecoveryForClaimRequest true "Apply recovery request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/apply-for-claim/{claimId} [post]
func (h *RecoveryHandler) ApplyRecoveryForClaim(ctx *gin.Context) {
	claimID, err := uuid.Parse(ctx.Param("claimId"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid claim ID")
		return
	}

	var req reinsuranceSchema.ApplyRecoveryForClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.ApplyRecoveryForClaim(ctx.Request.Context(), claimID, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// AcknowledgeRecovery godoc
// @Summary      Acknowledge a recovery
// @Description  Acknowledge a recovery request by ID
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Param        request body reinsuranceSchema.RecoveryWorkflowRequest true "Workflow request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id}/acknowledge [put]
func (h *RecoveryHandler) AcknowledgeRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.AcknowledgeRecovery(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RequestInfo godoc
// @Summary      Request info for a recovery
// @Description  Request additional information for a recovery by ID
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Param        request body reinsuranceSchema.RecoveryWorkflowRequest true "Workflow request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id}/request-info [put]
func (h *RecoveryHandler) RequestInfo(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.RequestInfo(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ApproveRecovery godoc
// @Summary      Approve a recovery
// @Description  Approve a recovery request by ID
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Param        request body reinsuranceSchema.RecoveryWorkflowRequest true "Workflow request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id}/approve [put]
func (h *RecoveryHandler) ApproveRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.ApproveRecovery(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RecordPayment godoc
// @Summary      Record a recovery payment
// @Description  Record a payment received for a recovery by ID
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Param        request body reinsuranceSchema.RecordPaymentRequest true "Record payment request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id}/record-payment [put]
func (h *RecoveryHandler) RecordPayment(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecordPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.RecordPayment(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// WriteOffRecovery godoc
// @Summary      Write off a recovery
// @Description  Write off an unrecoverable recovery by ID
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Param        request body reinsuranceSchema.RecoveryWorkflowRequest true "Workflow request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id}/write-off [put]
func (h *RecoveryHandler) WriteOffRecovery(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	var req reinsuranceSchema.RecoveryWorkflowRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.recoverySvc.WriteOffRecovery(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// AgedAnalysis godoc
// @Summary      Get aged analysis of recoveries
// @Description  Retrieve aged analysis breakdown of outstanding recoveries
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/aged-analysis [get]
func (h *RecoveryHandler) AgedAnalysis(ctx *gin.Context) {
	resp := h.recoverySvc.GetAgedAnalysis(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetWorkflowEvents godoc
// @Summary      Get recovery workflow events
// @Description  Retrieve the workflow event history for a recovery
// @Tags         Recoveries
// @Accept       json
// @Produce      json
// @Param        id path string true "Recovery ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/recoveries/{id}/workflow [get]
func (h *RecoveryHandler) GetWorkflowEvents(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid recovery ID")
		return
	}

	resp := h.recoverySvc.GetWorkflowEvents(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

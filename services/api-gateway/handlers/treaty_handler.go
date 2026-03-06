package handlers

import (
	"net/http"

	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TreatyHandler struct {
	treatySvc service.TreatyService
}

func NewTreatyHandler(treatySvc service.TreatyService) *TreatyHandler {
	return &TreatyHandler{treatySvc: treatySvc}
}

// CreateTreaty godoc
// @Summary      Create a treaty
// @Description  Create a new reinsurance treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        request body reinsuranceSchema.CreateTreatyRequest true "Create treaty request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties [post]
func (h *TreatyHandler) CreateTreaty(ctx *gin.Context) {
	var req reinsuranceSchema.CreateTreatyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.CreateTreaty(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// GetTreaty godoc
// @Summary      Get a treaty
// @Description  Retrieve a single reinsurance treaty by ID
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id} [get]
func (h *TreatyHandler) GetTreaty(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.treatySvc.GetTreaty(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListTreaties godoc
// @Summary      List treaties
// @Description  Retrieve a paginated list of reinsurance treaties, optionally filtered by status or type
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Param        status query string false "Filter by treaty status"
// @Param        type query string false "Filter by treaty type"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties [get]
func (h *TreatyHandler) ListTreaties(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)

	status := ctx.Query("status")
	treatyType := ctx.Query("type")

	if status != "" {
		resp := h.treatySvc.ListTreatiesByStatus(ctx.Request.Context(), status, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.treatySvc.GetTreatyCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	if treatyType != "" {
		resp := h.treatySvc.ListTreatiesByType(ctx.Request.Context(), treatyType, pagination.Page, pagination.PageSize)
		if resp.Error != nil {
			utils.RespondError(ctx, resp.StatusCode, resp.Message)
			return
		}
		countResp := h.treatySvc.GetTreatyCount(ctx.Request.Context())
		utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
		return
	}

	resp := h.treatySvc.ListTreaties(ctx.Request.Context(), pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	countResp := h.treatySvc.GetTreatyCount(ctx.Request.Context())
	utils.RespondPaginated(ctx, resp.Message, resp.Data, pagination.Page, pagination.PageSize, countResp.Data)
}

// UpdateTreaty godoc
// @Summary      Update a treaty
// @Description  Update an existing reinsurance treaty by ID
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        request body reinsuranceSchema.UpdateTreatyRequest true "Update treaty request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id} [put]
func (h *TreatyHandler) UpdateTreaty(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	var req reinsuranceSchema.UpdateTreatyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.UpdateTreaty(ctx.Request.Context(), id, req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ActivateTreaty godoc
// @Summary      Activate a treaty
// @Description  Activate a draft reinsurance treaty by ID
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/activate [put]
func (h *TreatyHandler) ActivateTreaty(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.treatySvc.ActivateTreaty(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// TerminateTreaty godoc
// @Summary      Terminate a treaty
// @Description  Terminate an active reinsurance treaty by ID
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/terminate [put]
func (h *TreatyHandler) TerminateTreaty(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.treatySvc.TerminateTreaty(ctx.Request.Context(), id, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ExpireOverdue godoc
// @Summary      Expire overdue treaties
// @Description  Mark all overdue treaties as expired
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/expire [post]
func (h *TreatyHandler) ExpireOverdue(ctx *gin.Context) {
	resp := h.treatySvc.ExpireOverdue(ctx.Request.Context())
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// --- Participants ---

// AddParticipant godoc
// @Summary      Add a participant to a treaty
// @Description  Add a new reinsurer participant to a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        request body reinsuranceSchema.AddParticipantRequest true "Add participant request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/participants [post]
func (h *TreatyHandler) AddParticipant(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	var req reinsuranceSchema.AddParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.AddParticipant(ctx.Request.Context(), treatyID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListParticipants godoc
// @Summary      List treaty participants
// @Description  Retrieve all participants for a specific treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/participants [get]
func (h *TreatyHandler) ListParticipants(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.treatySvc.ListParticipants(ctx.Request.Context(), treatyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateParticipant godoc
// @Summary      Update a treaty participant
// @Description  Update an existing participant in a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        participantId path string true "Participant ID"
// @Param        request body reinsuranceSchema.UpdateParticipantRequest true "Update participant request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/participants/{participantId} [put]
func (h *TreatyHandler) UpdateParticipant(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid participant ID")
		return
	}

	var req reinsuranceSchema.UpdateParticipantRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.UpdateParticipant(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RemoveParticipant godoc
// @Summary      Remove a treaty participant
// @Description  Remove a participant from a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        participantId path string true "Participant ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/participants/{participantId} [delete]
func (h *TreatyHandler) RemoveParticipant(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid participant ID")
		return
	}

	resp := h.treatySvc.RemoveParticipant(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// --- Layers ---

// AddLayer godoc
// @Summary      Add a layer to a treaty
// @Description  Add a new layer to a reinsurance treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        request body reinsuranceSchema.AddLayerRequest true "Add layer request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/layers [post]
func (h *TreatyHandler) AddLayer(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	var req reinsuranceSchema.AddLayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.AddLayer(ctx.Request.Context(), treatyID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListLayers godoc
// @Summary      List treaty layers
// @Description  Retrieve all layers for a specific treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/layers [get]
func (h *TreatyHandler) ListLayers(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.treatySvc.ListLayers(ctx.Request.Context(), treatyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// UpdateLayer godoc
// @Summary      Update a treaty layer
// @Description  Update an existing layer in a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        layerId path string true "Layer ID"
// @Param        request body reinsuranceSchema.UpdateLayerRequest true "Update layer request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/layers/{layerId} [put]
func (h *TreatyHandler) UpdateLayer(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid layer ID")
		return
	}

	var req reinsuranceSchema.UpdateLayerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.UpdateLayer(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RemoveLayer godoc
// @Summary      Remove a treaty layer
// @Description  Remove a layer from a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        layerId path string true "Layer ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/layers/{layerId} [delete]
func (h *TreatyHandler) RemoveLayer(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid layer ID")
		return
	}

	resp := h.treatySvc.RemoveLayer(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// --- Profit Commission Rules ---

// AddProfitCommissionRule godoc
// @Summary      Add a profit commission rule
// @Description  Add a new profit commission rule to a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        request body reinsuranceSchema.AddProfitCommissionRuleRequest true "Add profit commission rule request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/profit-commission-rules [post]
func (h *TreatyHandler) AddProfitCommissionRule(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	var req reinsuranceSchema.AddProfitCommissionRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.treatySvc.AddProfitCommissionRule(ctx.Request.Context(), treatyID, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// ListProfitCommissionRules godoc
// @Summary      List profit commission rules
// @Description  Retrieve all profit commission rules for a specific treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/profit-commission-rules [get]
func (h *TreatyHandler) ListProfitCommissionRules(ctx *gin.Context) {
	treatyID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid treaty ID")
		return
	}

	resp := h.treatySvc.ListProfitCommissionRules(ctx.Request.Context(), treatyID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

// RemoveProfitCommissionRule godoc
// @Summary      Remove a profit commission rule
// @Description  Remove a profit commission rule from a treaty
// @Tags         Treaties
// @Accept       json
// @Produce      json
// @Param        id path string true "Treaty ID"
// @Param        ruleId path string true "Profit Commission Rule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/treaties/{id}/profit-commission-rules/{ruleId} [delete]
func (h *TreatyHandler) RemoveProfitCommissionRule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid profit commission rule ID")
		return
	}

	resp := h.treatySvc.RemoveProfitCommissionRule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

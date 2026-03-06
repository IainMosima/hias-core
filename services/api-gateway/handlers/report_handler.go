package handlers

import (
	"fmt"
	"net/http"

	reportSchema "github.com/bitbiz/hias-core/domains/reporting/schema"
	"github.com/bitbiz/hias-core/domains/reporting/service"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportSvc service.ReportService
}

func NewReportHandler(reportSvc service.ReportService) *ReportHandler {
	return &ReportHandler{reportSvc: reportSvc}
}

func (h *ReportHandler) ListDefinitions(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	category := ctx.Query("category")
	role := getUserRole(ctx)

	resp := h.reportSvc.ListDefinitions(ctx.Request.Context(), category, role, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) GetDefinition(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid report definition ID")
		return
	}

	resp := h.reportSvc.GetDefinition(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) CreateAdHocDefinition(ctx *gin.Context) {
	var req reportSchema.CreateAdHocReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.reportSvc.CreateAdHocDefinition(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) GenerateReport(ctx *gin.Context) {
	var req reportSchema.GenerateReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.reportSvc.GenerateReport(ctx.Request.Context(), req, getUserID(ctx), getUserRole(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) PreviewReport(ctx *gin.Context) {
	var req reportSchema.PreviewReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	role := getUserRole(ctx)
	resp := h.reportSvc.PreviewReport(ctx.Request.Context(), req.ReportCode, req.Parameters, role)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) DrillDown(ctx *gin.Context) {
	var req reportSchema.DrillDownRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.reportSvc.DrillDown(ctx.Request.Context(), req, getUserID(ctx), getUserRole(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) ListGeneratedReports(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	userID := getUserID(ctx)

	var defID *uuid.UUID
	if defIDStr := ctx.Query("definition_id"); defIDStr != "" {
		id, err := uuid.Parse(defIDStr)
		if err != nil {
			utils.RespondError(ctx, http.StatusBadRequest, "Invalid definition ID")
			return
		}
		defID = &id
	}

	resp := h.reportSvc.ListGeneratedReports(ctx.Request.Context(), defID, pagination.Page, pagination.PageSize, userID)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) GetGeneratedReport(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid report ID")
		return
	}

	resp := h.reportSvc.GetGeneratedReport(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) DownloadReport(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid report ID")
		return
	}

	data, format, reportNumber, err := h.reportSvc.DownloadReport(ctx.Request.Context(), id)
	if err != nil {
		utils.RespondError(ctx, http.StatusNotFound, "Report file not found")
		return
	}

	var contentType string
	var ext string
	switch format {
	case "CSV":
		contentType = "text/csv"
		ext = "csv"
	case "XLSX":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		ext = "xlsx"
	case "PDF":
		contentType = "application/pdf"
		ext = "pdf"
	default:
		contentType = "application/octet-stream"
		ext = "bin"
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.%s", reportNumber, ext))
	ctx.Data(http.StatusOK, contentType, data)
}

func (h *ReportHandler) CreateSchedule(ctx *gin.Context) {
	var req reportSchema.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.reportSvc.CreateSchedule(ctx.Request.Context(), req, getUserID(ctx))
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) ListSchedules(ctx *gin.Context) {
	pagination := utils.GetPaginationParams(ctx)
	defIDStr := ctx.Query("definition_id")
	if defIDStr == "" {
		utils.RespondError(ctx, http.StatusBadRequest, "definition_id query parameter is required")
		return
	}

	defID, err := uuid.Parse(defIDStr)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid definition ID")
		return
	}

	resp := h.reportSvc.ListSchedules(ctx.Request.Context(), defID, pagination.Page, pagination.PageSize)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) UpdateSchedule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	var req reportSchema.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	resp := h.reportSvc.UpdateSchedule(ctx.Request.Context(), id, req)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) DeleteSchedule(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Invalid schedule ID")
		return
	}

	resp := h.reportSvc.DeleteSchedule(ctx.Request.Context(), id)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

func (h *ReportHandler) GetManagementDashboard(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "year")

	resp := h.reportSvc.GetManagementDashboard(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

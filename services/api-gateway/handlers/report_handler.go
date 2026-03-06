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

// ListDefinitions godoc
// @Summary      List report definitions
// @Description  Retrieve a paginated list of report definitions, optionally filtered by category
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        category query string false "Filter by report category"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/definitions [get]
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

// GetDefinition godoc
// @Summary      Get a report definition
// @Description  Retrieve a single report definition by ID
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        id path string true "Report Definition ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/definitions/{id} [get]
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

// CreateAdHocDefinition godoc
// @Summary      Create an ad-hoc report definition
// @Description  Create a new ad-hoc report definition
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        request body reportSchema.CreateAdHocReportRequest true "Create ad-hoc report request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/definitions/adhoc [post]
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

// GenerateReport godoc
// @Summary      Generate a report
// @Description  Generate a report from a definition with specified parameters
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        request body reportSchema.GenerateReportRequest true "Generate report request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/generate [post]
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

// PreviewReport godoc
// @Summary      Preview a report
// @Description  Preview report data without generating a full report
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        request body reportSchema.PreviewReportRequest true "Preview report request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/preview [post]
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

// DrillDown godoc
// @Summary      Drill down into report data
// @Description  Drill down into a specific report row for detailed data
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        request body reportSchema.DrillDownRequest true "Drill down request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/drilldown [post]
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

// ListGeneratedReports godoc
// @Summary      List generated reports
// @Description  Retrieve a paginated list of generated reports, optionally filtered by definition
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        definition_id query string false "Filter by report definition ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/generated [get]
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

// GetGeneratedReport godoc
// @Summary      Get a generated report
// @Description  Retrieve a single generated report by ID
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        id path string true "Generated Report ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/generated/{id} [get]
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

// DownloadReport godoc
// @Summary      Download a generated report
// @Description  Download the file for a generated report by ID
// @Tags         Reports
// @Accept       json
// @Produce      application/octet-stream
// @Param        id path string true "Generated Report ID"
// @Success      200 {file} file
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/generated/{id}/download [get]
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

// CreateSchedule godoc
// @Summary      Create a report schedule
// @Description  Create a new scheduled report generation
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        request body reportSchema.CreateScheduleRequest true "Create schedule request"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/schedules [post]
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

// ListSchedules godoc
// @Summary      List report schedules
// @Description  Retrieve a paginated list of report schedules for a definition
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        definition_id query string true "Report Definition ID"
// @Param        page query int false "Page number"
// @Param        page_size query int false "Page size"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/schedules [get]
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

// UpdateSchedule godoc
// @Summary      Update a report schedule
// @Description  Update an existing report schedule by ID
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        id path string true "Schedule ID"
// @Param        request body reportSchema.UpdateScheduleRequest true "Update schedule request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/schedules/{id} [put]
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

// DeleteSchedule godoc
// @Summary      Delete a report schedule
// @Description  Delete a report schedule by ID
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        id path string true "Schedule ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/schedules/{id} [delete]
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

// GetManagementDashboard godoc
// @Summary      Get management dashboard
// @Description  Retrieve aggregated management dashboard data for a given period
// @Tags         Reports
// @Accept       json
// @Produce      json
// @Param        period query string false "Dashboard period (e.g., year, quarter, month)" default(year)
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Security     BearerAuth
// @Router       /api/v1/reports/dashboard [get]
func (h *ReportHandler) GetManagementDashboard(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "year")

	resp := h.reportSvc.GetManagementDashboard(ctx.Request.Context(), period)
	if resp.Error != nil {
		utils.RespondError(ctx, resp.StatusCode, resp.Message)
		return
	}
	utils.RespondSuccess(ctx, resp.StatusCode, resp.Message, resp.Data)
}

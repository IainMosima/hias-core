package reporting

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	analyticsRepo "github.com/bitbiz/hias-core/domains/analytics/repository"
	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	notificationService "github.com/bitbiz/hias-core/domains/notification/service"
	"github.com/bitbiz/hias-core/domains/reporting/entity"
	"github.com/bitbiz/hias-core/domains/reporting/repository"
	reportSchema "github.com/bitbiz/hias-core/domains/reporting/schema"
	reportService "github.com/bitbiz/hias-core/domains/reporting/service"
	reportingInfra "github.com/bitbiz/hias-core/infrastructures/reporting"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
	cronParser "github.com/robfig/cron/v3"
)

type reportServiceImpl struct {
	reportRepo      repository.ReportRepository
	reportDataRepo  repository.ReportDataRepository
	exporter        reportingInfra.ReportExporter
	notificationSvc notificationService.NotificationService
	auditSvc        auditService.AuditService
	analyticsRepo   analyticsRepo.AnalyticsRepository
}

func NewReportService(
	reportRepo repository.ReportRepository,
	reportDataRepo repository.ReportDataRepository,
	exporter reportingInfra.ReportExporter,
	notificationSvc notificationService.NotificationService,
	auditSvc auditService.AuditService,
	analyticsRepo analyticsRepo.AnalyticsRepository,
) reportService.ReportService {
	return &reportServiceImpl{
		reportRepo:      reportRepo,
		reportDataRepo:  reportDataRepo,
		exporter:        exporter,
		notificationSvc: notificationSvc,
		auditSvc:        auditSvc,
		analyticsRepo:   analyticsRepo,
	}
}

func (s *reportServiceImpl) ListDefinitions(ctx context.Context, category, role string, page, pageSize int) *schema.ServiceResponse[[]reportSchema.ReportDefinitionResponse] {
	offset := (page - 1) * pageSize
	defs, err := s.reportRepo.ListDefinitions(ctx, category, "", pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reportSchema.ReportDefinitionResponse](http.StatusInternalServerError, "Failed to list report definitions", err)
	}

	// Filter by role access
	var filtered []reportSchema.ReportDefinitionResponse
	for _, def := range defs {
		if role == "" || hasRole(def.AllowedRoles, role) {
			filtered = append(filtered, reportSchema.ToReportDefinitionResponse(def))
		}
	}
	if filtered == nil {
		filtered = []reportSchema.ReportDefinitionResponse{}
	}

	return schema.NewServiceResponse(filtered, http.StatusOK, "Report definitions retrieved")
}

func (s *reportServiceImpl) GetDefinition(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reportSchema.ReportDefinitionResponse] {
	def, err := s.reportRepo.GetDefinition(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportDefinitionResponse](http.StatusNotFound, "Report definition not found", err)
	}
	return schema.NewServiceResponse(reportSchema.ToReportDefinitionResponse(def), http.StatusOK, "Report definition retrieved")
}

func (s *reportServiceImpl) CreateAdHocDefinition(ctx context.Context, req reportSchema.CreateAdHocReportRequest, createdBy uuid.UUID) *schema.ServiceResponse[reportSchema.ReportDefinitionResponse] {
	code := fmt.Sprintf("ADHOC_%s_%d", req.Category, time.Now().UnixMilli())

	def := &entity.ReportDefinition{
		Code:              code,
		Name:              req.Name,
		Description:       req.Description,
		Category:          req.Category,
		ReportType:        string(shared.ReportTypeAdHoc),
		DefaultParameters: req.Filters,
		AllowedRoles:      req.AllowedRoles,
		Columns:           req.Columns,
		IsActive:          true,
		CreatedBy:         createdBy,
	}

	created, err := s.reportRepo.CreateDefinition(ctx, def)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportDefinitionResponse](http.StatusInternalServerError, "Failed to create ad-hoc report definition", err)
	}

	s.auditSvc.LogEvent(ctx, createdBy, string(shared.AuditEntityTypeReportDefinition), created.ID, string(shared.AuditActionCreate), nil, nil, "", "")

	return schema.NewServiceResponse(reportSchema.ToReportDefinitionResponse(created), http.StatusCreated, "Ad-hoc report definition created")
}

func (s *reportServiceImpl) GenerateReport(ctx context.Context, req reportSchema.GenerateReportRequest, generatedBy uuid.UUID, role string) *schema.ServiceResponse[reportSchema.GeneratedReportResponse] {
	// 1. Get definition
	def, err := s.reportRepo.GetDefinitionByCode(ctx, req.ReportCode)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusNotFound, "Report definition not found", err)
	}

	if role != "" && !hasRole(def.AllowedRoles, role) {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusForbidden, "Access denied to this report", fmt.Errorf("role %s not allowed", role))
	}

	// 2. Get report data
	data, err := s.fetchReportData(ctx, def.Code, req.Parameters)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to fetch report data", err)
	}

	// 3. Create generated report record
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	report := &entity.GeneratedReport{
		ReportDefinitionID: def.ID,
		ReportNumber:       utils.GenerateReportNumber(),
		Name:               def.Name,
		Parameters:         req.Parameters,
		Format:             req.Format,
		Status:             string(shared.ReportStatusGenerating),
		GeneratedBy:        generatedBy,
		ExpiresAt:          &expiresAt,
	}

	created, err := s.reportRepo.CreateGenerated(ctx, report)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to create report record", err)
	}

	// 4. Export to requested format
	var fileData []byte
	switch req.Format {
	case string(shared.ReportFormatCSV):
		fileData, err = s.exporter.ExportCSV(def.Columns, data)
	case string(shared.ReportFormatXLSX):
		fileData, err = s.exporter.ExportXLSX(def.Columns, data, def.Name)
	case string(shared.ReportFormatPDF):
		fileData, err = s.exporter.ExportPDF(def.Columns, data, def.Name, nil)
	}

	if err != nil {
		s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusFailed), 0, 0, err.Error())
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to export report", err)
	}

	// 5. Store file and update status
	if err := s.reportRepo.StoreReportFile(ctx, created.ID, fileData, int64(len(fileData))); err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to store report file", err)
	}

	updated, err := s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusCompleted), len(data), int64(len(fileData)), "")
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to update report status", err)
	}

	s.auditSvc.LogEvent(ctx, generatedBy, string(shared.AuditEntityTypeGeneratedReport), updated.ID, string(shared.AuditActionCreate), nil, nil, "", "")

	return schema.NewServiceResponse(reportSchema.ToGeneratedReportResponse(updated), http.StatusCreated, "Report generated successfully")
}

func (s *reportServiceImpl) PreviewReport(ctx context.Context, reportCode string, params json.RawMessage, role string) *schema.ServiceResponse[reportSchema.ReportPreviewResponse] {
	def, err := s.reportRepo.GetDefinitionByCode(ctx, reportCode)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportPreviewResponse](http.StatusNotFound, "Report definition not found", err)
	}

	if role != "" && !hasRole(def.AllowedRoles, role) {
		return schema.NewServiceErrorResponse[reportSchema.ReportPreviewResponse](http.StatusForbidden, "Access denied to this report", fmt.Errorf("role %s not allowed", role))
	}

	data, err := s.fetchReportData(ctx, def.Code, params)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportPreviewResponse](http.StatusInternalServerError, "Failed to fetch report data", err)
	}

	// Limit preview to 50 rows
	previewData := data
	if len(previewData) > 50 {
		previewData = previewData[:50]
	}

	resp := reportSchema.ReportPreviewResponse{
		Columns:  def.Columns,
		Data:     previewData,
		RowCount: len(data),
	}

	return schema.NewServiceResponse(resp, http.StatusOK, "Report preview retrieved")
}

func (s *reportServiceImpl) DrillDown(ctx context.Context, req reportSchema.DrillDownRequest, generatedBy uuid.UUID, role string) *schema.ServiceResponse[reportSchema.GeneratedReportResponse] {
	entityID, err := uuid.Parse(req.EntityID)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusBadRequest, "Invalid entity ID", err)
	}

	// Check role access against definition
	def, _ := s.reportRepo.GetDefinitionByCode(ctx, req.ReportCode)
	if def != nil && role != "" && !hasRole(def.AllowedRoles, role) {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusForbidden, "Access denied to this report", fmt.Errorf("role %s not allowed", role))
	}

	start, end := parseDateRange(req.StartDate, req.EndDate)

	var data []map[string]interface{}
	var drillTitle string

	switch req.ReportCode {
	case "CLAIMS_EXPERIENCE", "CLAIMS_REGISTER":
		data, err = s.reportDataRepo.DrillDownClaimsByPolicy(ctx, entityID, start, end)
		drillTitle = "Claims Drill-Down"
	case "PREMIUM_REGISTER", "PREMIUM_DEBTORS_AGEING":
		data, err = s.reportDataRepo.DrillDownPaymentsByPolicy(ctx, entityID, start, end)
		drillTitle = "Payments Drill-Down"
	default:
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusBadRequest, "Drill-down not supported for this report", fmt.Errorf("unsupported: %s", req.ReportCode))
	}

	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to fetch drill-down data", err)
	}

	// Build column definitions for drill-down
	drillColumns := buildDrillDownColumns(req.ReportCode)

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	defID := uuid.Nil
	if def != nil {
		defID = def.ID
	}

	report := &entity.GeneratedReport{
		ReportDefinitionID: defID,
		ReportNumber:       utils.GenerateReportNumber(),
		Name:               drillTitle,
		Parameters:         json.RawMessage(fmt.Sprintf(`{"entity_id":"%s"}`, req.EntityID)),
		Format:             req.Format,
		Status:             string(shared.ReportStatusGenerating),
		GeneratedBy:        generatedBy,
		ExpiresAt:          &expiresAt,
	}

	created, err := s.reportRepo.CreateGenerated(ctx, report)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to create report record", err)
	}

	var fileData []byte
	switch req.Format {
	case string(shared.ReportFormatCSV):
		fileData, err = s.exporter.ExportCSV(drillColumns, data)
	case string(shared.ReportFormatXLSX):
		fileData, err = s.exporter.ExportXLSX(drillColumns, data, drillTitle)
	case string(shared.ReportFormatPDF):
		fileData, err = s.exporter.ExportPDF(drillColumns, data, drillTitle, nil)
	}

	if err != nil {
		s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusFailed), 0, 0, err.Error())
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to export drill-down report", err)
	}

	if err := s.reportRepo.StoreReportFile(ctx, created.ID, fileData, int64(len(fileData))); err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to store drill-down report file", err)
	}

	updated, err := s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusCompleted), len(data), int64(len(fileData)), "")
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to update drill-down report status", err)
	}

	s.auditSvc.LogEvent(ctx, generatedBy, string(shared.AuditEntityTypeGeneratedReport), updated.ID, string(shared.AuditActionCreate), nil, nil, "", "")

	return schema.NewServiceResponse(reportSchema.ToGeneratedReportResponse(updated), http.StatusCreated, "Drill-down report generated")
}

func (s *reportServiceImpl) ListGeneratedReports(ctx context.Context, defID *uuid.UUID, page, pageSize int, userID uuid.UUID) *schema.ServiceResponse[[]reportSchema.GeneratedReportResponse] {
	offset := (page - 1) * pageSize
	reports, err := s.reportRepo.ListGenerated(ctx, defID, "", &userID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reportSchema.GeneratedReportResponse](http.StatusInternalServerError, "Failed to list generated reports", err)
	}

	responses := make([]reportSchema.GeneratedReportResponse, len(reports))
	for i, r := range reports {
		responses[i] = reportSchema.ToGeneratedReportResponse(r)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Generated reports retrieved")
}

func (s *reportServiceImpl) GetGeneratedReport(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reportSchema.GeneratedReportResponse] {
	report, err := s.reportRepo.GetGenerated(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.GeneratedReportResponse](http.StatusNotFound, "Generated report not found", err)
	}
	return schema.NewServiceResponse(reportSchema.ToGeneratedReportResponse(report), http.StatusOK, "Generated report retrieved")
}

func (s *reportServiceImpl) DownloadReport(ctx context.Context, id uuid.UUID) ([]byte, string, string, error) {
	return s.reportRepo.GetReportFile(ctx, id)
}

// --- Schedules ---

func (s *reportServiceImpl) CreateSchedule(ctx context.Context, req reportSchema.CreateScheduleRequest, createdBy uuid.UUID) *schema.ServiceResponse[reportSchema.ReportScheduleResponse] {
	defID, err := uuid.Parse(req.ReportDefinitionID)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportScheduleResponse](http.StatusBadRequest, "Invalid report definition ID", err)
	}

	def, err := s.reportRepo.GetDefinition(ctx, defID)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportScheduleResponse](http.StatusNotFound, "Report definition not found", err)
	}

	recipients := make([]uuid.UUID, len(req.Recipients))
	for i, r := range req.Recipients {
		id, err := uuid.Parse(r)
		if err != nil {
			return schema.NewServiceErrorResponse[reportSchema.ReportScheduleResponse](http.StatusBadRequest, fmt.Sprintf("Invalid recipient ID: %s", r), err)
		}
		recipients[i] = id
	}

	nextRun := computeNextRun(req.CronExpression)
	sched := &entity.ReportSchedule{
		ReportDefinitionID: defID,
		Name:               req.Name,
		CronExpression:     req.CronExpression,
		Parameters:         req.Parameters,
		ExportFormat:       req.ExportFormat,
		Recipients:         recipients,
		IsActive:           true,
		NextRunAt:          &nextRun,
		CreatedBy:          createdBy,
	}

	created, err := s.reportRepo.CreateSchedule(ctx, sched)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportScheduleResponse](http.StatusInternalServerError, "Failed to create report schedule", err)
	}

	s.auditSvc.LogEvent(ctx, createdBy, string(shared.AuditEntityTypeReportSchedule), created.ID, string(shared.AuditActionCreate), nil, nil, "", "")

	return schema.NewServiceResponse(reportSchema.ToReportScheduleResponse(created, def.Name), http.StatusCreated, "Report schedule created")
}

func (s *reportServiceImpl) ListSchedules(ctx context.Context, defID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reportSchema.ReportScheduleResponse] {
	offset := (page - 1) * pageSize
	scheds, err := s.reportRepo.ListSchedulesByDefinition(ctx, defID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reportSchema.ReportScheduleResponse](http.StatusInternalServerError, "Failed to list report schedules", err)
	}

	def, _ := s.reportRepo.GetDefinition(ctx, defID)
	defName := ""
	if def != nil {
		defName = def.Name
	}

	responses := make([]reportSchema.ReportScheduleResponse, len(scheds))
	for i, sched := range scheds {
		responses[i] = reportSchema.ToReportScheduleResponse(sched, defName)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Report schedules retrieved")
}

func (s *reportServiceImpl) UpdateSchedule(ctx context.Context, id uuid.UUID, req reportSchema.UpdateScheduleRequest) *schema.ServiceResponse[reportSchema.ReportScheduleResponse] {
	existing, err := s.reportRepo.GetSchedule(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportScheduleResponse](http.StatusNotFound, "Report schedule not found", err)
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.CronExpression != "" {
		existing.CronExpression = req.CronExpression
		nextRun := computeNextRun(req.CronExpression)
		existing.NextRunAt = &nextRun
	}
	if req.Parameters != nil {
		existing.Parameters = req.Parameters
	}
	if req.ExportFormat != "" {
		existing.ExportFormat = req.ExportFormat
	}
	if req.Recipients != nil {
		recipients := make([]uuid.UUID, len(req.Recipients))
		for i, r := range req.Recipients {
			recipients[i], _ = uuid.Parse(r)
		}
		existing.Recipients = recipients
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	updated, err := s.reportRepo.UpdateSchedule(ctx, existing)
	if err != nil {
		return schema.NewServiceErrorResponse[reportSchema.ReportScheduleResponse](http.StatusInternalServerError, "Failed to update report schedule", err)
	}

	def, _ := s.reportRepo.GetDefinition(ctx, updated.ReportDefinitionID)
	defName := ""
	if def != nil {
		defName = def.Name
	}

	return schema.NewServiceResponse(reportSchema.ToReportScheduleResponse(updated, defName), http.StatusOK, "Report schedule updated")
}

func (s *reportServiceImpl) DeleteSchedule(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[string] {
	err := s.reportRepo.DeleteSchedule(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to delete report schedule", err)
	}
	return schema.NewServiceResponse("deleted", http.StatusOK, "Report schedule deleted")
}

// --- Management Dashboard ---

func (s *reportServiceImpl) GetManagementDashboard(ctx context.Context, period string) *schema.ServiceResponse[reportSchema.ManagementDashboardResponse] {
	start, end := periodToDateRange(period)

	claimsVol, err := s.analyticsRepo.GetClaimsVolume(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get claims volume: %v", err)
	}
	approvalRate, err := s.analyticsRepo.GetApprovalRate(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get approval rate: %v", err)
	}
	avgTAT, err := s.analyticsRepo.GetAverageTAT(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get average TAT: %v", err)
	}
	lossRatio, err := s.analyticsRepo.GetLossRatio(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get loss ratio: %v", err)
	}
	totalPremium, err := s.analyticsRepo.GetTotalPremiumCollected(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get total premium collected: %v", err)
	}
	totalClaimsPaid, err := s.analyticsRepo.GetTotalClaimsPaid(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get total claims paid: %v", err)
	}
	activePolicies, err := s.analyticsRepo.GetActivePolicyCount(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get active policy count: %v", err)
	}
	totalMembers, err := s.analyticsRepo.GetTotalMemberCount(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get total member count: %v", err)
	}
	renewalRate, err := s.analyticsRepo.GetRenewalRate(ctx, start, end)
	if err != nil {
		log.Printf("Dashboard: failed to get renewal rate: %v", err)
	}
	outstandingPremium, err := s.reportDataRepo.GetOutstandingPremium(ctx)
	if err != nil {
		log.Printf("Dashboard: failed to get outstanding premium: %v", err)
	}
	slaBreachCount, err := s.reportDataRepo.GetSLABreachCount(ctx)
	if err != nil {
		log.Printf("Dashboard: failed to get SLA breach count: %v", err)
	}

	var claimsVolume int64
	if claimsVol != nil {
		claimsVolume = claimsVol.TotalClaims
	}

	// Compute premium growth
	prevStart := start.AddDate(-1, 0, 0)
	prevEnd := end.AddDate(-1, 0, 0)
	prevPremium, err := s.analyticsRepo.GetTotalPremiumCollected(ctx, prevStart, prevEnd)
	if err != nil {
		log.Printf("Dashboard: failed to get previous period premium: %v", err)
	}
	var premiumGrowth float64
	if prevPremium > 0 {
		premiumGrowth = float64(totalPremium-prevPremium) / float64(prevPremium) * 100
	}

	dashboard := reportSchema.ManagementDashboardResponse{
		LossRatio:          lossRatio,
		ClaimsVolume:       claimsVolume,
		ApprovalRate:       approvalRate,
		AvgTATHours:        avgTAT,
		TotalPremium:       totalPremium,
		TotalClaimsPaid:    totalClaimsPaid,
		ActivePolicies:     activePolicies,
		TotalMembers:       totalMembers,
		RenewalRate:        renewalRate,
		PremiumGrowth:      premiumGrowth,
		OutstandingPremium: outstandingPremium,
		SLABreachCount:     slaBreachCount,
	}

	return schema.NewServiceResponse(dashboard, http.StatusOK, "Management dashboard retrieved")
}

// --- Scheduled Execution ---

func (s *reportServiceImpl) ExecuteDueSchedules(ctx context.Context) error {
	schedules, err := s.reportRepo.ListDueSchedules(ctx)
	if err != nil {
		return fmt.Errorf("failed to list due schedules: %w", err)
	}

	for _, sched := range schedules {
		def, err := s.reportRepo.GetDefinition(ctx, sched.ReportDefinitionID)
		if err != nil {
			log.Printf("Schedule %s: failed to get definition: %v", sched.Name, err)
			continue
		}

		data, err := s.fetchReportData(ctx, def.Code, sched.Parameters)
		if err != nil {
			log.Printf("Schedule %s: failed to fetch data: %v", sched.Name, err)
			continue
		}

		expiresAt := time.Now().Add(30 * 24 * time.Hour)
		report := &entity.GeneratedReport{
			ReportDefinitionID: def.ID,
			ScheduleID:         &sched.ID,
			ReportNumber:       utils.GenerateReportNumber(),
			Name:               fmt.Sprintf("%s (Scheduled)", def.Name),
			Parameters:         sched.Parameters,
			Format:             sched.ExportFormat,
			Status:             string(shared.ReportStatusGenerating),
			GeneratedBy:        sched.CreatedBy,
			ExpiresAt:          &expiresAt,
		}

		created, err := s.reportRepo.CreateGenerated(ctx, report)
		if err != nil {
			log.Printf("Schedule %s: failed to create report: %v", sched.Name, err)
			continue
		}

		var fileData []byte
		switch sched.ExportFormat {
		case string(shared.ReportFormatCSV):
			fileData, err = s.exporter.ExportCSV(def.Columns, data)
		case string(shared.ReportFormatXLSX):
			fileData, err = s.exporter.ExportXLSX(def.Columns, data, def.Name)
		case string(shared.ReportFormatPDF):
			fileData, err = s.exporter.ExportPDF(def.Columns, data, def.Name, nil)
		}

		if err != nil {
			s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusFailed), 0, 0, err.Error())
			log.Printf("Schedule %s: export failed: %v", sched.Name, err)
			continue
		}

		if err := s.reportRepo.StoreReportFile(ctx, created.ID, fileData, int64(len(fileData))); err != nil {
			log.Printf("Schedule %s: failed to store report file: %v", sched.Name, err)
			s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusFailed), 0, 0, err.Error())
			continue
		}

		if _, err := s.reportRepo.UpdateGeneratedStatus(ctx, created.ID, string(shared.ReportStatusCompleted), len(data), int64(len(fileData)), ""); err != nil {
			log.Printf("Schedule %s: failed to update report status: %v", sched.Name, err)
			continue
		}

		// Notify recipients
		if s.notificationSvc != nil && len(sched.Recipients) > 0 {
			s.notificationSvc.SendBulk(ctx,
				sched.Recipients,
				string(shared.NotificationChannelInApp),
				"REPORT",
				fmt.Sprintf("Scheduled Report: %s", def.Name),
				fmt.Sprintf("Your scheduled report '%s' has been generated. Report #%s is ready for download.", def.Name, created.ReportNumber),
			)
		}

		// Update schedule last run
		nextRun := computeNextRun(sched.CronExpression)
		s.reportRepo.UpdateScheduleLastRun(ctx, sched.ID, time.Now(), nextRun)

		log.Printf("Schedule %s: report %s generated (%d rows)", sched.Name, created.ReportNumber, len(data))
	}

	return nil
}

// --- Internal helpers ---

func (s *reportServiceImpl) fetchReportData(ctx context.Context, code string, params json.RawMessage) ([]map[string]interface{}, error) {
	p := parseReportParams(params)
	start, end := parseDateRange(p.StartDate, p.EndDate)
	limit := 10000
	offset := 0

	switch code {
	case "CLAIMS_EXPERIENCE":
		return s.reportDataRepo.GetClaimsExperienceData(ctx, start, end, p.PolicyID)
	case "CLAIMS_REGISTER":
		return s.reportDataRepo.GetClaimsRegisterData(ctx, start, end, p.Status, limit, offset)
	case "PREMIUM_DEBTORS_AGEING":
		return s.reportDataRepo.GetPremiumDebtorsAgeingData(ctx)
	case "PREMIUM_REGISTER":
		return s.reportDataRepo.GetPremiumRegisterData(ctx, start, end, limit, offset)
	case "MEMBERSHIP":
		return s.reportDataRepo.GetMembershipData(ctx, p.PolicyID, p.Status, limit, offset)
	case "PROVIDER_PERFORMANCE":
		return s.reportDataRepo.GetProviderPerformanceData(ctx, start, end)
	case "LOSS_RATIO":
		return s.reportDataRepo.GetLossRatioData(ctx, start, end)
	case "RENEWAL":
		return s.reportDataRepo.GetRenewalData(ctx, start, end)
	default:
		return nil, fmt.Errorf("unknown report code: %s", code)
	}
}

type reportParams struct {
	StartDate string     `json:"start_date"`
	EndDate   string     `json:"end_date"`
	PolicyID  *uuid.UUID `json:"policy_id,omitempty"`
	Status    string     `json:"status"`
	Period    string     `json:"period"`
}

func parseReportParams(params json.RawMessage) reportParams {
	var p reportParams
	if params != nil {
		json.Unmarshal(params, &p)
	}
	return p
}

func parseDateRange(startStr, endStr string) (time.Time, time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	end := now

	if startStr != "" {
		if t, err := utils.ParseFlexibleDate(startStr); err == nil {
			start = t
		}
	}
	if endStr != "" {
		if t, err := utils.ParseFlexibleDate(endStr); err == nil {
			end = t.Add(24*time.Hour - time.Second)
		}
	}

	return start, end
}

func periodToDateRange(period string) (time.Time, time.Time) {
	now := time.Now()
	end := now

	switch period {
	case "month":
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC), end
	case "quarter":
		q := (now.Month()-1)/3*3 + 1
		return time.Date(now.Year(), q, 1, 0, 0, 0, 0, time.UTC), end
	case "year":
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), end
	default:
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC), end
	}
}

func computeNextRun(cronExpr string) time.Time {
	parser := cronParser.NewParser(cronParser.Minute | cronParser.Hour | cronParser.Dom | cronParser.Month | cronParser.Dow)
	sched, err := parser.Parse(cronExpr)
	if err != nil {
		log.Printf("computeNextRun: failed to parse cron expression %q: %v, defaulting to 1 hour", cronExpr, err)
		return time.Now().Add(1 * time.Hour)
	}
	return sched.Next(time.Now())
}

func hasRole(allowedRoles []string, role string) bool {
	for _, r := range allowedRoles {
		if r == role {
			return true
		}
	}
	return false
}

func buildDrillDownColumns(reportCode string) json.RawMessage {
	switch reportCode {
	case "CLAIMS_EXPERIENCE", "CLAIMS_REGISTER":
		return json.RawMessage(`[
			{"name":"claim_number","label":"Claim No","type":"string"},
			{"name":"member_name","label":"Member","type":"string"},
			{"name":"provider_name","label":"Provider","type":"string"},
			{"name":"claim_type","label":"Type","type":"string"},
			{"name":"service_date","label":"Service Date","type":"date"},
			{"name":"total_amount","label":"Claimed","type":"money"},
			{"name":"approved_amount","label":"Approved","type":"money"},
			{"name":"co_pay_amount","label":"Co-Pay","type":"money"},
			{"name":"status","label":"Status","type":"string"},
			{"name":"created_at","label":"Submitted","type":"datetime"}
		]`)
	case "PREMIUM_REGISTER", "PREMIUM_DEBTORS_AGEING":
		return json.RawMessage(`[
			{"name":"reference_number","label":"Reference","type":"string"},
			{"name":"amount","label":"Amount","type":"money"},
			{"name":"method","label":"Method","type":"string"},
			{"name":"status","label":"Status","type":"string"},
			{"name":"paid_at","label":"Paid At","type":"datetime"},
			{"name":"invoice_number","label":"Invoice","type":"string"},
			{"name":"invoice_amount","label":"Invoice Amt","type":"money"},
			{"name":"due_date","label":"Due Date","type":"date"}
		]`)
	default:
		return json.RawMessage(`[]`)
	}
}

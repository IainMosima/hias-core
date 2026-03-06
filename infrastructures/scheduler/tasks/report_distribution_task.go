package tasks

import (
	"context"
	"fmt"
	"log"

	reportService "github.com/bitbiz/hias-core/domains/reporting/service"
)

// ReportDistributionTask executes due report schedules and distributes them.
type ReportDistributionTask struct {
	schedule  string
	reportSvc reportService.ReportService
}

func NewReportDistributionTask(
	schedule string,
	reportSvc reportService.ReportService,
) *ReportDistributionTask {
	return &ReportDistributionTask{
		schedule:  schedule,
		reportSvc: reportSvc,
	}
}

func (t *ReportDistributionTask) Name() string     { return "report-distribution" }
func (t *ReportDistributionTask) Schedule() string { return t.schedule }

func (t *ReportDistributionTask) Execute(ctx context.Context) error {
	err := t.reportSvc.ExecuteDueSchedules(ctx)
	if err != nil {
		return fmt.Errorf("report distribution failed: %w", err)
	}
	log.Println("Report distribution task completed")
	return nil
}

// ReportCleanupTask deletes expired generated reports.
type ReportCleanupTask struct {
	schedule   string
	reportRepo interface {
		DeleteExpiredReports(ctx context.Context) error
	}
}

func NewReportCleanupTask(
	schedule string,
	reportRepo interface {
		DeleteExpiredReports(ctx context.Context) error
	},
) *ReportCleanupTask {
	return &ReportCleanupTask{
		schedule:   schedule,
		reportRepo: reportRepo,
	}
}

func (t *ReportCleanupTask) Name() string     { return "report-cleanup" }
func (t *ReportCleanupTask) Schedule() string { return t.schedule }

func (t *ReportCleanupTask) Execute(ctx context.Context) error {
	err := t.reportRepo.DeleteExpiredReports(ctx)
	if err != nil {
		return fmt.Errorf("report cleanup failed: %w", err)
	}
	log.Println("Report cleanup task completed")
	return nil
}

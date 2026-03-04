package billing

import (
	"context"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/domains/billing/entity"
	"github.com/bitbiz/hias-core/domains/billing/repository"
	billingSchema "github.com/bitbiz/hias-core/domains/billing/schema"
	"github.com/bitbiz/hias-core/domains/billing/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	policyRepo "github.com/bitbiz/hias-core/domains/policy/repository"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type installmentServiceImpl struct {
	scheduleRepo    repository.InstallmentScheduleRepository
	installmentRepo repository.InstallmentRepository
	policyRepo      policyRepo.PolicyRepository
}

func NewInstallmentService(
	scheduleRepo repository.InstallmentScheduleRepository,
	installmentRepo repository.InstallmentRepository,
	policyRepo policyRepo.PolicyRepository,
) service.InstallmentService {
	return &installmentServiceImpl{
		scheduleRepo:    scheduleRepo,
		installmentRepo: installmentRepo,
		policyRepo:      policyRepo,
	}
}

func (s *installmentServiceImpl) CreateSchedule(ctx context.Context, req billingSchema.CreateInstallmentScheduleRequest, createdBy uuid.UUID) *schema.ServiceResponse[billingSchema.InstallmentScheduleResponse] {
	policyID, err := uuid.Parse(req.PolicyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InstallmentScheduleResponse](http.StatusBadRequest, "Invalid policy ID", err)
	}

	pol, err := s.policyRepo.GetByID(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InstallmentScheduleResponse](http.StatusNotFound, "Policy not found", err)
	}

	// Calculate installment count based on frequency
	var totalInstallments int
	switch req.Frequency {
	case string(shared.BillingFrequencyMonthly):
		totalInstallments = shared.InstallmentsPerMonth
	case string(shared.BillingFrequencyQuarterly):
		totalInstallments = shared.InstallmentsPerQuarter
	case string(shared.BillingFrequencySemiAnnual):
		totalInstallments = shared.InstallmentsPerSemiAnnual
	case string(shared.BillingFrequencyAnnual):
		totalInstallments = shared.InstallmentsPerAnnual
	default:
		return schema.NewServiceErrorResponse[billingSchema.InstallmentScheduleResponse](http.StatusBadRequest, "Invalid frequency. Use: monthly, quarterly, semi_annual, annual", nil)
	}

	amountPerInstallment := pol.PremiumAmount / int64(totalInstallments)

	startDate := req.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	schedule := &entity.InstallmentSchedule{
		PolicyID:             policyID,
		Frequency:            req.Frequency,
		TotalInstallments:    totalInstallments,
		AmountPerInstallment: amountPerInstallment,
		StartDate:            startDate,
		Status:               string(shared.InstallmentScheduleStatusActive),
		CreatedBy:            createdBy,
	}

	createdSchedule, err := s.scheduleRepo.Create(ctx, schedule)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InstallmentScheduleResponse](http.StatusInternalServerError, "Failed to create installment schedule", err)
	}

	// Auto-generate installments
	var installmentResponses []billingSchema.InstallmentResponse
	for i := 1; i <= totalInstallments; i++ {
		var dueDate time.Time
		switch req.Frequency {
		case string(shared.BillingFrequencyMonthly):
			dueDate = startDate.AddDate(0, i-1, 0)
		case string(shared.BillingFrequencyQuarterly):
			dueDate = startDate.AddDate(0, (i-1)*3, 0)
		case string(shared.BillingFrequencySemiAnnual):
			dueDate = startDate.AddDate(0, (i-1)*6, 0)
		case string(shared.BillingFrequencyAnnual):
			dueDate = startDate.AddDate(i-1, 0, 0)
		}

		inst := &entity.Installment{
			ScheduleID:        createdSchedule.ID,
			InstallmentNumber: i,
			DueDate:           dueDate,
			Amount:            amountPerInstallment,
			Status:            string(shared.InstallmentStatusPending),
		}

		createdInst, err := s.installmentRepo.Create(ctx, inst)
		if err != nil {
			return schema.NewServiceErrorResponse[billingSchema.InstallmentScheduleResponse](http.StatusInternalServerError, "Failed to create installment", err)
		}
		installmentResponses = append(installmentResponses, billingSchema.ToInstallmentResponse(createdInst))
	}

	resp := billingSchema.ToInstallmentScheduleResponse(createdSchedule)
	resp.Installments = installmentResponses

	return schema.NewServiceResponse(resp, http.StatusCreated, "Installment schedule created")
}

func (s *installmentServiceImpl) GetSchedulesByPolicy(ctx context.Context, policyID uuid.UUID) *schema.ServiceResponse[[]billingSchema.InstallmentScheduleResponse] {
	schedules, err := s.scheduleRepo.GetByPolicy(ctx, policyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.InstallmentScheduleResponse](http.StatusInternalServerError, "Failed to get installment schedules", err)
	}

	responses := make([]billingSchema.InstallmentScheduleResponse, len(schedules))
	for i, sched := range schedules {
		responses[i] = billingSchema.ToInstallmentScheduleResponse(sched)
		installments, _ := s.installmentRepo.ListBySchedule(ctx, sched.ID)
		instResponses := make([]billingSchema.InstallmentResponse, len(installments))
		for j, inst := range installments {
			instResponses[j] = billingSchema.ToInstallmentResponse(inst)
		}
		responses[i].Installments = instResponses
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Installment schedules retrieved")
}

func (s *installmentServiceImpl) ListInstallmentsBySchedule(ctx context.Context, scheduleID uuid.UUID) *schema.ServiceResponse[[]billingSchema.InstallmentResponse] {
	installments, err := s.installmentRepo.ListBySchedule(ctx, scheduleID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]billingSchema.InstallmentResponse](http.StatusInternalServerError, "Failed to list installments", err)
	}

	responses := make([]billingSchema.InstallmentResponse, len(installments))
	for i, inst := range installments {
		responses[i] = billingSchema.ToInstallmentResponse(inst)
	}

	return schema.NewServiceResponse(responses, http.StatusOK, "Installments retrieved")
}

func (s *installmentServiceImpl) MarkInstallmentPaid(ctx context.Context, installmentID uuid.UUID, invoiceID uuid.UUID) *schema.ServiceResponse[billingSchema.InstallmentResponse] {
	paid, err := s.installmentRepo.MarkPaid(ctx, installmentID, invoiceID)
	if err != nil {
		return schema.NewServiceErrorResponse[billingSchema.InstallmentResponse](http.StatusInternalServerError, "Failed to mark installment paid", err)
	}

	return schema.NewServiceResponse(billingSchema.ToInstallmentResponse(paid), http.StatusOK, "Installment marked as paid")
}

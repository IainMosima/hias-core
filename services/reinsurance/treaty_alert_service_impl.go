package reinsurance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/bitbiz/hias-core/domains/reinsurance/repository"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/google/uuid"
)

type treatyAlertServiceImpl struct {
	alertRepo    repository.TreatyAlertRepository
	layerRepo    repository.TreatyLayerRepository
	recoveryRepo repository.RecoveryRepository
	treatyRepo   repository.TreatyRepository
}

func NewTreatyAlertService(
	alertRepo repository.TreatyAlertRepository,
	layerRepo repository.TreatyLayerRepository,
	recoveryRepo repository.RecoveryRepository,
	treatyRepo repository.TreatyRepository,
) service.TreatyAlertService {
	return &treatyAlertServiceImpl{
		alertRepo:    alertRepo,
		layerRepo:    layerRepo,
		recoveryRepo: recoveryRepo,
		treatyRepo:   treatyRepo,
	}
}

func (s *treatyAlertServiceImpl) CheckTreatyLimits(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse] {
	layers, err := s.layerRepo.ListByTreaty(ctx, treatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyAlertResponse](http.StatusInternalServerError, "Failed to fetch layers", err)
	}

	var alerts []reinsuranceSchema.TreatyAlertResponse
	for _, layer := range layers {
		if layer.AggregateLimit == nil || *layer.AggregateLimit == 0 {
			continue
		}

		usagePercent := layer.AggregateUsed * 100 / *layer.AggregateLimit

		if usagePercent >= 100 {
			alert := &entity.TreatyAlert{
				TreatyID:       treatyID,
				TreatyLayerID:  layer.ID,
				AlertType:      string(shared.TreatyAlertTypeLimitBreach),
				Severity:       string(shared.TreatyAlertSeverityCritical),
				Message:        fmt.Sprintf("Layer %d aggregate limit breached: %d%% used", layer.LayerNumber, usagePercent),
				ThresholdValue: *layer.AggregateLimit,
				CurrentValue:   layer.AggregateUsed,
			}
			created, err := s.alertRepo.Create(ctx, alert)
			if err == nil {
				alerts = append(alerts, reinsuranceSchema.ToTreatyAlertResponse(created))
			}
		} else if usagePercent >= int64(shared.AggregateWarningPercent) {
			alert := &entity.TreatyAlert{
				TreatyID:       treatyID,
				TreatyLayerID:  layer.ID,
				AlertType:      string(shared.TreatyAlertTypeAggregateWarning),
				Severity:       string(shared.TreatyAlertSeverityHigh),
				Message:        fmt.Sprintf("Layer %d aggregate at %d%% of limit", layer.LayerNumber, usagePercent),
				ThresholdValue: *layer.AggregateLimit * int64(shared.AggregateWarningPercent) / 100,
				CurrentValue:   layer.AggregateUsed,
			}
			created, err := s.alertRepo.Create(ctx, alert)
			if err == nil {
				alerts = append(alerts, reinsuranceSchema.ToTreatyAlertResponse(created))
			}
		}
	}

	return schema.NewServiceResponse(alerts, http.StatusOK, fmt.Sprintf("Created %d alerts", len(alerts)))
}

func (s *treatyAlertServiceImpl) CheckCatastropheThresholds(ctx context.Context, treatyID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse] {
	totalRecoverable, _ := s.recoveryRepo.GetTotalRecoverableByTreaty(ctx, treatyID)

	var alerts []reinsuranceSchema.TreatyAlertResponse
	if totalRecoverable > shared.CatastropheThresholdCents {
		alert := &entity.TreatyAlert{
			TreatyID:       treatyID,
			AlertType:      string(shared.TreatyAlertTypeCatastropheThreshold),
			Severity:       string(shared.TreatyAlertSeverityCritical),
			Message:        fmt.Sprintf("Total recoverable amount %d exceeds catastrophe threshold %d", totalRecoverable, shared.CatastropheThresholdCents),
			ThresholdValue: shared.CatastropheThresholdCents,
			CurrentValue:   totalRecoverable,
		}
		created, err := s.alertRepo.Create(ctx, alert)
		if err == nil {
			alerts = append(alerts, reinsuranceSchema.ToTreatyAlertResponse(created))
		}
	}

	return schema.NewServiceResponse(alerts, http.StatusOK, "Catastrophe threshold check complete")
}

func (s *treatyAlertServiceImpl) ListAlerts(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse] {
	offset := (page - 1) * pageSize
	alerts, err := s.alertRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyAlertResponse](http.StatusInternalServerError, "Failed to list alerts", err)
	}

	responses := make([]reinsuranceSchema.TreatyAlertResponse, len(alerts))
	for i, a := range alerts {
		responses[i] = reinsuranceSchema.ToTreatyAlertResponse(a)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Alerts retrieved successfully")
}

func (s *treatyAlertServiceImpl) ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse] {
	offset := (page - 1) * pageSize
	alerts, err := s.alertRepo.ListByTreaty(ctx, treatyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyAlertResponse](http.StatusInternalServerError, "Failed to list alerts", err)
	}

	responses := make([]reinsuranceSchema.TreatyAlertResponse, len(alerts))
	for i, a := range alerts {
		responses[i] = reinsuranceSchema.ToTreatyAlertResponse(a)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Alerts retrieved successfully")
}

func (s *treatyAlertServiceImpl) ListUnacknowledged(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse] {
	offset := (page - 1) * pageSize
	alerts, err := s.alertRepo.ListUnacknowledged(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyAlertResponse](http.StatusInternalServerError, "Failed to list unacknowledged alerts", err)
	}

	responses := make([]reinsuranceSchema.TreatyAlertResponse, len(alerts))
	for i, a := range alerts {
		responses[i] = reinsuranceSchema.ToTreatyAlertResponse(a)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Unacknowledged alerts retrieved successfully")
}

func (s *treatyAlertServiceImpl) AcknowledgeAlert(ctx context.Context, id uuid.UUID, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.TreatyAlertResponse] {
	updated, err := s.alertRepo.Acknowledge(ctx, id, userID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.TreatyAlertResponse](http.StatusInternalServerError, "Failed to acknowledge alert", err)
	}

	resp := reinsuranceSchema.ToTreatyAlertResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Alert acknowledged successfully")
}

func (s *treatyAlertServiceImpl) CheckExpiryWarnings(ctx context.Context) *schema.ServiceResponse[[]reinsuranceSchema.TreatyAlertResponse] {
	expiringTreaties, err := s.treatyRepo.ListExpiring(ctx, 30, 1000, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.TreatyAlertResponse](http.StatusInternalServerError, "Failed to list expiring treaties", err)
	}

	now := time.Now()
	var alerts []reinsuranceSchema.TreatyAlertResponse
	for _, treaty := range expiringTreaties {
		daysUntilExpiry := int(treaty.ExpiryDate.Sub(now).Hours() / 24)

		alert := &entity.TreatyAlert{
			TreatyID:       treaty.ID,
			AlertType:      string(shared.TreatyAlertTypeExpiryWarning),
			Severity:       string(shared.TreatyAlertSeverityMedium),
			Message:        fmt.Sprintf("Treaty %s (%s) is expiring on %s", treaty.TreatyNumber, treaty.Name, treaty.ExpiryDate.Format("2006-01-02")),
			ThresholdValue: 30,
			CurrentValue:   int64(daysUntilExpiry),
		}
		created, err := s.alertRepo.Create(ctx, alert)
		if err == nil {
			alerts = append(alerts, reinsuranceSchema.ToTreatyAlertResponse(created))
		}
	}

	return schema.NewServiceResponse(alerts, http.StatusOK, fmt.Sprintf("Created %d expiry warning alerts", len(alerts)))
}

func (s *treatyAlertServiceImpl) CountUnacknowledged(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.alertRepo.CountUnacknowledged(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count unacknowledged alerts", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Unacknowledged alert count retrieved")
}

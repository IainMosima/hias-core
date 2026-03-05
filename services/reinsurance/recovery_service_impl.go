package reinsurance

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	auditService "github.com/bitbiz/hias-core/domains/audit/service"
	schema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/entity"
	"github.com/bitbiz/hias-core/domains/reinsurance/repository"
	reinsuranceSchema "github.com/bitbiz/hias-core/domains/reinsurance/schema"
	"github.com/bitbiz/hias-core/domains/reinsurance/service"
	"github.com/bitbiz/hias-core/shared"
	"github.com/bitbiz/hias-core/shared/utils"
	"github.com/google/uuid"
)

type recoveryServiceImpl struct {
	recoveryRepo      repository.RecoveryRepository
	workflowEventRepo repository.RecoveryWorkflowEventRepository
	treatyRepo        repository.TreatyRepository
	layerRepo         repository.TreatyLayerRepository
	participantRepo   repository.TreatyParticipantRepository
	cessionRepo       repository.CessionRepository
	auditSvc          auditService.AuditService
}

func NewRecoveryService(
	recoveryRepo repository.RecoveryRepository,
	workflowEventRepo repository.RecoveryWorkflowEventRepository,
	treatyRepo repository.TreatyRepository,
	layerRepo repository.TreatyLayerRepository,
	participantRepo repository.TreatyParticipantRepository,
	cessionRepo repository.CessionRepository,
	auditSvc auditService.AuditService,
) service.RecoveryService {
	return &recoveryServiceImpl{
		recoveryRepo:      recoveryRepo,
		workflowEventRepo: workflowEventRepo,
		treatyRepo:        treatyRepo,
		layerRepo:         layerRepo,
		participantRepo:   participantRepo,
		cessionRepo:       cessionRepo,
		auditSvc:          auditSvc,
	}
}

func (s *recoveryServiceImpl) CreateRecovery(ctx context.Context, req reinsuranceSchema.CreateRecoveryRequest, createdBy uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	claimID, err := uuid.Parse(req.ClaimID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusBadRequest, "Invalid claim ID", err)
	}

	treatyID, err := uuid.Parse(req.TreatyID)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusBadRequest, "Invalid treaty ID", err)
	}

	var treatyLayerID uuid.UUID
	if req.TreatyLayerID != "" {
		treatyLayerID, err = uuid.Parse(req.TreatyLayerID)
		if err != nil {
			return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusBadRequest, "Invalid treaty layer ID", err)
		}
	}

	var cessionID uuid.UUID
	if req.CessionID != "" {
		cessionID, err = uuid.Parse(req.CessionID)
		if err != nil {
			return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusBadRequest, "Invalid cession ID", err)
		}
	}

	recovery := &entity.ReinsuranceRecovery{
		RecoveryNumber:    utils.GenerateRecoveryNumber(),
		ClaimID:           claimID,
		TreatyID:          treatyID,
		TreatyLayerID:     treatyLayerID,
		CessionID:         cessionID,
		GrossClaimAmount:  req.GrossAmount,
		RecoverableAmount: req.RecoverableAmount,
		RecoveredAmount:   0,
		OutstandingAmount: req.RecoverableAmount,
		Status:            string(shared.RecoveryStatusNotified),
		WorkflowStatus:    string(shared.RecoveryWorkflowNotification),
		Notes:             req.Notes,
		CreatedBy:         createdBy,
	}

	created, err := s.recoveryRepo.Create(ctx, recovery)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to create recovery", err)
	}

	// Create initial workflow event
	workflowEvent := &entity.RecoveryWorkflowEvent{
		RecoveryID:  created.ID,
		FromStatus:  "",
		ToStatus:    string(shared.RecoveryStatusNotified),
		EventType:   "NOTIFICATION",
		Notes:       req.Notes,
		PerformedBy: createdBy,
	}
	s.workflowEventRepo.Create(ctx, workflowEvent)

	logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeReinsuranceRecovery), created.ID, string(shared.AuditActionCreate))

	resp := reinsuranceSchema.ToRecoveryResponse(created)
	return schema.NewServiceResponse(resp, http.StatusCreated, "Recovery created successfully")
}

func (s *recoveryServiceImpl) ApplyRecoveryForClaim(ctx context.Context, claimID uuid.UUID, req reinsuranceSchema.ApplyRecoveryForClaimRequest, createdBy uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse] {
	// Find all ACTIVE treaties
	activeTreaties, err := s.treatyRepo.ListActive(ctx, 1000, 0)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to list active treaties", err)
	}

	if len(activeTreaties) == 0 {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryResponse](http.StatusNotFound, "No active treaties found", fmt.Errorf("no active treaties"))
	}

	approvedAmount := req.ApprovedAmount
	var responses []reinsuranceSchema.RecoveryResponse

	for _, treaty := range activeTreaties {
		if treaty.TreatyType == string(shared.TreatyTypeQuotaShare) {
			// QUOTA_SHARE: recoverable = approvedAmount * totalParticipantShare / 100
			totalShare, err := s.participantRepo.GetTotalShareByTreaty(ctx, treaty.ID)
			if err != nil || totalShare <= 0 {
				continue
			}

			recoverable := int64(float64(approvedAmount) * totalShare / 100)
			if recoverable <= 0 {
				continue
			}

			recovery := &entity.ReinsuranceRecovery{
				RecoveryNumber:    utils.GenerateRecoveryNumber(),
				ClaimID:           claimID,
				TreatyID:          treaty.ID,
				GrossClaimAmount:  approvedAmount,
				RecoverableAmount: recoverable,
				RecoveredAmount:   0,
				OutstandingAmount: recoverable,
				Status:            string(shared.RecoveryStatusNotified),
				WorkflowStatus:    string(shared.RecoveryWorkflowNotification),
				CreatedBy:         createdBy,
			}

			created, err := s.recoveryRepo.Create(ctx, recovery)
			if err != nil {
				continue
			}

			// Create workflow event
			workflowEvent := &entity.RecoveryWorkflowEvent{
				RecoveryID:  created.ID,
				FromStatus:  "",
				ToStatus:    string(shared.RecoveryStatusNotified),
				EventType:   "NOTIFICATION",
				PerformedBy: createdBy,
			}
			s.workflowEventRepo.Create(ctx, workflowEvent)

			logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeReinsuranceRecovery), created.ID, string(shared.AuditActionCreate))
			responses = append(responses, reinsuranceSchema.ToRecoveryResponse(created))

		} else if treaty.TreatyType == string(shared.TreatyTypeXOL) {
			// XOL: For each layer (sorted by layer_number), calculate recovery
			layers, err := s.layerRepo.ListByTreaty(ctx, treaty.ID)
			if err != nil || len(layers) == 0 {
				continue
			}

			// Sort layers by layer_number ascending
			sort.Slice(layers, func(i, j int) bool {
				return layers[i].LayerNumber < layers[j].LayerNumber
			})

			for _, layer := range layers {
				excess := approvedAmount - layer.AttachmentPoint
				if excess <= 0 {
					continue
				}

				layerExposure := excess
				if layerExposure > layer.LayerLimit {
					layerExposure = layer.LayerLimit
				}

				recoverable := layerExposure - layer.DeductibleAmount
				if recoverable <= 0 {
					continue
				}

				// Check aggregate limit
				if layer.AggregateLimit != nil {
					remaining := *layer.AggregateLimit - layer.AggregateUsed
					if remaining <= 0 {
						continue
					}
					if recoverable > remaining {
						recoverable = remaining
					}
				}

				recovery := &entity.ReinsuranceRecovery{
					RecoveryNumber:    utils.GenerateRecoveryNumber(),
					ClaimID:           claimID,
					TreatyID:          treaty.ID,
					TreatyLayerID:     layer.ID,
					GrossClaimAmount:  approvedAmount,
					RecoverableAmount: recoverable,
					RecoveredAmount:   0,
					OutstandingAmount: recoverable,
					Status:            string(shared.RecoveryStatusNotified),
					WorkflowStatus:    string(shared.RecoveryWorkflowNotification),
					CreatedBy:         createdBy,
				}

				created, err := s.recoveryRepo.Create(ctx, recovery)
				if err != nil {
					continue
				}

				// Update aggregate_used on the layer
				newAggregateUsed := layer.AggregateUsed + recoverable
				s.layerRepo.UpdateAggregateUsed(ctx, layer.ID, newAggregateUsed)

				// Create workflow event
				workflowEvent := &entity.RecoveryWorkflowEvent{
					RecoveryID:  created.ID,
					FromStatus:  "",
					ToStatus:    string(shared.RecoveryStatusNotified),
					EventType:   "NOTIFICATION",
					PerformedBy: createdBy,
				}
				s.workflowEventRepo.Create(ctx, workflowEvent)

				logAudit(ctx, s.auditSvc, createdBy, string(shared.AuditEntityTypeReinsuranceRecovery), created.ID, string(shared.AuditActionCreate))
				responses = append(responses, reinsuranceSchema.ToRecoveryResponse(created))
			}
		}
	}

	if len(responses) == 0 {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryResponse](http.StatusNotFound, "No recoveries could be applied for this claim", fmt.Errorf("no applicable treaties or layers"))
	}

	return schema.NewServiceResponse(responses, http.StatusCreated, fmt.Sprintf("Applied %d recoveries for claim", len(responses)))
}

// transitionRecovery is a helper that handles common workflow transition logic.
func (s *recoveryServiceImpl) transitionRecovery(
	ctx context.Context,
	id uuid.UUID,
	validFromStatuses []string,
	toStatus string,
	workflowStatus string,
	eventType string,
	notes string,
	userID uuid.UUID,
) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	existing, err := s.recoveryRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusNotFound, "Recovery not found", err)
	}

	// Validate current status
	validTransition := false
	for _, valid := range validFromStatuses {
		if existing.Status == valid {
			validTransition = true
			break
		}
	}
	if !validTransition {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](
			http.StatusBadRequest,
			fmt.Sprintf("Cannot transition from %s to %s", existing.Status, toStatus),
			fmt.Errorf("invalid status transition: %s -> %s", existing.Status, toStatus),
		)
	}

	fromStatus := existing.Status

	updated, err := s.recoveryRepo.UpdateStatus(ctx, id, toStatus, workflowStatus)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to update recovery status", err)
	}

	// Create workflow event
	workflowEvent := &entity.RecoveryWorkflowEvent{
		RecoveryID:  id,
		FromStatus:  fromStatus,
		ToStatus:    toStatus,
		EventType:   eventType,
		Notes:       notes,
		PerformedBy: userID,
	}
	s.workflowEventRepo.Create(ctx, workflowEvent)

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeReinsuranceRecovery), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToRecoveryResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, fmt.Sprintf("Recovery %s successfully", eventType))
}

func (s *recoveryServiceImpl) AcknowledgeRecovery(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	return s.transitionRecovery(
		ctx, id,
		[]string{string(shared.RecoveryStatusNotified)},
		string(shared.RecoveryStatusAcknowledged),
		string(shared.RecoveryWorkflowAcknowledgment),
		"ACKNOWLEDGMENT",
		req.Notes,
		userID,
	)
}

func (s *recoveryServiceImpl) RequestInfo(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	return s.transitionRecovery(
		ctx, id,
		[]string{string(shared.RecoveryStatusAcknowledged)},
		string(shared.RecoveryStatusInfoRequested),
		string(shared.RecoveryWorkflowInfoRequest),
		"INFO_REQUEST",
		req.Notes,
		userID,
	)
}

func (s *recoveryServiceImpl) ApproveRecovery(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	return s.transitionRecovery(
		ctx, id,
		[]string{string(shared.RecoveryStatusAcknowledged), string(shared.RecoveryStatusInfoRequested)},
		string(shared.RecoveryStatusApproved),
		string(shared.RecoveryWorkflowApproval),
		"APPROVAL",
		req.Notes,
		userID,
	)
}

func (s *recoveryServiceImpl) RecordPayment(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecordPaymentRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	existing, err := s.recoveryRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusNotFound, "Recovery not found", err)
	}

	if existing.Status != string(shared.RecoveryStatusApproved) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](
			http.StatusBadRequest,
			"Only APPROVED recoveries can receive payments",
			fmt.Errorf("invalid status: %s", existing.Status),
		)
	}

	fromStatus := existing.Status
	newRecoveredAmount := existing.RecoveredAmount + req.Amount
	newOutstandingAmount := existing.RecoverableAmount - newRecoveredAmount
	if newOutstandingAmount < 0 {
		newOutstandingAmount = 0
	}

	// Update recovered and outstanding amounts
	updated, err := s.recoveryRepo.UpdateRecoveredAmount(ctx, id, newRecoveredAmount, newOutstandingAmount)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to record payment", err)
	}

	// Update status to PAID
	updated, err = s.recoveryRepo.UpdateStatus(ctx, id, string(shared.RecoveryStatusPaid), string(shared.RecoveryWorkflowPayment))
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to update recovery status", err)
	}

	// Create workflow event
	workflowEvent := &entity.RecoveryWorkflowEvent{
		RecoveryID:  id,
		FromStatus:  fromStatus,
		ToStatus:    string(shared.RecoveryStatusPaid),
		EventType:   "PAYMENT",
		Notes:       req.Notes,
		PerformedBy: userID,
	}
	s.workflowEventRepo.Create(ctx, workflowEvent)

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeReinsuranceRecovery), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToRecoveryResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Payment recorded successfully")
}

func (s *recoveryServiceImpl) WriteOffRecovery(ctx context.Context, id uuid.UUID, req reinsuranceSchema.RecoveryWorkflowRequest, userID uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryResponse] {
	existing, err := s.recoveryRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusNotFound, "Recovery not found", err)
	}

	if existing.Status == string(shared.RecoveryStatusPaid) {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](
			http.StatusBadRequest,
			"Cannot write off a PAID recovery",
			fmt.Errorf("invalid status: %s", existing.Status),
		)
	}

	fromStatus := existing.Status

	// Keep the current workflow status when writing off
	updated, err := s.recoveryRepo.UpdateStatus(ctx, id, string(shared.RecoveryStatusWrittenOff), existing.WorkflowStatus)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to write off recovery", err)
	}

	// Create workflow event
	workflowEvent := &entity.RecoveryWorkflowEvent{
		RecoveryID:  id,
		FromStatus:  fromStatus,
		ToStatus:    string(shared.RecoveryStatusWrittenOff),
		EventType:   "WRITE_OFF",
		Notes:       req.Notes,
		PerformedBy: userID,
	}
	s.workflowEventRepo.Create(ctx, workflowEvent)

	logAudit(ctx, s.auditSvc, userID, string(shared.AuditEntityTypeReinsuranceRecovery), id, string(shared.AuditActionStateChange))

	resp := reinsuranceSchema.ToRecoveryResponse(updated)
	return schema.NewServiceResponse(resp, http.StatusOK, "Recovery written off successfully")
}

func (s *recoveryServiceImpl) GetRecovery(ctx context.Context, id uuid.UUID) *schema.ServiceResponse[reinsuranceSchema.RecoveryDetailResponse] {
	recovery, err := s.recoveryRepo.GetByID(ctx, id)
	if err != nil {
		return schema.NewServiceErrorResponse[reinsuranceSchema.RecoveryDetailResponse](http.StatusNotFound, "Recovery not found", err)
	}

	events, _ := s.workflowEventRepo.ListByRecovery(ctx, id)
	eventResponses := make([]reinsuranceSchema.RecoveryWorkflowEventResponse, len(events))
	for i, e := range events {
		eventResponses[i] = reinsuranceSchema.ToRecoveryWorkflowEventResponse(e)
	}

	resp := reinsuranceSchema.RecoveryDetailResponse{
		RecoveryResponse: reinsuranceSchema.ToRecoveryResponse(recovery),
		WorkflowEvents:   eventResponses,
	}
	return schema.NewServiceResponse(resp, http.StatusOK, "Recovery retrieved successfully")
}

func (s *recoveryServiceImpl) GetWorkflowEvents(ctx context.Context, recoveryID uuid.UUID) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryWorkflowEventResponse] {
	events, err := s.workflowEventRepo.ListByRecovery(ctx, recoveryID)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryWorkflowEventResponse](http.StatusInternalServerError, "Failed to list workflow events", err)
	}

	responses := make([]reinsuranceSchema.RecoveryWorkflowEventResponse, len(events))
	for i, e := range events {
		responses[i] = reinsuranceSchema.ToRecoveryWorkflowEventResponse(e)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Workflow events retrieved successfully")
}

func (s *recoveryServiceImpl) ListByClaim(ctx context.Context, claimID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse] {
	offset := (page - 1) * pageSize
	recoveries, err := s.recoveryRepo.ListByClaim(ctx, claimID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to list recoveries by claim", err)
	}

	responses := make([]reinsuranceSchema.RecoveryResponse, len(recoveries))
	for i, r := range recoveries {
		responses[i] = reinsuranceSchema.ToRecoveryResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Recoveries retrieved successfully")
}

func (s *recoveryServiceImpl) ListByTreaty(ctx context.Context, treatyID uuid.UUID, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse] {
	offset := (page - 1) * pageSize
	recoveries, err := s.recoveryRepo.ListByTreaty(ctx, treatyID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to list recoveries by treaty", err)
	}

	responses := make([]reinsuranceSchema.RecoveryResponse, len(recoveries))
	for i, r := range recoveries {
		responses[i] = reinsuranceSchema.ToRecoveryResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Recoveries retrieved successfully")
}

func (s *recoveryServiceImpl) ListOutstanding(ctx context.Context, page, pageSize int) *schema.ServiceResponse[[]reinsuranceSchema.RecoveryResponse] {
	offset := (page - 1) * pageSize
	recoveries, err := s.recoveryRepo.ListOutstanding(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.RecoveryResponse](http.StatusInternalServerError, "Failed to list outstanding recoveries", err)
	}

	responses := make([]reinsuranceSchema.RecoveryResponse, len(recoveries))
	for i, r := range recoveries {
		responses[i] = reinsuranceSchema.ToRecoveryResponse(r)
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Outstanding recoveries retrieved successfully")
}

func (s *recoveryServiceImpl) GetRecoveryCount(ctx context.Context) *schema.ServiceResponse[int64] {
	count, err := s.recoveryRepo.Count(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to count recoveries", err)
	}
	return schema.NewServiceResponse(count, http.StatusOK, "Recovery count retrieved")
}

func (s *recoveryServiceImpl) GetAgedAnalysis(ctx context.Context) *schema.ServiceResponse[[]reinsuranceSchema.AgedRecoveryBucketResponse] {
	buckets, err := s.recoveryRepo.GetAgedAnalysis(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[[]reinsuranceSchema.AgedRecoveryBucketResponse](http.StatusInternalServerError, "Failed to get aged analysis", err)
	}

	responses := make([]reinsuranceSchema.AgedRecoveryBucketResponse, len(buckets))
	for i, b := range buckets {
		responses[i] = reinsuranceSchema.AgedRecoveryBucketResponse{
			Bucket:           b.Bucket,
			Count:            b.Count,
			TotalOutstanding: b.TotalOutstanding,
		}
	}
	return schema.NewServiceResponse(responses, http.StatusOK, "Aged analysis retrieved successfully")
}

func (s *recoveryServiceImpl) GetTotalRecoverableAmount(ctx context.Context) *schema.ServiceResponse[int64] {
	total, err := s.recoveryRepo.GetTotalRecoverableAmountAll(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get total recoverable amount", err)
	}
	return schema.NewServiceResponse(total, http.StatusOK, "Total recoverable amount retrieved")
}

func (s *recoveryServiceImpl) GetTotalRecoveredAmount(ctx context.Context) *schema.ServiceResponse[int64] {
	total, err := s.recoveryRepo.GetTotalRecoveredAmountAll(ctx)
	if err != nil {
		return schema.NewServiceErrorResponse[int64](http.StatusInternalServerError, "Failed to get total recovered amount", err)
	}
	return schema.NewServiceResponse(total, http.StatusOK, "Total recovered amount retrieved")
}

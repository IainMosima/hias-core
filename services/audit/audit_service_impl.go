package audit

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bitbiz/hias-core/domains/audit/entity"
	"github.com/bitbiz/hias-core/domains/audit/repository"
	auditSchema "github.com/bitbiz/hias-core/domains/audit/schema"
	"github.com/bitbiz/hias-core/domains/audit/service"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type auditServiceImpl struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(
	auditRepo repository.AuditRepository,
) service.AuditService {
	return &auditServiceImpl{
		auditRepo: auditRepo,
	}
}

func (s *auditServiceImpl) LogEvent(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string, oldValue, newValue json.RawMessage, ipAddress, userAgent string) *schema.ServiceResponse[string] {
	event := &entity.AuditEvent{
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValue:   oldValue,
		NewValue:   newValue,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	created, err := s.auditRepo.Create(ctx, event)
	if err != nil {
		return schema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to log audit event", err)
	}

	return schema.NewServiceResponse(created.ID.String(), http.StatusCreated, "Audit event logged")
}

func (s *auditServiceImpl) ListEvents(ctx context.Context, page, pageSize int) *schema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	events, err := s.auditRepo.List(ctx, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list audit events", err)
	}

	responses := auditSchema.ToAuditEventResponseList(events)
	return schema.NewServiceResponse[interface{}](responses, http.StatusOK, "Audit events retrieved")
}

func (s *auditServiceImpl) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	events, err := s.auditRepo.ListByEntity(ctx, entityType, entityID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list audit events by entity", err)
	}

	responses := auditSchema.ToAuditEventResponseList(events)
	return schema.NewServiceResponse[interface{}](responses, http.StatusOK, "Audit events retrieved")
}

func (s *auditServiceImpl) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	events, err := s.auditRepo.ListByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return schema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list audit events by user", err)
	}

	responses := auditSchema.ToAuditEventResponseList(events)
	return schema.NewServiceResponse[interface{}](responses, http.StatusOK, "Audit events retrieved")
}

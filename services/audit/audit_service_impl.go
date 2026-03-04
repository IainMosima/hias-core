package audit

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bitbiz/hias-core/domains/audit/entity"
	"github.com/bitbiz/hias-core/domains/audit/repository"
	"github.com/bitbiz/hias-core/domains/audit/schema"
	authSchema "github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type auditServiceImpl struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) *auditServiceImpl {
	return &auditServiceImpl{auditRepo: auditRepo}
}

func (s *auditServiceImpl) LogEvent(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string, oldValue, newValue json.RawMessage, ipAddress, userAgent string) *authSchema.ServiceResponse[string] {
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

	_, err := s.auditRepo.Create(ctx, event)
	if err != nil {
		log.Printf("Failed to log audit event: %v", err)
		return authSchema.NewServiceErrorResponse[string](http.StatusInternalServerError, "Failed to log audit event", err)
	}

	return authSchema.NewServiceResponse("Event logged", http.StatusCreated, "Audit event created")
}

func (s *auditServiceImpl) ListEvents(ctx context.Context, page, pageSize int) *authSchema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	events, err := s.auditRepo.List(ctx, pageSize, offset)
	if err != nil {
		return authSchema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list audit events", err)
	}

	count, _ := s.auditRepo.Count(ctx)

	result := map[string]interface{}{
		"events":      toAuditResponseList(events),
		"total_count": count,
		"page":        page,
		"page_size":   pageSize,
	}

	return authSchema.NewServiceResponse[interface{}](result, http.StatusOK, "Audit events retrieved")
}

func (s *auditServiceImpl) ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID, page, pageSize int) *authSchema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	events, err := s.auditRepo.ListByEntity(ctx, entityType, entityID, pageSize, offset)
	if err != nil {
		return authSchema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list audit events by entity", err)
	}

	count, _ := s.auditRepo.CountByEntity(ctx, entityType, entityID)

	result := map[string]interface{}{
		"events":      toAuditResponseList(events),
		"total_count": count,
		"page":        page,
		"page_size":   pageSize,
	}

	return authSchema.NewServiceResponse[interface{}](result, http.StatusOK, "Audit events retrieved")
}

func (s *auditServiceImpl) ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) *authSchema.ServiceResponse[interface{}] {
	offset := (page - 1) * pageSize
	events, err := s.auditRepo.ListByUser(ctx, userID, pageSize, offset)
	if err != nil {
		return authSchema.NewServiceErrorResponse[interface{}](http.StatusInternalServerError, "Failed to list audit events by user", err)
	}

	result := map[string]interface{}{
		"events":    toAuditResponseList(events),
		"page":      page,
		"page_size": pageSize,
	}

	return authSchema.NewServiceResponse[interface{}](result, http.StatusOK, "Audit events retrieved")
}

func toAuditResponseList(events []*entity.AuditEvent) []schema.AuditEventResponse {
	responses := make([]schema.AuditEventResponse, len(events))
	for i, e := range events {
		responses[i] = schema.ToAuditEventResponse(e)
	}
	return responses
}

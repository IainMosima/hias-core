package service

import (
	"context"
	"encoding/json"
	"github.com/bitbiz/hias-core/domains/identity/schema"
	"github.com/google/uuid"
)

type AuditService interface {
	LogEvent(ctx context.Context, userID uuid.UUID, entityType string, entityID uuid.UUID, action string, oldValue, newValue json.RawMessage, ipAddress, userAgent string) *schema.ServiceResponse[string]
	ListEvents(ctx context.Context, page, pageSize int) *schema.ServiceResponse[interface{}]
	ListByEntity(ctx context.Context, entityType string, entityID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}]
	ListByUser(ctx context.Context, userID uuid.UUID, page, pageSize int) *schema.ServiceResponse[interface{}]
}

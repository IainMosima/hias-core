package schema

import (
	"encoding/json"
	"time"

	"github.com/bitbiz/hias-core/domains/audit/entity"
	"github.com/google/uuid"
)

type AuditEventResponse struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	EntityType string          `json:"entity_type"`
	EntityID   uuid.UUID       `json:"entity_id"`
	Action     string          `json:"action"`
	OldValue   json.RawMessage `json:"old_value,omitempty"`
	NewValue   json.RawMessage `json:"new_value,omitempty"`
	IPAddress  string          `json:"ip_address"`
	UserAgent  string          `json:"user_agent"`
	CreatedAt  time.Time       `json:"created_at"`
}

func ToAuditEventResponse(event *entity.AuditEvent) AuditEventResponse {
	return AuditEventResponse{
		ID:         event.ID,
		UserID:     event.UserID,
		EntityType: event.EntityType,
		EntityID:   event.EntityID,
		Action:     event.Action,
		OldValue:   event.OldValue,
		NewValue:   event.NewValue,
		IPAddress:  event.IPAddress,
		UserAgent:  event.UserAgent,
		CreatedAt:  event.CreatedAt,
	}
}

func ToAuditEventResponseList(events []*entity.AuditEvent) []AuditEventResponse {
	responses := make([]AuditEventResponse, len(events))
	for i, event := range events {
		responses[i] = ToAuditEventResponse(event)
	}
	return responses
}

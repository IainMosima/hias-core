package schema

import (
	"encoding/json"
	"time"

	"github.com/bitbiz/hias-core/domains/notification/entity"
	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	Channel    string          `json:"channel"`
	Type       string          `json:"type"`
	Subject    string          `json:"subject"`
	Body       string          `json:"body"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	Status     string          `json:"status"`
	RetryCount int             `json:"retry_count"`
	MaxRetries int             `json:"max_retries"`
	SentAt     *time.Time      `json:"sent_at,omitempty"`
	ReadAt     *time.Time      `json:"read_at,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

func ToNotificationResponse(n *entity.Notification) NotificationResponse {
	return NotificationResponse{
		ID:         n.ID,
		UserID:     n.UserID,
		Channel:    n.Channel,
		Type:       n.Type,
		Subject:    n.Subject,
		Body:       n.Body,
		Metadata:   n.Metadata,
		Status:     n.Status,
		RetryCount: n.RetryCount,
		MaxRetries: n.MaxRetries,
		SentAt:     n.SentAt,
		ReadAt:     n.ReadAt,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

func ToNotificationResponseList(notifications []*entity.Notification) []NotificationResponse {
	responses := make([]NotificationResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = ToNotificationResponse(n)
	}
	return responses
}

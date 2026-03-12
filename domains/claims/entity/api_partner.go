package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type APIPartner struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	PartnerType        string          `json:"partner_type"`
	APIKey             string          `json:"api_key"`
	APISecretHash      string          `json:"-"`
	ProviderID         uuid.UUID       `json:"provider_id,omitempty"`
	IsActive           bool            `json:"is_active"`
	RateLimitPerMinute int             `json:"rate_limit_per_minute"`
	AllowedClaimTypes  []string        `json:"allowed_claim_types"`
	WebhookURL         string          `json:"webhook_url,omitempty"`
	ContactEmail       string          `json:"contact_email,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

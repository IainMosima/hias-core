package schema

import (
	"encoding/json"

	"github.com/bitbiz/hias-core/domains/claims/entity"
	"github.com/google/uuid"
)

type CreateAPIPartnerRequest struct {
	Name               string          `json:"name" binding:"required"`
	PartnerType        string          `json:"partner_type" binding:"required"`
	ProviderID         string          `json:"provider_id"`
	RateLimitPerMinute int             `json:"rate_limit_per_minute"`
	AllowedClaimTypes  []string        `json:"allowed_claim_types"`
	WebhookURL         string          `json:"webhook_url"`
	ContactEmail       string          `json:"contact_email"`
	Metadata           json.RawMessage `json:"metadata"`
}

type UpdateAPIPartnerRequest struct {
	Name               string          `json:"name" binding:"required"`
	PartnerType        string          `json:"partner_type" binding:"required"`
	ProviderID         string          `json:"provider_id"`
	RateLimitPerMinute int             `json:"rate_limit_per_minute"`
	AllowedClaimTypes  []string        `json:"allowed_claim_types"`
	WebhookURL         string          `json:"webhook_url"`
	ContactEmail       string          `json:"contact_email"`
	Metadata           json.RawMessage `json:"metadata"`
}

type APIPartnerResponse struct {
	ID                 uuid.UUID       `json:"id"`
	Name               string          `json:"name"`
	PartnerType        string          `json:"partner_type"`
	APIKey             string          `json:"api_key"`
	ProviderID         uuid.UUID       `json:"provider_id,omitempty"`
	IsActive           bool            `json:"is_active"`
	RateLimitPerMinute int             `json:"rate_limit_per_minute"`
	AllowedClaimTypes  []string        `json:"allowed_claim_types"`
	WebhookURL         string          `json:"webhook_url,omitempty"`
	ContactEmail       string          `json:"contact_email,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	CreatedAt          string          `json:"created_at"`
	UpdatedAt          string          `json:"updated_at"`
}

type CreateAPIPartnerResponse struct {
	APIPartnerResponse
	APISecret string `json:"api_secret"`
}

func ToAPIPartnerResponse(p *entity.APIPartner) APIPartnerResponse {
	return APIPartnerResponse{
		ID:                 p.ID,
		Name:               p.Name,
		PartnerType:        p.PartnerType,
		APIKey:             p.APIKey,
		ProviderID:         p.ProviderID,
		IsActive:           p.IsActive,
		RateLimitPerMinute: p.RateLimitPerMinute,
		AllowedClaimTypes:  p.AllowedClaimTypes,
		WebhookURL:         p.WebhookURL,
		ContactEmail:       p.ContactEmail,
		Metadata:           p.Metadata,
		CreatedAt:          p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:          p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

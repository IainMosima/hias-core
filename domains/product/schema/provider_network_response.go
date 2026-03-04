package schema

import (
	"time"

	"github.com/bitbiz/hias-core/domains/product/entity"
	"github.com/google/uuid"
)

type ProviderNetworkResponse struct {
	ID              uuid.UUID `json:"id"`
	PlanID          uuid.UUID `json:"plan_id"`
	ProviderID      uuid.UUID `json:"provider_id"`
	BenefitCategory string    `json:"benefit_category,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToProviderNetworkResponse(n *entity.ProviderNetwork) ProviderNetworkResponse {
	return ProviderNetworkResponse{
		ID: n.ID, PlanID: n.PlanID, ProviderID: n.ProviderID,
		BenefitCategory: n.BenefitCategory, Status: n.Status,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}
}

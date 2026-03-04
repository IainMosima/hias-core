package entity

import (
	"time"

	"github.com/google/uuid"
)

type ApprovalLimit struct {
	ID                    uuid.UUID `json:"id"`
	RoleName              string    `json:"role_name"`
	MaxDiscountPercentage int64     `json:"max_discount_percentage"`
	MaxDiscountAmount     int64     `json:"max_discount_amount"`
	MaxLoadingPercentage  int64     `json:"max_loading_percentage"`
	MaxLoadingAmount      int64     `json:"max_loading_amount"`
	EscalationRole        string    `json:"escalation_role"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

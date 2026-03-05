package entity

import (
	"time"

	"github.com/google/uuid"
)

type Bordereau struct {
	ID              uuid.UUID `json:"id"`
	BordereauNumber string    `json:"bordereau_number"`
	TreatyID        uuid.UUID `json:"treaty_id"`
	BordereauType   string    `json:"bordereau_type"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	TotalGross      int64     `json:"total_gross"`
	TotalCeded      int64     `json:"total_ceded"`
	TotalCommission int64     `json:"total_commission"`
	ItemCount       int       `json:"item_count"`
	Status          string    `json:"status"`
	CreatedBy       uuid.UUID `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

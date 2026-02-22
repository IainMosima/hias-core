package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	CognitoSub string    `json:"cognito_sub"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	NationalID string    `json:"national_id"`
	RoleID     uuid.UUID `json:"role_id"`
	RoleName   string    `json:"role_name,omitempty"`
	Status     string    `json:"status"`
	CreatedBy  uuid.UUID `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

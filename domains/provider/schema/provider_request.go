package schema

import "time"

type RegisterProviderRequest struct {
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=hospital clinic pharmacy lab"`
	Tier          string `json:"tier"`
	LicenseNumber string `json:"license_number" binding:"required"`
	County        string `json:"county"`
	Address       string `json:"address"`
	Phone         string `json:"phone" binding:"required"`
	Email         string `json:"email" binding:"required,email"`
	ContactPerson string `json:"contact_person"`
}

type UpdateProviderRequest struct {
	Name          *string `json:"name"`
	County        *string `json:"county"`
	Address       *string `json:"address"`
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	ContactPerson *string `json:"contact_person"`
}

type CreateContractRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
	Terms     string    `json:"terms"`
}

type CreateRateCardRequest struct {
	ProcedureCode string    `json:"procedure_code" binding:"required"`
	ProcedureName string    `json:"procedure_name" binding:"required"`
	RateAmount    int64     `json:"rate_amount" binding:"required,min=1"`
	EffectiveDate time.Time `json:"effective_date"`
	AgeFrom       int       `json:"age_from"`
	AgeTo         int       `json:"age_to"`
	Gender        string    `json:"gender"`
	Relationship  string    `json:"relationship"`
}

type BulkCreateRateCardRequest struct {
	RateCards []CreateRateCardRequest `json:"rate_cards" binding:"required,min=1"`
}

type UpdateAccreditationRequest struct {
	AccreditationStatus string `json:"accreditation_status" binding:"required,oneof=NONE PENDING ACCREDITED EXPIRED REVOKED"`
	AccreditationExpiry string `json:"accreditation_expiry"` // YYYY-MM-DD
	AccreditationBody   string `json:"accreditation_body"`
}

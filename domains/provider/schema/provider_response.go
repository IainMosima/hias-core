package schema

import (
	"github.com/bitbiz/hias-core/domains/provider/entity"
	"github.com/google/uuid"
	"time"
)

type ProviderResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	LicenseNumber string    `json:"license_number"`
	Status        string    `json:"status"`
	Tier          string    `json:"tier"`
	County        string    `json:"county"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	ContactPerson string    `json:"contact_person"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ContractResponse struct {
	ID         uuid.UUID `json:"id"`
	ProviderID uuid.UUID `json:"provider_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Terms      string    `json:"terms"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type RateCardResponse struct {
	ID            uuid.UUID `json:"id"`
	ProviderID    uuid.UUID `json:"provider_id"`
	ProcedureCode string    `json:"procedure_code"`
	ProcedureName string    `json:"procedure_name"`
	RateAmount    int64     `json:"rate_amount"`
	EffectiveDate time.Time `json:"effective_date"`
	AgeFrom       int       `json:"age_from"`
	AgeTo         int       `json:"age_to"`
	Gender        string    `json:"gender,omitempty"`
	Relationship  string    `json:"relationship,omitempty"`
}

func ToProviderResponse(p *entity.Provider) ProviderResponse {
	return ProviderResponse{
		ID: p.ID, Name: p.Name, Type: p.Type, LicenseNumber: p.LicenseNumber,
		Status: p.Status, Tier: p.Tier, County: p.County, Phone: p.Phone, Email: p.Email,
		ContactPerson: p.ContactPerson, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func ToContractResponse(c *entity.Contract) ContractResponse {
	return ContractResponse{
		ID: c.ID, ProviderID: c.ProviderID, StartDate: c.StartDate,
		EndDate: c.EndDate, Terms: c.Terms, Status: c.Status, CreatedAt: c.CreatedAt,
	}
}

func ToRateCardResponse(r *entity.RateCard) RateCardResponse {
	return RateCardResponse{
		ID: r.ID, ProviderID: r.ProviderID, ProcedureCode: r.ProcedureCode,
		ProcedureName: r.ProcedureName, RateAmount: r.RateAmount,
		EffectiveDate: r.EffectiveDate, AgeFrom: r.AgeFrom, AgeTo: r.AgeTo,
		Gender: r.Gender, Relationship: r.Relationship,
	}
}

package schema

type CreateProviderNetworkRequest struct {
	ProviderID      string `json:"provider_id" binding:"required,uuid"`
	BenefitCategory string `json:"benefit_category"`
}

type UpdateProviderNetworkStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

package bank

type TransferRequest struct {
	AccountNumber    string `json:"account_number"`
	BankCode         string `json:"bank_code"`
	BeneficiaryName  string `json:"beneficiary_name"`
	Amount           int64  `json:"amount"`
	Currency         string `json:"currency"`
	Reference        string `json:"reference"`
	Narration        string `json:"narration"`
}

type TransferResponse struct {
	TransactionID string `json:"transaction_id"`
	Reference     string `json:"reference"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

type StatusResponse struct {
	TransactionID string `json:"transaction_id"`
	Reference     string `json:"reference"`
	Status        string `json:"status"` // PENDING, COMPLETED, FAILED
	Message       string `json:"message"`
}

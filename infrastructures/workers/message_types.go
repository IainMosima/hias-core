package workers

import "time"

// ExtractionResultMessage is received from the AI service after document processing.
type ExtractionResultMessage struct {
	DocumentID     string                 `json:"document_id"`
	ClaimID        string                 `json:"claim_id,omitempty"`
	Status         string                 `json:"status"` // completed, failed
	ExtractedData  map[string]interface{} `json:"extracted_data,omitempty"`
	Confidence     float64                `json:"confidence"`
	Error          string                 `json:"error,omitempty"`
	ProcessedAt    time.Time              `json:"processed_at"`
}

// ClaimSubmittedMessage is published when a new claim is submitted.
type ClaimSubmittedMessage struct {
	ClaimID     string    `json:"claim_id"`
	ClaimNumber string    `json:"claim_number"`
	PolicyID    string    `json:"policy_id"`
	MemberID    string    `json:"member_id"`
	ProviderID  string    `json:"provider_id"`
	TotalAmount int64     `json:"total_amount"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// ClaimApprovedMessage is published when a claim is approved.
type ClaimApprovedMessage struct {
	ClaimID        string    `json:"claim_id"`
	ClaimNumber    string    `json:"claim_number"`
	ProviderID     string    `json:"provider_id"`
	ApprovedAmount int64     `json:"approved_amount"`
	ApprovedAt     time.Time `json:"approved_at"`
}

// PaymentWebhookMessage is received from M-Pesa/bank webhook.
type PaymentWebhookMessage struct {
	PaymentID        string                 `json:"payment_id"`
	TransactionID    string                 `json:"transaction_id"`
	Status           string                 `json:"status"`
	Amount           int64                  `json:"amount"`
	GatewayResponse  map[string]interface{} `json:"gateway_response"`
	ProcessedAt      time.Time              `json:"processed_at"`
}

// PreAuthSubmittedMessage is published when a pre-auth is submitted.
type PreAuthSubmittedMessage struct {
	PreAuthID  string    `json:"preauth_id"`
	PolicyID   string    `json:"policy_id"`
	MemberID   string    `json:"member_id"`
	ProviderID string    `json:"provider_id"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// NotificationMessage is consumed by the notification handler.
type NotificationMessage struct {
	UserID  string `json:"user_id"`
	Channel string `json:"channel"` // SMS, EMAIL, IN_APP, PUSH
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

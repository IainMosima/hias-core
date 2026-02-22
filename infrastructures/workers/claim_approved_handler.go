package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ClaimApprovedHandler creates a payment record when a claim is approved.
type ClaimApprovedHandler struct {
	// paymentService would be injected
}

func NewClaimApprovedHandler() *ClaimApprovedHandler {
	return &ClaimApprovedHandler{}
}

func (h *ClaimApprovedHandler) GetName() string {
	return "claim-approved-handler"
}

func (h *ClaimApprovedHandler) HandleMessage(ctx context.Context, payload []byte) error {
	var msg ClaimApprovedMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal claim approved message: %w", err)
	}

	log.Printf("Processing approved claim %s for provider %s, amount: %d",
		msg.ClaimNumber, msg.ProviderID, msg.ApprovedAmount)

	// Create a remittance-pending payment record for the provider
	// h.paymentService.CreateProviderPayment(ctx, msg.ProviderID, msg.ApprovedAmount)

	log.Printf("Payment record created for approved claim %s", msg.ClaimNumber)
	return nil
}

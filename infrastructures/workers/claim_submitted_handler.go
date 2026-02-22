package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ClaimSubmittedHandler processes newly submitted claims through validation and adjudication.
type ClaimSubmittedHandler struct {
	// claimService, validatorService, adjudicatorService would be injected
}

func NewClaimSubmittedHandler() *ClaimSubmittedHandler {
	return &ClaimSubmittedHandler{}
}

func (h *ClaimSubmittedHandler) GetName() string {
	return "claim-submitted-handler"
}

func (h *ClaimSubmittedHandler) HandleMessage(ctx context.Context, payload []byte) error {
	var msg ClaimSubmittedMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal claim submitted message: %w", err)
	}

	log.Printf("Processing submitted claim %s (amount: %d)", msg.ClaimNumber, msg.TotalAmount)

	// Step 1: Validate the claim
	// valid, errors, err := h.validatorService.ValidateClaim(ctx, claim, lineItems)

	// Step 2: If valid, run adjudication
	// result, err := h.adjudicatorService.Adjudicate(ctx, claim, lineItems)

	// Step 3: Update claim status based on adjudication result

	log.Printf("Claim %s processed successfully", msg.ClaimNumber)
	return nil
}

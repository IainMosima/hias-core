package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ExtractionResultHandler processes document extraction results from the AI service.
type ExtractionResultHandler struct {
	// claimService would be injected for auto-creating claims
}

func NewExtractionResultHandler() *ExtractionResultHandler {
	return &ExtractionResultHandler{}
}

func (h *ExtractionResultHandler) GetName() string {
	return "extraction-result-handler"
}

func (h *ExtractionResultHandler) HandleMessage(ctx context.Context, payload []byte) error {
	var msg ExtractionResultMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal extraction result: %w", err)
	}

	log.Printf("Processing extraction result for document %s, status: %s, confidence: %.2f",
		msg.DocumentID, msg.Status, msg.Confidence)

	if msg.Status == "failed" {
		log.Printf("Document extraction failed for %s: %s", msg.DocumentID, msg.Error)
		// Queue for manual review
		return nil
	}

	// High confidence → auto-create claim
	if msg.Confidence >= 0.85 {
		log.Printf("High confidence extraction (%.2f) for document %s — auto-creating claim", msg.Confidence, msg.DocumentID)
		// TODO: call claimService.SubmitClaim with extracted data
		return nil
	}

	// Low confidence → queue for human review
	log.Printf("Low confidence extraction (%.2f) for document %s — queued for review", msg.Confidence, msg.DocumentID)
	return nil
}

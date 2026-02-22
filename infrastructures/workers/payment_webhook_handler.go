package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// PaymentWebhookHandler processes M-Pesa/bank payment webhook confirmations.
type PaymentWebhookHandler struct {
	// paymentService would be injected
}

func NewPaymentWebhookHandler() *PaymentWebhookHandler {
	return &PaymentWebhookHandler{}
}

func (h *PaymentWebhookHandler) GetName() string {
	return "payment-webhook-handler"
}

func (h *PaymentWebhookHandler) HandleMessage(ctx context.Context, payload []byte) error {
	var msg PaymentWebhookMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal payment webhook message: %w", err)
	}

	log.Printf("Processing payment webhook for payment %s, tx: %s, status: %s",
		msg.PaymentID, msg.TransactionID, msg.Status)

	switch msg.Status {
	case "CONFIRMED":
		// h.paymentService.ConfirmPayment(ctx, msg.PaymentID, msg.TransactionID)
		log.Printf("Payment %s confirmed", msg.PaymentID)
	case "FAILED":
		// h.paymentService.FailPayment(ctx, msg.PaymentID, msg.GatewayResponse)
		log.Printf("Payment %s failed", msg.PaymentID)
	default:
		log.Printf("Unknown payment status: %s for payment %s", msg.Status, msg.PaymentID)
	}

	return nil
}

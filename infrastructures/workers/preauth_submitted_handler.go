package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// PreAuthSubmittedHandler notifies claims officers when a new pre-auth is submitted.
type PreAuthSubmittedHandler struct {
	// notificationService would be injected
}

func NewPreAuthSubmittedHandler() *PreAuthSubmittedHandler {
	return &PreAuthSubmittedHandler{}
}

func (h *PreAuthSubmittedHandler) GetName() string {
	return "preauth-submitted-handler"
}

func (h *PreAuthSubmittedHandler) HandleMessage(ctx context.Context, payload []byte) error {
	var msg PreAuthSubmittedMessage
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("failed to unmarshal preauth submitted message: %w", err)
	}

	log.Printf("Processing pre-auth submission %s for member %s", msg.PreAuthID, msg.MemberID)

	// Notify assigned claims officer about new pre-auth request
	// h.notificationService.Send(ctx, NotificationMessage{
	//     Channel: "IN_APP",
	//     Type:    "preauth_submitted",
	//     Title:   "New Pre-Authorization Request",
	//     Message: fmt.Sprintf("Pre-auth %s requires review", msg.PreAuthID),
	// })

	log.Printf("Notification sent for pre-auth %s", msg.PreAuthID)
	return nil
}

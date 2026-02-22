package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SMSAdapter struct {
	apiKey   string
	username string
	senderID string
	baseURL  string
	client   *http.Client
}

func NewSMSAdapter(apiKey, username, senderID string) *SMSAdapter {
	return &SMSAdapter{
		apiKey:   apiKey,
		username: username,
		senderID: senderID,
		baseURL:  "https://api.africastalking.com/version1",
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *SMSAdapter) Send(ctx context.Context, phone, message string) error {
	payload := map[string]string{
		"username": a.username,
		"to":       phone,
		"message":  message,
		"from":     a.senderID,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.baseURL+"/messaging", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	req.Header.Set("apiKey", a.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("SMS send failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS API returned status %d", resp.StatusCode)
	}

	return nil
}

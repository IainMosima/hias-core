package iprs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Adapter struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewAdapter(baseURL, apiKey string) *Adapter {
	return &Adapter{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *Adapter) VerifyNationalID(ctx context.Context, nationalID string) (*VerifyResponse, error) {
	url := fmt.Sprintf("%s/api/v1/verify/%s", a.baseURL, nationalID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPRS request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("IPRS verification failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPRS returned status %d", resp.StatusCode)
	}

	var verifyResp VerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&verifyResp); err != nil {
		return nil, fmt.Errorf("failed to decode IPRS response: %w", err)
	}

	verifyResp.Verified = true
	return &verifyResp, nil
}

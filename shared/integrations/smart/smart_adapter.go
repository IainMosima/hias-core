package smart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Adapter struct {
	baseURL      string
	apiKey       string
	apiSecret    string
	facilityCode string
	httpClient   *http.Client
}

func NewAdapter(baseURL, apiKey, apiSecret, facilityCode string) *Adapter {
	return &Adapter{
		baseURL:      baseURL,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		facilityCode: facilityCode,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *Adapter) SubmitClaim(ctx context.Context, claim ClaimSubmission) (*ClaimSubmissionResponse, error) {
	claim.FacilityCode = a.facilityCode

	body, err := json.Marshal(claim)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal claim: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/claims", a.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create SMART request: %w", err)
	}

	req.Header.Set("X-API-Key", a.apiKey)
	req.Header.Set("X-API-Secret", a.apiSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SMART claim submission failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("SMART returned status %d", resp.StatusCode)
	}

	var result ClaimSubmissionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode SMART response: %w", err)
	}

	return &result, nil
}

func (a *Adapter) GetClaimStatus(ctx context.Context, referenceNumber string) (*ClaimStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/claims/%s/status", a.baseURL, referenceNumber)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create SMART status request: %w", err)
	}

	req.Header.Set("X-API-Key", a.apiKey)
	req.Header.Set("X-API-Secret", a.apiSecret)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SMART status query failed: %w", err)
	}
	defer resp.Body.Close()

	var result ClaimStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode SMART status response: %w", err)
	}

	return &result, nil
}

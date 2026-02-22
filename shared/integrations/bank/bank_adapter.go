package bank

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Adapter struct {
	baseURL    string
	apiKey     string
	accountNo  string
	httpClient *http.Client
}

func NewAdapter(baseURL, apiKey, accountNo string) *Adapter {
	return &Adapter{
		baseURL:    baseURL,
		apiKey:     apiKey,
		accountNo:  accountNo,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *Adapter) InitiateTransfer(ctx context.Context, req TransferRequest) (*TransferResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/transfers", a.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create bank transfer request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("bank transfer request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("bank API returned status %d", resp.StatusCode)
	}

	var result TransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode bank response: %w", err)
	}

	return &result, nil
}

func (a *Adapter) CheckStatus(ctx context.Context, transactionID string) (*StatusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/transfers/%s/status", a.baseURL, transactionID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bank status query failed: %w", err)
	}
	defer resp.Body.Close()

	var result StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode bank status response: %w", err)
	}

	return &result, nil
}

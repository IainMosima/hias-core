package mpesa

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bitbiz/hias-core/shared/utils"
)

type Adapter struct {
	config     Config
	httpClient *http.Client
	token      string
	tokenExpiry time.Time
}

func NewAdapter(config Config) *Adapter {
	return &Adapter{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (a *Adapter) getAccessToken(ctx context.Context) (string, error) {
	if a.token != "" && time.Now().Before(a.tokenExpiry) {
		return a.token, nil
	}

	url := fmt.Sprintf("%s/oauth/v1/generate?grant_type=client_credentials", a.config.BaseURL())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(a.config.ConsumerKey + ":" + a.config.ConsumerSecret))
	req.Header.Set("Authorization", "Basic "+credentials)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp AuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	a.token = tokenResp.AccessToken
	a.tokenExpiry = time.Now().Add(50 * time.Minute) // tokens expire in ~1h
	return a.token, nil
}

func (a *Adapter) STKPush(ctx context.Context, req STKPushRequest) (*STKPushResponse, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")
	password := base64.StdEncoding.EncodeToString([]byte(a.config.Shortcode + a.config.Passkey + timestamp))

	payload := map[string]interface{}{
		"BusinessShortCode": a.config.Shortcode,
		"Password":          password,
		"Timestamp":         timestamp,
		"TransactionType":   "CustomerPayBillOnline",
		"Amount":            req.Amount / 100, // convert from cents
		"PartyA":            req.PhoneNumber,
		"PartyB":            a.config.Shortcode,
		"PhoneNumber":       req.PhoneNumber,
		"CallBackURL":       a.config.CallbackURL,
		"AccountReference":  req.AccountRef,
		"TransactionDesc":   req.Description,
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/mpesa/stkpush/v1/processrequest", a.config.BaseURL())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create STK push request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("STK push request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	utils.LogInfo("M-Pesa STK Push response: %s", string(respBody))

	var stkResp STKPushResponse
	if err := json.Unmarshal(respBody, &stkResp); err != nil {
		return nil, fmt.Errorf("failed to decode STK push response: %w", err)
	}

	return &stkResp, nil
}

func (a *Adapter) B2C(ctx context.Context, req B2CRequest) (*B2CResponse, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"InitiatorName":   "apiuser",
		"CommandID":       "BusinessPayment",
		"Amount":          req.Amount / 100,
		"PartyA":          a.config.Shortcode,
		"PartyB":          req.PhoneNumber,
		"Remarks":         req.Remarks,
		"Occasion":        req.Occasion,
		"QueueTimeOutURL": a.config.CallbackURL + "/timeout",
		"ResultURL":       a.config.CallbackURL + "/result",
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/mpesa/b2c/v1/paymentrequest", a.config.BaseURL())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create B2C request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("B2C request failed: %w", err)
	}
	defer resp.Body.Close()

	var b2cResp B2CResponse
	if err := json.NewDecoder(resp.Body).Decode(&b2cResp); err != nil {
		return nil, fmt.Errorf("failed to decode B2C response: %w", err)
	}

	return &b2cResp, nil
}

func (a *Adapter) QueryStatus(ctx context.Context, req StatusQueryRequest) (*StatusQueryResponse, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Format("20060102150405")
	password := base64.StdEncoding.EncodeToString([]byte(a.config.Shortcode + a.config.Passkey + timestamp))

	payload := map[string]interface{}{
		"BusinessShortCode": a.config.Shortcode,
		"Password":          password,
		"Timestamp":         timestamp,
		"CheckoutRequestID": req.CheckoutRequestID,
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/mpesa/stkpushquery/v1/query", a.config.BaseURL())

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create status query request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("status query failed: %w", err)
	}
	defer resp.Body.Close()

	var queryResp StatusQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&queryResp); err != nil {
		return nil, fmt.Errorf("failed to decode status query response: %w", err)
	}

	return &queryResp, nil
}

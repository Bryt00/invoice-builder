package paystack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

func NewClient(secretKey string) *Client {
	return &Client{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type InitializeRequest struct {
	Email       string         `json:"email"`
	Amount      int            `json:"amount"` // in sub-units (e.g. cents/kobo)
	Reference   string         `json:"reference"`
	CallbackURL string         `json:"callback_url"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type InitializeResponseData struct {
	AuthorizationURL string `json:"authorization_url"`
	AccessCode       string `json:"access_code"`
	Reference        string `json:"reference"`
}

type InitializeResponse struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    InitializeResponseData `json:"data"`
}

type VerifyResponseData struct {
	ID              uint64         `json:"id"`
	Domain          string         `json:"domain"`
	Status          string         `json:"status"`
	Reference       string         `json:"reference"`
	Amount          float64        `json:"amount"`
	Message         string         `json:"message"`
	GatewayResponse string         `json:"gateway_response"`
	PaidAt          string         `json:"paid_at"`
	Channel         string         `json:"channel"`
	Currency        string         `json:"currency"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

type VerifyResponse struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    VerifyResponseData `json:"data"`
}

func (c *Client) InitializeTransaction(req *InitializeRequest) (*InitializeResponse, error) {
	url := fmt.Sprintf("%s/transaction/initialize", c.baseURL)
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.secretKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var initResp InitializeResponse
	if err := json.Unmarshal(body, &initResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !initResp.Status {
		return nil, fmt.Errorf("paystack initialization error: %s", initResp.Message)
	}

	return &initResp, nil
}

func (c *Client) VerifyTransaction(reference string) (*VerifyResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", c.baseURL, reference)

	httpReq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.secretKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var verifyResp VerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !verifyResp.Status {
		return nil, fmt.Errorf("paystack verification error: %s", verifyResp.Message)
	}

	return &verifyResp, nil
}

// VerifyWebhookSignature validates the HMAC SHA-512 signature sent by Paystack.
func VerifyWebhookSignature(body []byte, signatureHeader, secretKey string) bool {
	h := hmac.New(sha512.New, []byte(secretKey))
	h.Write(body)
	expectedSignature := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSignature), []byte(signatureHeader))
}

// Package patentn provides a Go SDK for the Patent N Licensing API
//
// Patent Application #19/429,654
//
// Example:
//
//	client := patentn.NewClient("patent_n_prod_sk_cashapp_...")
//
//	// Detect error
//	detection, err := client.Detect(context.Background(), &patentn.DetectRequest{
//	    ErrorCode:   "OR_CCR_61",
//	    MerchantMCC: patentn.String("5411"),
//	    CardType:    patentn.String("prepaid_debit"),
//	    Amount:      patentn.Float64(50.00),
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Execute bypass
//	if detection.Data.BypassRecommended {
//	    result, err := client.Bypass(context.Background(), &patentn.BypassRequest{
//	        TransactionID: "tx_1234567890",
//	        ErrorCode:     "OR_CCR_61",
//	        Amount:        50.00,
//	        MerchantMCC:   patentn.String("5411"),
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    fmt.Println("Bypass successful:", result.Data.BypassSuccessful)
//	}
package patentn

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultBaseURL = "https://api.patent-n.example.com"
	defaultTimeout = 30 * time.Second
	sdkVersion     = "1.0.0"
)

// Client is the Patent N API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	rateLimit  *RateLimitInfo
}

// Config holds client configuration
type Config struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// RateLimitInfo contains rate limit information from API headers
type RateLimitInfo struct {
	LimitMinute int
	LimitHour   int
	LimitDay    int
	Remaining   int
	Reset       int64
}

// Error represents a Patent N API error
type Error struct {
	Message    string
	StatusCode int
	ErrorCode  string
	Details    map[string]interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("patent-n: %s (status: %d, code: %s)", e.Message, e.StatusCode, e.ErrorCode)
}

// NewClient creates a new Patent N API client
func NewClient(apiKey string, opts ...func(*Config)) *Client {
	config := &Config{
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
		Timeout: defaultTimeout,
	}

	for _, opt := range opts {
		opt(config)
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return &Client{
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		httpClient: config.HTTPClient,
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(url string) func(*Config) {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithTimeout sets a custom timeout
func WithTimeout(timeout time.Duration) func(*Config) {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// Detect detects errors and returns bypass recommendation
func (c *Client) Detect(ctx context.Context, req *DetectRequest) (*DetectResponse, error) {
	var resp DetectResponse
	err := c.doRequest(ctx, "POST", "/api/patent-n/detect", req, &resp)
	return &resp, err
}

// Bypass executes OR_CCR_61 bypass
func (c *Client) Bypass(ctx context.Context, req *BypassRequest) (*BypassResponse, error) {
	var resp BypassResponse
	err := c.doRequest(ctx, "POST", "/api/patent-n/bypass", req, &resp)
	return &resp, err
}

// GetMetrics retrieves performance metrics
func (c *Client) GetMetrics(ctx context.Context, req *MetricsRequest) (*MetricsResponse, error) {
	// Build query params
	params := url.Values{}
	if req.LicenseeID != nil {
		params.Set("licensee_id", *req.LicenseeID)
	}
	if req.StartDate != nil {
		params.Set("start_date", req.StartDate.Format(time.RFC3339))
	}
	if req.EndDate != nil {
		params.Set("end_date", req.EndDate.Format(time.RFC3339))
	}
	if req.TimePeriod != nil {
		params.Set("time_period", *req.TimePeriod)
	}

	endpoint := "/api/patent-n/metrics"
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	var resp MetricsResponse
	err := c.doRequest(ctx, "GET", endpoint, nil, &resp)
	return &resp, err
}

// IngestErrors ingests error logs in batch
func (c *Client) IngestErrors(ctx context.Context, req *ErrorLogBatchRequest) (*ErrorLogBatchResponse, error) {
	var resp ErrorLogBatchResponse
	err := c.doRequest(ctx, "POST", "/api/patent-n/licensee/errors", req, &resp)
	return &resp, err
}

// GetRateLimitInfo returns the current rate limit information
func (c *Client) GetRateLimitInfo() *RateLimitInfo {
	return c.rateLimit
}

// HasRateLimitRemaining checks if rate limit has remaining requests
func (c *Client) HasRateLimitRemaining() bool {
	return c.rateLimit == nil || c.rateLimit.Remaining > 0
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	// Build URL
	reqURL := c.baseURL + endpoint

	// Encode body
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "patent-n-sdk-go/"+sdkVersion)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	// Update rate limit info
	c.updateRateLimitInfo(resp.Header)

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	// Handle errors
	if resp.StatusCode >= 400 {
		var errResp struct {
			Error     string                 `json:"error"`
			ErrorCode string                 `json:"errorCode"`
			Details   map[string]interface{} `json:",inline"`
		}
		_ = json.Unmarshal(respBody, &errResp)

		return &Error{
			Message:    errResp.Error,
			StatusCode: resp.StatusCode,
			ErrorCode:  errResp.ErrorCode,
			Details:    errResp.Details,
		}
	}

	// Decode response
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// updateRateLimitInfo extracts rate limit info from headers
func (c *Client) updateRateLimitInfo(headers http.Header) {
	limitMinute, _ := strconv.Atoi(headers.Get("X-RateLimit-Limit-Minute"))
	limitHour, _ := strconv.Atoi(headers.Get("X-RateLimit-Limit-Hour"))
	limitDay, _ := strconv.Atoi(headers.Get("X-RateLimit-Limit-Day"))
	remaining, _ := strconv.Atoi(headers.Get("X-RateLimit-Remaining"))
	reset, _ := strconv.ParseInt(headers.Get("X-RateLimit-Reset"), 10, 64)

	if limitMinute > 0 {
		c.rateLimit = &RateLimitInfo{
			LimitMinute: limitMinute,
			LimitHour:   limitHour,
			LimitDay:    limitDay,
			Remaining:   remaining,
			Reset:       reset,
		}
	}
}

// Helper functions for optional fields

// String returns a pointer to a string
func String(s string) *string {
	return &s
}

// Float64 returns a pointer to a float64
func Float64(f float64) *float64 {
	return &f
}

// Int returns a pointer to an int
func Int(i int) *int {
	return &i
}

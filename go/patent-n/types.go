package patentn

import "time"

// DetectRequest is the request for error detection
type DetectRequest struct {
	ErrorCode     string   `json:"error_code"`
	MerchantMCC   *string  `json:"merchant_mcc,omitempty"`
	CardType      *string  `json:"card_type,omitempty"`
	Amount        *float64 `json:"amount,omitempty"`
	TransactionID *string  `json:"transaction_id,omitempty"`
}

// DetectResponse is the response from error detection
type DetectResponse struct {
	Success      bool                   `json:"success"`
	Data         DetectData             `json:"data"`
	Licensee     *LicenseeInfo          `json:"licensee,omitempty"`
	Timestamp    string                 `json:"timestamp"`
	ResponseTime string                 `json:"responseTime"`
}

// DetectData contains error detection results
type DetectData struct {
	ErrorFamily          string  `json:"error_family"`
	BypassRecommended    bool    `json:"bypass_recommended"`
	Confidence           float64 `json:"confidence"`
	RecommendedStrategy  string  `json:"recommended_strategy"`
}

// BypassRequest is the request for bypass execution
type BypassRequest struct {
	TransactionID string   `json:"transaction_id"`
	ErrorCode     string   `json:"error_code"`
	Amount        float64  `json:"amount"`
	Currency      *string  `json:"currency,omitempty"`
	MerchantName  *string  `json:"merchant_name,omitempty"`
	MerchantMCC   *string  `json:"merchant_mcc,omitempty"`
	MerchantID    *string  `json:"merchant_id,omitempty"`
	CardBIN       *string  `json:"card_bin,omitempty"`
	CardType      *string  `json:"card_type,omitempty"`
	CardIssuer    *string  `json:"card_issuer,omitempty"`
	UserIDHash    *string  `json:"user_id_hash,omitempty"`
	UserBalance   *float64 `json:"user_balance,omitempty"`
}

// BypassResponse is the response from bypass execution
type BypassResponse struct {
	Success      bool                   `json:"success"`
	Data         BypassData             `json:"data"`
	Timestamp    string                 `json:"timestamp"`
	ResponseTime string                 `json:"responseTime"`
}

// BypassData contains bypass execution results
type BypassData struct {
	BypassSuccessful     bool                   `json:"bypass_successful"`
	BypassTimeMs         int                    `json:"bypass_time_ms"`
	RetryCount           int                    `json:"retry_count"`
	TransformedMetadata  TransformedMetadata    `json:"transformed_metadata"`
	ErrorLogID           string                 `json:"error_log_id"`
	Message              string                 `json:"message"`
}

// TransformedMetadata contains ISO 8583 transformation details
type TransformedMetadata struct {
	Field60Modified  bool `json:"field_60_modified"`
	BINSwapped       bool `json:"bin_swapped"`
	CardTypeChanged  bool `json:"card_type_changed"`
}

// MetricsRequest is the request for performance metrics
type MetricsRequest struct {
	LicenseeID *string     `json:"licensee_id,omitempty"`
	StartDate  *time.Time  `json:"start_date,omitempty"`
	EndDate    *time.Time  `json:"end_date,omitempty"`
	TimePeriod *string     `json:"time_period,omitempty"`
}

// MetricsResponse is the response containing performance metrics
type MetricsResponse struct {
	Success bool         `json:"success"`
	Data    []MetricData `json:"data"`
}

// MetricData contains aggregated metrics
type MetricData struct {
	TimePeriod          string  `json:"time_period"`
	TotalRequests       int     `json:"total_requests"`
	SuccessfulBypasses  int     `json:"successful_bypasses"`
	FailedRequests      int     `json:"failed_requests"`
	AvgRetryTimeMs      float64 `json:"avg_retry_time_ms"`
	RevenueGenerated    float64 `json:"revenue_generated"`
}

// ErrorLogBatchRequest is the request for batch error ingestion
type ErrorLogBatchRequest struct {
	Errors []map[string]interface{} `json:"errors"`
}

// ErrorLogBatchResponse is the response from batch error ingestion
type ErrorLogBatchResponse struct {
	Success bool                `json:"success"`
	Data    ErrorLogBatchData   `json:"data"`
}

// ErrorLogBatchData contains ingestion results
type ErrorLogBatchData struct {
	TotalReceived         int `json:"total_received"`
	SuccessfullyIngested  int `json:"successfully_ingested"`
	Failed                int `json:"failed"`
}

// LicenseeInfo contains licensee information
type LicenseeInfo struct {
	ID   string `json:"id"`
	Tier string `json:"tier"`
}

/**
 * Patent N TypeScript SDK
 * 
 * Official TypeScript client for Patent N Licensing API
 * Patent Application #19/429,654
 * 
 * @example
 * ```typescript
 * import { PatentNClient } from '@patent-n/sdk-typescript';
 * 
 * const client = new PatentNClient({
 *   apiKey: 'patent_n_prod_sk_cashapp_...',
 *   baseURL: 'https://api.patent-n.example.com',
 * });
 * 
 * // Detect error
 * const detection = await client.detect({
 *   error_code: 'OR_CCR_61',
 *   merchant_mcc: '5411',
 *   card_type: 'prepaid_debit',
 *   amount: 50.00,
 * });
 * 
 * // Execute bypass
 * if (detection.bypass_recommended) {
 *   const result = await client.bypass({
 *     transaction_id: 'tx_1234567890',
 *     error_code: 'OR_CCR_61',
 *     amount: 50.00,
 *     merchant_mcc: '5411',
 *   });
 *   console.log(result.bypass_successful);
 * }
 * ```
 */

import axios, { AxiosInstance, AxiosError, AxiosRequestConfig } from 'axios';

// ============================================================================
// Types
// ============================================================================

export interface PatentNConfig {
  apiKey: string;
  baseURL?: string;
  timeout?: number;
  retryAttempts?: number;
  retryDelay?: number;
}

export interface DetectRequest {
  error_code: string;
  merchant_mcc?: string;
  card_type?: 'prepaid_debit' | 'prepaid_credit' | 'debit' | 'credit';
  amount?: number;
  transaction_id?: string;
}

export interface DetectResponse {
  success: boolean;
  data: {
    error_family: string;
    bypass_recommended: boolean;
    confidence: number;
    recommended_strategy: string;
  };
  licensee?: {
    id: string;
    tier: 'standard' | 'premium' | 'enterprise';
  };
  timestamp: string;
  responseTime: string;
}

export interface BypassRequest {
  transaction_id: string;
  error_code: string;
  amount: number;
  currency?: string;
  merchant_name?: string;
  merchant_mcc?: string;
  merchant_id?: string;
  card_bin?: string;
  card_type?: string;
  card_issuer?: string;
  user_id_hash?: string;
  user_balance?: number;
}

export interface BypassResponse {
  success: boolean;
  data: {
    bypass_successful: boolean;
    bypass_time_ms: number;
    retry_count: number;
    transformed_metadata: {
      field_60_modified: boolean;
      bin_swapped: boolean;
      card_type_changed: boolean;
    };
    error_log_id: string;
    message: string;
  };
  timestamp: string;
  responseTime: string;
}

export interface MetricsRequest {
  licensee_id?: string;
  start_date?: Date;
  end_date?: Date;
  time_period?: 'hourly' | 'daily' | 'monthly';
}

export interface MetricsResponse {
  success: boolean;
  data: Array<{
    time_period: string;
    total_requests: number;
    successful_bypasses: number;
    failed_requests: number;
    avg_retry_time_ms: number;
    revenue_generated: number;
  }>;
}

export interface ErrorLogBatchRequest {
  errors: Array<{
    error_code: string;
    amount: number;
    timestamp: Date;
    [key: string]: any;
  }>;
}

export interface ErrorLogBatchResponse {
  success: boolean;
  data: {
    total_received: number;
    successfully_ingested: number;
    failed: number;
  };
}

export interface RateLimitInfo {
  limitMinute: number;
  limitHour: number;
  limitDay: number;
  remaining: number;
  reset: number;
}

export class PatentNError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public errorCode?: string,
    public details?: any
  ) {
    super(message);
    this.name = 'PatentNError';
  }
}

// ============================================================================
// Client
// ============================================================================

export class PatentNClient {
  private client: AxiosInstance;
  private config: Required<PatentNConfig>;
  private rateLimitInfo?: RateLimitInfo;

  constructor(config: PatentNConfig) {
    this.config = {
      apiKey: config.apiKey,
      baseURL: config.baseURL || 'https://api.patent-n.example.com',
      timeout: config.timeout || 30000,
      retryAttempts: config.retryAttempts || 3,
      retryDelay: config.retryDelay || 1000,
    };

    this.client = axios.create({
      baseURL: this.config.baseURL,
      timeout: this.config.timeout,
      headers: {
        'X-API-Key': this.config.apiKey,
        'Content-Type': 'application/json',
        'User-Agent': '@patent-n/sdk-typescript/1.0.0',
      },
    });

    // Response interceptor to capture rate limit info
    this.client.interceptors.response.use(
      (response) => {
        this.updateRateLimitInfo(response.headers);
        return response;
      },
      (error) => {
        if (error.response) {
          this.updateRateLimitInfo(error.response.headers);
        }
        return Promise.reject(this.handleError(error));
      }
    );
  }

  /**
   * Detect errors and get bypass recommendation
   */
  async detect(request: DetectRequest): Promise<DetectResponse> {
    const response = await this.request<DetectResponse>('POST', '/api/patent-n/detect', request);
    return response.data;
  }

  /**
   * Execute bypass for OR_CCR_61 error
   */
  async bypass(request: BypassRequest): Promise<BypassResponse> {
    const response = await this.request<BypassResponse>('POST', '/api/patent-n/bypass', request);
    return response.data;
  }

  /**
   * Get performance metrics
   */
  async getMetrics(request?: MetricsRequest): Promise<MetricsResponse> {
    const params: any = {};
    
    if (request?.licensee_id) params.licensee_id = request.licensee_id;
    if (request?.start_date) params.start_date = request.start_date.toISOString();
    if (request?.end_date) params.end_date = request.end_date.toISOString();
    if (request?.time_period) params.time_period = request.time_period;

    const response = await this.request<MetricsResponse>('GET', '/api/patent-n/metrics', undefined, { params });
    return response.data;
  }

  /**
   * Ingest error logs in batch
   */
  async ingestErrors(request: ErrorLogBatchRequest): Promise<ErrorLogBatchResponse> {
    const response = await this.request<ErrorLogBatchResponse>('POST', '/api/patent-n/licensee/errors', request);
    return response.data;
  }

  /**
   * Get current rate limit information
   */
  getRateLimitInfo(): RateLimitInfo | undefined {
    return this.rateLimitInfo;
  }

  /**
   * Check if rate limit is available
   */
  hasRateLimitRemaining(): boolean {
    return !this.rateLimitInfo || this.rateLimitInfo.remaining > 0;
  }

  /**
   * Wait until rate limit resets
   */
  async waitForRateLimitReset(): Promise<void> {
    if (!this.rateLimitInfo) return;

    const now = Date.now();
    const waitTime = this.rateLimitInfo.reset - now;

    if (waitTime > 0) {
      await new Promise(resolve => setTimeout(resolve, waitTime));
    }
  }

  // ============================================================================
  // Private Methods
  // ============================================================================

  private async request<T>(
    method: string,
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<{ data: T; headers: any }> {
    let lastError: any;

    for (let attempt = 0; attempt < this.config.retryAttempts; attempt++) {
      try {
        const response = await this.client.request<T>({
          method,
          url,
          data,
          ...config,
        });

        return {
          data: response.data,
          headers: response.headers,
        };
      } catch (error: any) {
        lastError = error;

        // Don't retry on client errors (4xx)
        if (error.response && error.response.status < 500) {
          throw error;
        }

        // Wait before retry
        if (attempt < this.config.retryAttempts - 1) {
          await this.sleep(this.config.retryDelay * (attempt + 1));
        }
      }
    }

    throw lastError;
  }

  private updateRateLimitInfo(headers: any): void {
    const limitMinute = parseInt(headers['x-ratelimit-limit-minute'] || '0');
    const limitHour = parseInt(headers['x-ratelimit-limit-hour'] || '0');
    const limitDay = parseInt(headers['x-ratelimit-limit-day'] || '0');
    const remaining = parseInt(headers['x-ratelimit-remaining'] || '0');
    const reset = parseInt(headers['x-ratelimit-reset'] || '0');

    if (limitMinute > 0) {
      this.rateLimitInfo = {
        limitMinute,
        limitHour,
        limitDay,
        remaining,
        reset,
      };
    }
  }

  private handleError(error: AxiosError): PatentNError {
    if (error.response) {
      const data: any = error.response.data;
      return new PatentNError(
        data.error || 'API request failed',
        error.response.status,
        data.errorCode,
        data
      );
    } else if (error.request) {
      return new PatentNError(
        'No response from API',
        0,
        'NETWORK_ERROR'
      );
    } else {
      return new PatentNError(
        error.message,
        0,
        'REQUEST_ERROR'
      );
    }
  }

  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

// ============================================================================
// Exports
// ============================================================================

export default PatentNClient;

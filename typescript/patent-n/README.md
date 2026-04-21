# Patent N TypeScript SDK

Official TypeScript/JavaScript client for the Patent N Licensing API.

**Patent Application:** #19/429,654  
**Status:** Patent Pending

## Installation

```bash
npm install @patent-n/sdk-typescript
# or
yarn add @patent-n/sdk-typescript
```

## Quick Start

```typescript
import { PatentNClient } from '@patent-n/sdk-typescript';

const client = new PatentNClient({
  apiKey: 'patent_n_prod_sk_cashapp_...',
  baseURL: 'https://api.patent-n.example.com',
});

// Detect error
const detection = await client.detect({
  error_code: 'OR_CCR_61',
  merchant_mcc: '5411',
  card_type: 'prepaid_debit',
  amount: 50.00,
});

console.log('Bypass recommended:', detection.data.bypass_recommended);
console.log('Confidence:', detection.data.confidence);

// Execute bypass
if (detection.data.bypass_recommended) {
  const result = await client.bypass({
    transaction_id: 'tx_1234567890',
    error_code: 'OR_CCR_61',
    amount: 50.00,
    merchant_mcc: '5411',
  });
  
  console.log('Bypass successful:', result.data.bypass_successful);
  console.log('Bypass time:', result.data.bypass_time_ms, 'ms');
}
```

## Features

- ✅ Full TypeScript support with type definitions
- ✅ Automatic retries with exponential backoff
- ✅ Rate limiting support with automatic tracking
- ✅ Promise-based API
- ✅ Built-in error handling
- ✅ Request/response logging

## API Reference

### `PatentNClient`

#### Constructor

```typescript
new PatentNClient(config: PatentNConfig)
```

**Config Options:**
- `apiKey` (required): Your Patent N API key
- `baseURL` (optional): API base URL (default: `https://api.patent-n.example.com`)
- `timeout` (optional): Request timeout in ms (default: `30000`)
- `retryAttempts` (optional): Max retry attempts (default: `3`)
- `retryDelay` (optional): Base delay between retries in ms (default: `1000`)

#### Methods

##### `detect(request: DetectRequest): Promise<DetectResponse>`

Detect errors and get bypass recommendation.

##### `bypass(request: BypassRequest): Promise<BypassResponse>`

Execute OR_CCR_61 bypass.

##### `getMetrics(request?: MetricsRequest): Promise<MetricsResponse>`

Get performance metrics.

##### `ingestErrors(request: ErrorLogBatchRequest): Promise<ErrorLogBatchResponse>`

Ingest error logs in batch (up to 1000 errors).

##### `getRateLimitInfo(): RateLimitInfo | undefined`

Get current rate limit information.

##### `hasRateLimitRemaining(): boolean`

Check if rate limit has remaining requests.

##### `waitForRateLimitReset(): Promise<void>`

Wait until rate limit resets.

## Rate Limiting

The SDK automatically tracks rate limits from API response headers:

```typescript
const result = await client.detect({ error_code: 'OR_CCR_61' });

const rateLimit = client.getRateLimitInfo();
console.log('Remaining requests:', rateLimit.remaining);
console.log('Resets at:', new Date(rateLimit.reset));

// Wait for rate limit reset
if (!client.hasRateLimitRemaining()) {
  await client.waitForRateLimitReset();
}
```

## Error Handling

```typescript
try {
  const result = await client.bypass({ ... });
} catch (error) {
  if (error instanceof PatentNError) {
    console.error('API Error:', error.message);
    console.error('Status Code:', error.statusCode);
    console.error('Error Code:', error.errorCode);
  }
}
```

## License

MIT

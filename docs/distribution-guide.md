# Patent N SDK Distribution Guide

> Complete guide for distributing and onboarding licensees with Patent N SDKs

## Overview

This guide covers:
1. Publishing SDKs to internal registry
2. Cash App integration testing
3. Licensee onboarding process
4. API key generation and management
5. Rate limiting and tier management

---

## 1. SDK Publishing

### NPM Registry (TypeScript SDK)

#### Internal Registry Setup

```bash
# Configure internal npm registry
npm config set registry https://npm.internal.patent-n.example.com
npm config set //npm.internal.patent-n.example.com/:_authToken $NPM_TOKEN

# Publish TypeScript SDK
cd /home/user/webapp/sdks/typescript/patent-n
npm publish --access restricted
```

#### Package.json Configuration

```json
{
  "name": "@patent-n/sdk-typescript",
  "version": "1.0.0",
  "publishConfig": {
    "registry": "https://npm.internal.patent-n.example.com",
    "access": "restricted"
  }
}
```

### PyPI (Python SDK)

#### Internal PyPI Setup

```bash
# Configure internal PyPI
cat > ~/.pypirc << EOF
[distutils]
index-servers =
    patent-n

[patent-n]
repository: https://pypi.internal.patent-n.example.com
username: __token__
password: $PYPI_TOKEN
EOF

# Build and publish
cd /home/user/webapp/sdks/python/patent-n
python setup.py sdist bdist_wheel
twine upload --repository patent-n dist/*
```

### Go Module (Go SDK)

#### Private Module Configuration

```bash
# Configure Go private module
export GOPRIVATE=github.com/yourusername/patent-n-go-sdk

# Tag and push
cd /home/user/webapp/sdks/go/patent-n
git tag v1.0.0
git push origin v1.0.0
```

---

## 2. Cash App Integration Testing

### Phase 1: Sandbox Testing

#### Generate Sandbox API Key

```bash
# Use admin dashboard or API
curl -X POST https://admin.patent-n.example.com/api/admin/api-keys/generate \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "licensee_id": "cashapp-001",
    "name": "Cash App Sandbox Key",
    "tier": "enterprise",
    "environment": "sandbox"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "id": "key_abc123...",
    "key": "patent_licensing_sandbox_sk_cashapp_xyz789...",
    "key_prefix": "patent_licensing_sandbox_sk_cashapp_",
    "tier": "enterprise",
    "rate_limits": {
      "per_minute": 500,
      "per_hour": 25000,
      "per_day": 500000
    },
    "expires_at": "2026-01-23T00:00:00Z"
  }
}
```

#### TypeScript Integration Test

```typescript
// test-cash-app-integration.ts
import { PatentNClient } from '@patent-n/sdk-typescript';

const client = new PatentNClient({
  apiKey: 'patent_licensing_sandbox_sk_cashapp_xyz789...',
  baseURL: 'https://sandbox.patent-n.example.com',
});

async function testIntegration() {
  try {
    // 1. Detect OR_CCR_61 error
    console.log('Testing error detection...');
    const detection = await client.detect({
      error_code: 'OR_CCR_61',
      merchant_mcc: '5411', // Grocery store
      card_type: 'prepaid_debit',
      amount: 50.00,
      transaction_id: 'test_tx_001',
    });
    
    console.log('Detection result:', detection);
    console.log('Bypass recommended:', detection.data.bypass_recommended);
    console.log('Confidence:', detection.data.confidence);

    // 2. Execute bypass if recommended
    if (detection.data.bypass_recommended) {
      console.log('Executing bypass...');
      const bypass = await client.bypass({
        transaction_id: 'test_tx_001',
        error_code: 'OR_CCR_61',
        amount: 50.00,
        currency: 'USD',
        merchant_name: 'Safeway',
        merchant_mcc: '5411',
        card_bin: '123456',
        card_type: 'prepaid_debit',
        user_id_hash: 'hash_abc123',
      });
      
      console.log('Bypass result:', bypass);
      console.log('Success:', bypass.data.bypass_successful);
      console.log('Retry time:', bypass.data.bypass_time_ms, 'ms');
    }

    // 3. Get metrics
    console.log('Fetching metrics...');
    const metrics = await client.getMetrics();
    console.log('Metrics:', metrics);

    console.log('✅ All tests passed!');
  } catch (error) {
    console.error('❌ Test failed:', error);
    process.exit(1);
  }
}

testIntegration();
```

Run test:
```bash
npx ts-node test-cash-app-integration.ts
```

### Phase 2: Production Rollout

#### Generate Production API Key

```bash
curl -X POST https://admin.patent-n.example.com/api/admin/api-keys/generate \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "licensee_id": "cashapp-001",
    "name": "Cash App Production Key",
    "tier": "enterprise",
    "environment": "production"
  }'
```

#### Production Configuration

```typescript
// Cash App production config
const client = new PatentNClient({
  apiKey: process.env.PATENT_N_API_KEY, // From secure vault
  baseURL: 'https://api.patent-n.example.com',
  timeout: 5000, // 5 second timeout
  retryAttempts: 3,
});
```

---

## 3. Licensee Onboarding Process

### Step 1: License Agreement

1. **Negotiate Terms**
   - Tier selection (Standard/Premium/Enterprise)
   - Upfront payment amount
   - Revenue share percentage
   - SLA commitments

2. **Sign Agreement**
   - Execute license agreement
   - Upload signed contract to system
   - Activate license in admin dashboard

### Step 2: Technical Onboarding

#### 2.1 Create Licensee Account

```bash
curl -X POST https://admin.patent-n.example.com/api/admin/licensees \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chime",
    "company_type": "Digital Wallet",
    "tier": "premium",
    "email": "licensing@chime.com",
    "phone": "+1-844-244-6363",
    "website": "https://chime.com",
    "webhook_url": "https://api.chime.com/webhooks/patent-n",
    "address_line1": "101 California Street",
    "city": "San Francisco",
    "state": "CA",
    "country": "US",
    "postal_code": "94111"
  }'
```

#### 2.2 Generate API Keys

Generate 2-3 API keys per environment:
- **Production**: Primary + Backup
- **Sandbox**: Testing key

```bash
# Production Primary
curl -X POST https://admin.patent-n.example.com/api/admin/api-keys/generate \
  -d '{"licensee_id": "chime-001", "name": "Chime Production Primary", "tier": "premium"}'

# Production Backup
curl -X POST https://admin.patent-n.example.com/api/admin/api-keys/generate \
  -d '{"licensee_id": "chime-001", "name": "Chime Production Backup", "tier": "premium"}'
```

#### 2.3 Configure Webhooks

Test webhook endpoint:
```bash
curl -X POST https://admin.patent-n.example.com/api/admin/webhooks/test \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "licensee_id": "chime-001",
    "webhook_url": "https://api.chime.com/webhooks/patent-n",
    "event_type": "sla.violation"
  }'
```

### Step 3: SDK Installation

#### TypeScript/Node.js

```bash
# Install from internal registry
npm install @patent-n/sdk-typescript

# Or via package.json
{
  "dependencies": {
    "@patent-n/sdk-typescript": "^1.0.0"
  }
}
```

#### Python

```bash
# Install from internal PyPI
pip install patent-n-sdk --index-url https://pypi.internal.patent-n.example.com

# Or via requirements.txt
patent-n-sdk>=1.0.0
```

#### Go

```bash
# Install via go get
go get github.com/yourusername/patent-n-go-sdk

# Or in go.mod
require github.com/yourusername/patent-n-go-sdk v1.0.0
```

### Step 4: Integration Testing

Provide licensee with integration test suite:

```bash
# Clone test suite
git clone https://github.com/yourusername/patent-n-integration-tests

# Run tests
npm install
npm test -- --licensee=chime --environment=sandbox
```

### Step 5: Production Deployment

**Checklist:**
- [ ] Sandbox testing complete (100% pass rate)
- [ ] Load testing complete (1M+ requests/day)
- [ ] Webhook notifications tested
- [ ] Error handling verified
- [ ] Rate limiting tested
- [ ] Monitoring dashboards configured
- [ ] SLA metrics baseline established
- [ ] Documentation reviewed
- [ ] Production API keys secured in vault
- [ ] Deployment approved by both parties

---

## 4. API Key Management

### Generate New Key

```bash
POST /api/admin/api-keys/generate
{
  "licensee_id": "uuid",
  "name": "Key Name",
  "tier": "standard|premium|enterprise",
  "expires_at": "2026-01-01T00:00:00Z"
}
```

### Revoke Key

```bash
POST /api/admin/api-keys/revoke
{
  "key_id": "uuid",
  "reason": "Security incident"
}
```

### Rotate Key

```bash
POST /api/admin/api-keys/rotate
{
  "key_id": "uuid",
  "grace_period_hours": 24
}
```

**Key Rotation Process:**
1. Generate new key
2. Share new key with licensee
3. Licensee updates their systems
4. Old key remains valid for grace period (24 hours)
5. Old key automatically revoked after grace period

---

## 5. Rate Limiting

### Tier-Based Limits

| Tier | Per Minute | Per Hour | Per Day |
|------|-----------|----------|---------|
| Standard | 100 | 5,000 | 100,000 |
| Premium | 200 | 10,000 | 250,000 |
| Enterprise | 500 | 25,000 | 500,000 |

### Rate Limit Headers

Every API response includes rate limit info:

```
X-RateLimit-Limit-Minute: 500
X-RateLimit-Limit-Hour: 25000
X-RateLimit-Limit-Day: 500000
X-RateLimit-Remaining: 487
X-RateLimit-Reset: 1706055600000
```

### Handling Rate Limits

```typescript
const client = new PatentNClient({ apiKey: '...' });

// Check rate limit before request
if (client.hasRateLimitRemaining()) {
  const result = await client.detect({ error_code: 'OR_CCR_61' });
} else {
  // Wait until rate limit resets
  await client.waitForRateLimitReset();
  const result = await client.detect({ error_code: 'OR_CCR_61' });
}

// Get current rate limit info
const rateLimit = client.getRateLimitInfo();
console.log('Remaining:', rateLimit.remaining);
console.log('Resets at:', new Date(rateLimit.reset));
```

---

## 6. Monitoring & SLA

### Licensee Dashboard

Each licensee has access to:
- Real-time metrics (bypass rate, retry time)
- SLA compliance tracking
- API usage statistics
- Error logs (anonymized)
- Revenue/billing history

Access: `https://licensee.patent-n.example.com/dashboard`

### SLA Commitments

**Patent N Standard SLA:**
- Bypass Rate: ≥94%
- Retry Time: ≤3000ms
- Uptime: ≥99.9%
- Error Rate: ≤5%

**Notifications:**
- SLA violation alerts (webhooks + email)
- API key expiration (30 days, 7 days, 1 day)
- License renewal reminders (60 days, 30 days)

---

## 7. Support & Troubleshooting

### Contact

- **Technical Support**: support@patent-n.example.com
- **Licensing**: licensing@patent-n.example.com
- **Emergency (Critical SLA violations)**: +1-800-PATENT-N

### Common Issues

**Issue: 401 Unauthorized**
- Check API key is valid and not expired
- Verify API key format: `patent_licensing_[env]_sk_[licensee]_[random]`
- Ensure key is sent in `X-API-Key` header

**Issue: 429 Rate Limit Exceeded**
- Check rate limit headers
- Implement exponential backoff
- Consider upgrading tier

**Issue: Bypass not recommended**
- Verify error code is OR_CCR_61 family
- Check merchant MCC is supported
- Review card type (prepaid cards only)

---

## 8. Best Practices

### Security

- **Store API keys securely** (vault, secrets manager)
- **Never commit keys** to version control
- **Rotate keys regularly** (every 90 days)
- **Use separate keys** for production/sandbox
- **Monitor for unauthorized usage**

### Performance

- **Cache bypass recommendations** (30 seconds)
- **Implement circuit breakers** for API failures
- **Use batch error ingestion** for high volume
- **Monitor rate limit usage** proactively

### Error Handling

- **Implement retry logic** with exponential backoff
- **Log all API errors** for debugging
- **Monitor SLA compliance** continuously
- **Set up alerts** for SLA violations

---

## Appendix

### SDK Versions

- **TypeScript**: 1.0.0
- **Python**: 1.0.0
- **Go**: 1.0.0

### API Version

- **Current**: v1
- **Base URL**: https://api.patent-n.example.com

### Changelog

- **2025-01-23**: Initial release (Week 3 implementation)

---

**Need help?** Contact technical support at support@patent-n.example.com

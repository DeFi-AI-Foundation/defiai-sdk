"""
Patent N Python SDK

Official Python client for Patent N Licensing API
Patent Application #19/429,654

Example:
    >>> from patent_n_sdk import PatentNClient
    >>> 
    >>> client = PatentNClient(
    ...     api_key='patent_n_prod_sk_cashapp_...',
    ...     base_url='https://api.patent-n.example.com'
    ... )
    >>> 
    >>> # Detect error
    >>> detection = client.detect(
    ...     error_code='OR_CCR_61',
    ...     merchant_mcc='5411',
    ...     card_type='prepaid_debit',
    ...     amount=50.00
    ... )
    >>> 
    >>> # Execute bypass
    >>> if detection['data']['bypass_recommended']:
    ...     result = client.bypass(
    ...         transaction_id='tx_1234567890',
    ...         error_code='OR_CCR_61',
    ...         amount=50.00,
    ...         merchant_mcc='5411'
    ...     )
    ...     print(result['data']['bypass_successful'])
"""

from .client import PatentNClient, PatentNError, RateLimitInfo
from .types import (
    DetectRequest,
    DetectResponse,
    BypassRequest,
    BypassResponse,
    MetricsRequest,
    MetricsResponse,
    ErrorLogBatchRequest,
    ErrorLogBatchResponse,
)

__version__ = '1.0.0'
__all__ = [
    'PatentNClient',
    'PatentNError',
    'RateLimitInfo',
    'DetectRequest',
    'DetectResponse',
    'BypassRequest',
    'BypassResponse',
    'MetricsRequest',
    'MetricsResponse',
    'ErrorLogBatchRequest',
    'ErrorLogBatchResponse',
]

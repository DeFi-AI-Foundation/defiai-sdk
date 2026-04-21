"""Patent N API Client"""

import time
from typing import Optional, Dict, Any, List
from datetime import datetime
import requests
from requests.adapters import HTTPAdapter
from urllib3.util.retry import Retry

from .types import (
    DetectRequest, DetectResponse,
    BypassRequest, BypassResponse,
    MetricsRequest, MetricsResponse,
    ErrorLogBatchRequest, ErrorLogBatchResponse
)


class RateLimitInfo:
    """Rate limit information from API response headers"""
    
    def __init__(self, headers: Dict[str, str]):
        self.limit_minute = int(headers.get('X-RateLimit-Limit-Minute', 0))
        self.limit_hour = int(headers.get('X-RateLimit-Limit-Hour', 0))
        self.limit_day = int(headers.get('X-RateLimit-Limit-Day', 0))
        self.remaining = int(headers.get('X-RateLimit-Remaining', 0))
        self.reset = int(headers.get('X-RateLimit-Reset', 0))
    
    def has_remaining(self) -> bool:
        """Check if rate limit has remaining requests"""
        return self.remaining > 0
    
    def wait_time(self) -> float:
        """Calculate wait time in seconds until rate limit resets"""
        now = time.time() * 1000  # Convert to milliseconds
        wait = (self.reset - now) / 1000  # Convert back to seconds
        return max(0, wait)


class PatentNError(Exception):
    """Patent N API Error"""
    
    def __init__(self, message: str, status_code: int = 0, 
                 error_code: Optional[str] = None, details: Optional[Dict] = None):
        super().__init__(message)
        self.message = message
        self.status_code = status_code
        self.error_code = error_code
        self.details = details or {}


class PatentNClient:
    """Patent N API Client
    
    Args:
        api_key: Your Patent N API key
        base_url: API base URL (default: https://api.patent-n.example.com)
        timeout: Request timeout in seconds (default: 30)
        max_retries: Maximum retry attempts (default: 3)
        retry_delay: Base delay between retries in seconds (default: 1)
    
    Example:
        >>> client = PatentNClient(api_key='patent_n_prod_sk_...')
        >>> result = client.detect(error_code='OR_CCR_61')
    """
    
    def __init__(
        self,
        api_key: str,
        base_url: str = 'https://api.patent-n.example.com',
        timeout: int = 30,
        max_retries: int = 3,
        retry_delay: float = 1.0
    ):
        self.api_key = api_key
        self.base_url = base_url.rstrip('/')
        self.timeout = timeout
        self.max_retries = max_retries
        self.retry_delay = retry_delay
        self._rate_limit_info: Optional[RateLimitInfo] = None
        
        # Configure session with retries
        self.session = requests.Session()
        retry_strategy = Retry(
            total=max_retries,
            backoff_factor=retry_delay,
            status_forcelist=[500, 502, 503, 504],
            allowed_methods=["GET", "POST"]
        )
        adapter = HTTPAdapter(max_retries=retry_strategy)
        self.session.mount("http://", adapter)
        self.session.mount("https://", adapter)
        
        # Set default headers
        self.session.headers.update({
            'X-API-Key': api_key,
            'Content-Type': 'application/json',
            'User-Agent': 'patent-n-sdk-python/1.0.0'
        })
    
    def detect(self, error_code: str, merchant_mcc: Optional[str] = None,
               card_type: Optional[str] = None, amount: Optional[float] = None,
               transaction_id: Optional[str] = None) -> DetectResponse:
        """Detect errors and get bypass recommendation
        
        Args:
            error_code: Transaction error code (e.g., 'OR_CCR_61')
            merchant_mcc: Merchant Category Code (e.g., '5411')
            card_type: Card type ('prepaid_debit', 'prepaid_credit', etc.)
            amount: Transaction amount
            transaction_id: Unique transaction identifier
        
        Returns:
            DetectResponse with bypass recommendation
        
        Raises:
            PatentNError: If API request fails
        """
        request_data: DetectRequest = {'error_code': error_code}
        
        if merchant_mcc:
            request_data['merchant_mcc'] = merchant_mcc
        if card_type:
            request_data['card_type'] = card_type
        if amount:
            request_data['amount'] = amount
        if transaction_id:
            request_data['transaction_id'] = transaction_id
        
        return self._request('POST', '/api/patent-n/detect', json=request_data)
    
    def bypass(self, transaction_id: str, error_code: str, amount: float,
               currency: str = 'USD', **kwargs) -> BypassResponse:
        """Execute bypass for OR_CCR_61 error
        
        Args:
            transaction_id: Unique transaction identifier
            error_code: Transaction error code
            amount: Transaction amount
            currency: Currency code (default: 'USD')
            **kwargs: Additional bypass parameters (merchant_name, merchant_mcc, etc.)
        
        Returns:
            BypassResponse with execution result
        
        Raises:
            PatentNError: If API request fails
        """
        request_data: BypassRequest = {
            'transaction_id': transaction_id,
            'error_code': error_code,
            'amount': amount,
            'currency': currency,
            **kwargs
        }
        
        return self._request('POST', '/api/patent-n/bypass', json=request_data)
    
    def get_metrics(self, licensee_id: Optional[str] = None,
                    start_date: Optional[datetime] = None,
                    end_date: Optional[datetime] = None,
                    time_period: Optional[str] = None) -> MetricsResponse:
        """Get performance metrics
        
        Args:
            licensee_id: Filter by licensee ID
            start_date: Start date for metrics
            end_date: End date for metrics
            time_period: Time period ('hourly', 'daily', 'monthly')
        
        Returns:
            MetricsResponse with performance data
        """
        params: Dict[str, Any] = {}
        
        if licensee_id:
            params['licensee_id'] = licensee_id
        if start_date:
            params['start_date'] = start_date.isoformat()
        if end_date:
            params['end_date'] = end_date.isoformat()
        if time_period:
            params['time_period'] = time_period
        
        return self._request('GET', '/api/patent-n/metrics', params=params)
    
    def ingest_errors(self, errors: List[Dict[str, Any]]) -> ErrorLogBatchResponse:
        """Ingest error logs in batch
        
        Args:
            errors: List of error log dictionaries
        
        Returns:
            ErrorLogBatchResponse with ingestion results
        """
        request_data: ErrorLogBatchRequest = {'errors': errors}
        return self._request('POST', '/api/patent-n/licensee/errors', json=request_data)
    
    @property
    def rate_limit_info(self) -> Optional[RateLimitInfo]:
        """Get current rate limit information"""
        return self._rate_limit_info
    
    def has_rate_limit_remaining(self) -> bool:
        """Check if rate limit has remaining requests"""
        return not self._rate_limit_info or self._rate_limit_info.has_remaining()
    
    def wait_for_rate_limit_reset(self):
        """Wait until rate limit resets"""
        if self._rate_limit_info:
            wait_time = self._rate_limit_info.wait_time()
            if wait_time > 0:
                time.sleep(wait_time)
    
    def _request(self, method: str, endpoint: str, **kwargs) -> Dict[str, Any]:
        """Make API request with error handling"""
        url = f"{self.base_url}{endpoint}"
        
        try:
            response = self.session.request(
                method=method,
                url=url,
                timeout=self.timeout,
                **kwargs
            )
            
            # Update rate limit info
            self._rate_limit_info = RateLimitInfo(response.headers)
            
            # Handle errors
            if not response.ok:
                error_data = response.json() if response.content else {}
                raise PatentNError(
                    message=error_data.get('error', 'API request failed'),
                    status_code=response.status_code,
                    error_code=error_data.get('errorCode'),
                    details=error_data
                )
            
            return response.json()
            
        except requests.exceptions.Timeout:
            raise PatentNError(
                message='Request timeout',
                status_code=0,
                error_code='TIMEOUT'
            )
        except requests.exceptions.ConnectionError:
            raise PatentNError(
                message='Connection error',
                status_code=0,
                error_code='CONNECTION_ERROR'
            )
        except requests.exceptions.RequestException as e:
            raise PatentNError(
                message=str(e),
                status_code=0,
                error_code='REQUEST_ERROR'
            )
    
    def __enter__(self):
        """Context manager entry"""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        """Context manager exit"""
        self.session.close()

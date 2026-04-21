"""Type definitions for Patent N SDK"""

from typing import TypedDict, Optional, List, Literal


class DetectRequest(TypedDict, total=False):
    error_code: str
    merchant_mcc: Optional[str]
    card_type: Optional[Literal['prepaid_debit', 'prepaid_credit', 'debit', 'credit']]
    amount: Optional[float]
    transaction_id: Optional[str]


class DetectResponse(TypedDict):
    success: bool
    data: dict
    licensee: Optional[dict]
    timestamp: str
    responseTime: str


class BypassRequest(TypedDict, total=False):
    transaction_id: str
    error_code: str
    amount: float
    currency: str
    merchant_name: Optional[str]
    merchant_mcc: Optional[str]
    merchant_id: Optional[str]
    card_bin: Optional[str]
    card_type: Optional[str]
    card_issuer: Optional[str]
    user_id_hash: Optional[str]
    user_balance: Optional[float]


class BypassResponse(TypedDict):
    success: bool
    data: dict
    timestamp: str
    responseTime: str


class MetricsRequest(TypedDict, total=False):
    licensee_id: Optional[str]
    start_date: Optional[str]
    end_date: Optional[str]
    time_period: Optional[Literal['hourly', 'daily', 'monthly']]


class MetricsResponse(TypedDict):
    success: bool
    data: List[dict]


class ErrorLogBatchRequest(TypedDict):
    errors: List[dict]


class ErrorLogBatchResponse(TypedDict):
    success: bool
    data: dict

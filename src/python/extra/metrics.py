# Metrics Module
# Contains metrics collection and analytics functionality

import time
from datetime import datetime, timezone
from typing import Dict, Optional
from dataclasses import dataclass, asdict
import aiohttp

# LogEntry definition (local definition for metrics module)
@dataclass
class LogEntry:
    timestamp: str
    method: str
    path: str
    user_agent: str
    is_ai_crawler: bool
    status_code: int
    remote_addr: str
    referer: str

# AI Crawlers list (local definition for metrics module)
AI_CRAWLERS = [
    "GPTBot",
    "OAI-SearchBot",
    "ChatGPT-User",
    "anthropic-ai",
    "ClaudeBot",
    "PerplexityBot",
    "Google-Extended",
    "Bytespider",
    "Amazonbot",
    "Applebot-Extended",
    "CCBot"
]

@dataclass
class Metrics:
    total_requests: int = 0
    ai_crawler_requests: int = 0
    proxy_requests: int = 0
    error_requests: int = 0
    start_time: str = ""
    last_request_time: str = ""
    status_code_counts: Dict[int, int] = None
    user_agent_counts: Dict[str, int] = None
    path_counts: Dict[str, int] = None
    ai_crawler_types: Dict[str, int] = None
    average_response_time: float = 0.0
    requests_by_hour: Dict[str, int] = None
    
    def __post_init__(self):
        if self.status_code_counts is None:
            self.status_code_counts = {}
        if self.user_agent_counts is None:
            self.user_agent_counts = {}
        if self.path_counts is None:
            self.path_counts = {}
        if self.ai_crawler_types is None:
            self.ai_crawler_types = {}
        if self.requests_by_hour is None:
            self.requests_by_hour = {}

@dataclass
class AnalyticsEntry:
    timestamp: str
    type: str  # "metrics" or "request"
    data: Optional[Metrics] = None
    log_entry: Optional[LogEntry] = None

class MetricsCollector:
    def __init__(self):
        self.metrics = Metrics(
            start_time=datetime.now(timezone.utc).isoformat(),
            status_code_counts={},
            user_agent_counts={},
            path_counts={},
            ai_crawler_types={},
            requests_by_hour={}
        )
        self.response_times = []
    
    def record_request(self, log_entry: LogEntry, response_time: float):
        self.metrics.total_requests += 1
        self.metrics.last_request_time = log_entry.timestamp
        
        # Update counts
        self.metrics.status_code_counts[log_entry.status_code] = \
            self.metrics.status_code_counts.get(log_entry.status_code, 0) + 1
        self.metrics.path_counts[log_entry.path] = \
            self.metrics.path_counts.get(log_entry.path, 0) + 1
        
        # Track requests by hour
        hour_key = datetime.fromisoformat(log_entry.timestamp.replace('Z', '+00:00')).strftime("%Y-%m-%d-%H")
        self.metrics.requests_by_hour[hour_key] = \
            self.metrics.requests_by_hour.get(hour_key, 0) + 1
        
        # Track user agents (truncate to avoid memory issues)
        user_agent = log_entry.user_agent[:100] if len(log_entry.user_agent) > 100 else log_entry.user_agent
        self.metrics.user_agent_counts[user_agent] = \
            self.metrics.user_agent_counts.get(user_agent, 0) + 1
        
        if log_entry.is_ai_crawler:
            self.metrics.ai_crawler_requests += 1
            # Track specific AI crawler types
            for crawler in AI_CRAWLERS:
                if crawler in log_entry.user_agent:
                    self.metrics.ai_crawler_types[crawler] = \
                        self.metrics.ai_crawler_types.get(crawler, 0) + 1
                    break
        else:
            self.metrics.proxy_requests += 1
        
        if log_entry.status_code >= 400:
            self.metrics.error_requests += 1
        
        # Track response times (keep only last 1000)
        self.response_times.append(response_time)
        if len(self.response_times) > 1000:
            self.response_times = self.response_times[1:]
        
        # Calculate average response time
        if self.response_times:
            self.metrics.average_response_time = sum(self.response_times) / len(self.response_times)
    
    def get_metrics(self) -> Metrics:
        return Metrics(
            total_requests=self.metrics.total_requests,
            ai_crawler_requests=self.metrics.ai_crawler_requests,
            proxy_requests=self.metrics.proxy_requests,
            error_requests=self.metrics.error_requests,
            start_time=self.metrics.start_time,
            last_request_time=self.metrics.last_request_time,
            average_response_time=self.metrics.average_response_time,
            status_code_counts=self.metrics.status_code_counts.copy(),
            user_agent_counts=self.metrics.user_agent_counts.copy(),
            path_counts=self.metrics.path_counts.copy(),
            ai_crawler_types=self.metrics.ai_crawler_types.copy(),
            requests_by_hour=self.metrics.requests_by_hour.copy()
        )

# Analytics service configuration
ANALYTICS_SERVICE_URL = 'https://httpbin.org/post'

async def forward_analytics(analytics_entry: AnalyticsEntry):
    """Send analytics data to the configured analytics service"""
    try:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=15)) as session:
            headers = {
                'Content-Type': 'application/json',
                'User-Agent': 'Cloudflare-Worker-Analytics/1.0'
            }
            
            analytics_data = {
                'timestamp': analytics_entry.timestamp,
                'type': analytics_entry.type
            }
            
            if analytics_entry.data:
                analytics_data['data'] = asdict(analytics_entry.data)
            if analytics_entry.log_entry:
                analytics_data['log_entry'] = asdict(analytics_entry.log_entry)
            
            async with session.post(ANALYTICS_SERVICE_URL, 
                                  json=analytics_data, 
                                  headers=headers) as response:
                if not (200 <= response.status < 300):
                    print(f"Analytics service returned non-success status: {response.status}")
    except Exception as e:
        print(f"Failed to forward analytics: {e}")

# Global metrics instance
global_metrics = MetricsCollector()

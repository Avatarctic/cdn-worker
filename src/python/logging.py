# Logging Module
# Contains log entry structure and log forwarding functionality

from dataclasses import dataclass, asdict
import aiohttp

# Log service configuration
LOG_SERVICE_URL = 'https://httpbin.org/post'

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

async def forward_log(log_entry: LogEntry):
    """Send log data to the configured log service"""
    try:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=15)) as session:
            headers = {
                'Content-Type': 'application/json',
                'User-Agent': 'Cloudflare-Worker-Logger/1.0'
            }
            
            log_data = asdict(log_entry)
            
            async with session.post(LOG_SERVICE_URL, 
                                  json=log_data, 
                                  headers=headers) as response:
                if not (200 <= response.status < 300):
                    print(f"Log service returned non-success status: {response.status}")
                    error_text = await response.text()
                    print(f"Log service error response: {error_text}")
    except Exception as e:
        print(f"Failed to forward log: {e}")

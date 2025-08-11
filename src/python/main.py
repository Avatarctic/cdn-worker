import asyncio
import json
import time
from datetime import datetime, timezone
from typing import Dict, List, Optional
from dataclasses import dataclass, asdict
from collections import defaultdict
import aiohttp

# Import functionality from separate modules
from logging import LogEntry, forward_log
from ai_detection import is_ai_crawler, AI_CRAWLERS

# Origin URL
ORIGIN_URL = "http://www.mysite.com/"

async def handle_request(request):
    """Main request handler function"""
    start_time = time.time()
    current_time = datetime.now(timezone.utc).isoformat()
    
    # Extract request information
    method = request.method
    url = str(request.url)
    path = request.url.path
    user_agent = request.headers.get('User-Agent', '')
    remote_addr = request.client.host if hasattr(request, 'client') else 'unknown'
    referer = request.headers.get('Referer', '')
    
    # Create log entry
    log_entry = LogEntry(
        timestamp=current_time,
        method=method,
        path=path,
        user_agent=user_agent,
        is_ai_crawler=is_ai_crawler(user_agent),
        status_code=200,  # Will be updated
        remote_addr=remote_addr,
        referer=referer
    )
    
    try:
        # Check for AI crawlers
        if is_ai_crawler(user_agent):
            # Serve AI-optimized content
            log_entry.status_code = 200
            response_content = """
            <html>
                <head><title>AI-Ready Content</title></head>
                <body>This is content optimized for AI crawlers.</body>
            </html>
            """
            
            # Log the AI crawler request
            print(f"AI Crawler detected: {user_agent} from {remote_addr} accessing {path}")
            
            # Send log
            await forward_log(log_entry)
            
            return {
                'statusCode': 200,
                'headers': {
                    'Content-Type': 'text/html; charset=utf-8'
                },
                'body': response_content
            }
        
        else:
            # Act as a transparent proxy
            try:
                origin_url = f"{ORIGIN_URL.rstrip('/')}{path}"
                
                async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=30)) as session:
                    # Forward the request to origin
                    async with session.request(
                        method=method,
                        url=origin_url,
                        headers=dict(request.headers),
                        data=await request.body() if hasattr(request, 'body') else None
                    ) as origin_response:
                        
                        log_entry.status_code = origin_response.status
                        
                        # Log the proxy request
                        print(f"Proxied request: {method} {path} -> {origin_response.status}")
                        
                        # Send log
                        await forward_log(log_entry)
                        
                        # Return proxied response
                        response_body = await origin_response.text()
                        response_headers = dict(origin_response.headers)
                        
                        return {
                            'statusCode': origin_response.status,
                            'headers': response_headers,
                            'body': response_body
                        }
            
            except Exception as proxy_error:
                print(f"Failed to fetch from origin: {proxy_error}")
                log_entry.status_code = 502
                
                await forward_log(log_entry)
                
                return {
                    'statusCode': 502,
                    'headers': {'Content-Type': 'text/plain'},
                    'body': 'Failed to fetch from origin'
                }
    
    except Exception as e:
        print(f"Error in handle_request: {e}")
        log_entry.status_code = 500
        
        await forward_log(log_entry)
        
        return {
            'statusCode': 500,
            'headers': {'Content-Type': 'text/plain'},
            'body': 'Internal Server Error'
        }

# For Cloudflare Workers with Pyodide
async def fetch_handler(request):
    """Entry point for Cloudflare Workers"""
    return await handle_request(request)

# Initialize and start background tasks
def initialize():
    """Initialize the worker"""
    print("Cloudflare Worker initialized")
    
# Call initialize when module loads
initialize()

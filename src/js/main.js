// Cloudflare Worker that executes Python-like logic in JavaScript
// This approach avoids the complexity of loading Pyodide at runtime

import { isAICrawler } from './ai_detection.js';
import { forwardLog } from './logging.js';

// Origin URL
const ORIGIN_URL = "http://www.mysite.com/";

// Main request handler
async function handleRequest(request) {
  const startTime = Date.now();
  const currentTime = new Date().toISOString();
  
  // Extract request information
  const url = new URL(request.url);
  const userAgent = request.headers.get('User-Agent') || '';
  const remoteAddr = request.headers.get('CF-Connecting-IP') || 'unknown';
  const referer = request.headers.get('Referer') || '';
  
  // Create log entry
  const logEntry = {
    timestamp: currentTime,
    method: request.method,
    path: url.pathname,
    userAgent: userAgent,
    isAICrawler: isAICrawler(userAgent),
    statusCode: 200, // Will be updated
    remoteAddr: remoteAddr,
    referer: referer
  };
  
  try {
    // Check for AI crawlers
    if (isAICrawler(userAgent)) {
      // Serve AI-optimized content
      logEntry.statusCode = 200;
      const responseContent = `
        <html>
          <head><title>AI-Ready Content</title></head>
          <body>This is content optimized for AI crawlers.</body>
        </html>
      `;
      
      // Log the AI crawler request
      console.log(`AI Crawler detected: ${userAgent} from ${remoteAddr} accessing ${url.pathname}`);
      
      // Send log
      forwardLog(logEntry);
      
      return new Response(responseContent, {
        status: 200,
        headers: { 'Content-Type': 'text/html; charset=utf-8' }
      });
    } else {
      // Act as a transparent proxy
      try {
        const originUrl = `${ORIGIN_URL.replace(/\/$/, '')}${url.pathname}${url.search}`;
        
        // Forward the request to origin
        const originResponse = await fetch(originUrl, {
          method: request.method,
          headers: request.headers,
          body: request.body
        });
        
        logEntry.statusCode = originResponse.status;
        
        // Log the proxy request
        console.log(`Proxied request: ${request.method} ${url.pathname} -> ${originResponse.status}`);
        
        // Send log
        forwardLog(logEntry);
        
        // Return proxied response
        return new Response(originResponse.body, {
          status: originResponse.status,
          headers: originResponse.headers
        });
        
      } catch (proxyError) {
        console.log(`Failed to fetch from origin: ${proxyError}`);
        logEntry.statusCode = 502;
        
        forwardLog(logEntry);
        
        return new Response('Failed to fetch from origin', {
          status: 502,
          headers: { 'Content-Type': 'text/plain' }
        });
      }
    }
  } catch (error) {
    console.log(`Error in handleRequest: ${error}`);
    logEntry.statusCode = 500;
    
    forwardLog(logEntry);
    
    return new Response('Internal Server Error', {
      status: 500,
      headers: { 'Content-Type': 'text/plain' }
    });
  }
}

// Main export for Cloudflare Workers
export default {
  async fetch(request, env, ctx) {
    return await handleRequest(request);
  }
};


// Metrics Module
// Contains metrics collection and analytics functionality

// AI Crawlers list (local definition for metrics module)
const aiCrawlers = [
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
];

// MetricsCollector class for tracking analytics data
class MetricsCollector {
  constructor() {
    this.metrics = {
      totalRequests: 0,
      aiCrawlerRequests: 0,
      proxyRequests: 0,
      errorRequests: 0,
      startTime: new Date(),
      lastRequestTime: new Date(),
      statusCodeCounts: {},
      userAgentCounts: {},
      pathCounts: {},
      aiCrawlerTypes: {},
      averageResponseTime: 0,
      requestsByHour: {},
    };
    this.responseTimes = []; // Keep last 1000
  }

  recordRequest(logEntry, responseTime) {
    this.metrics.totalRequests++;
    this.metrics.lastRequestTime = logEntry.timestamp;
    
    // Track status codes
    this.metrics.statusCodeCounts[logEntry.statusCode] = (this.metrics.statusCodeCounts[logEntry.statusCode] || 0) + 1;
    
    // Track paths
    this.metrics.pathCounts[logEntry.path] = (this.metrics.pathCounts[logEntry.path] || 0) + 1;
    
    // Track requests by hour
    const hourKey = logEntry.timestamp.toISOString().slice(0, 13); // YYYY-MM-DDTHH
    this.metrics.requestsByHour[hourKey] = (this.metrics.requestsByHour[hourKey] || 0) + 1;
    
    // Track user agents (truncate to avoid memory issues)
    let userAgent = logEntry.userAgent;
    if (userAgent && userAgent.length > 100) {
      userAgent = userAgent.substring(0, 100);
    }
    this.metrics.userAgentCounts[userAgent] = (this.metrics.userAgentCounts[userAgent] || 0) + 1;
    
    if (logEntry.isAICrawler) {
      this.metrics.aiCrawlerRequests++;
      // Track specific AI crawler types
      for (const crawler of aiCrawlers) {
        if (logEntry.userAgent && logEntry.userAgent.includes(crawler)) {
          this.metrics.aiCrawlerTypes[crawler] = (this.metrics.aiCrawlerTypes[crawler] || 0) + 1;
          break;
        }
      }
    } else {
      this.metrics.proxyRequests++;
    }
    
    if (logEntry.statusCode >= 400) {
      this.metrics.errorRequests++;
    }
    
    // Track response times (keep only last 1000)
    this.responseTimes.push(responseTime);
    if (this.responseTimes.length > 1000) {
      this.responseTimes.shift();
    }
    
    // Calculate average response time
    if (this.responseTimes.length > 0) {
      const total = this.responseTimes.reduce((sum, time) => sum + time, 0);
      this.metrics.averageResponseTime = total / this.responseTimes.length;
    }
  }

  getMetrics() {
    // Return a deep copy of metrics
    return JSON.parse(JSON.stringify(this.metrics));
  }
}

// Analytics service configuration
const ANALYTICS_SERVICE_URL = 'https://httpbin.org/post';

// Forward analytics data to the configured analytics service
async function forwardAnalytics(analyticsEntry) {
  try {
    const response = await fetch(ANALYTICS_SERVICE_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Cloudflare-Worker-Analytics/1.0',
      },
      body: JSON.stringify(analyticsEntry),
    });

    if (!response.ok) {
      console.error(`Analytics service returned non-success status: ${response.status}`);
    }
  } catch (error) {
    console.error('Failed to forward analytics:', error);
  }
}

// PERIODIC ANALYTICS FUNCTIONALITY (PRESERVED FOR FUTURE USE)
// Uncomment the function below to enable periodic analytics

// function startPeriodicAnalytics(metricsCollector) {
//   setInterval(async () => {
//     const metrics = metricsCollector.getMetrics();
//     const analytics = {
//       timestamp: new Date(),
//       type: "metrics",
//       data: metrics,
//     };
//     await forwardAnalytics(analytics);
//     console.log(`Sent periodic analytics: ${metrics.totalRequests} total requests, ${metrics.aiCrawlerRequests} AI crawler requests`);
//   }, 5 * 60 * 1000); // 5 minutes
// }

// Export functions and classes
export { MetricsCollector, forwardAnalytics };

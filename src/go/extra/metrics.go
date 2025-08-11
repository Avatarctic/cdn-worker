package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LogEntry represents a log entry (local definition for metrics module)
type LogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	UserAgent   string    `json:"user_agent"`
	IsAICrawler bool      `json:"is_ai_crawler"`
	StatusCode  int       `json:"status_code"`
	RemoteAddr  string    `json:"remote_addr"`
	Referer     string    `json:"referer"`
}

// AI Crawlers list (local definition for metrics module)
var aiCrawlers = []string{
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
	"CCBot",
}

// Metrics holds various analytics data
type Metrics struct {
	TotalRequests       int64            `json:"total_requests"`
	AICrawlerRequests   int64            `json:"ai_crawler_requests"`
	ProxyRequests       int64            `json:"proxy_requests"`
	ErrorRequests       int64            `json:"error_requests"`
	StartTime           time.Time        `json:"start_time"`
	LastRequestTime     time.Time        `json:"last_request_time"`
	StatusCodeCounts    map[int]int64    `json:"status_code_counts"`
	UserAgentCounts     map[string]int64 `json:"user_agent_counts"`
	PathCounts          map[string]int64 `json:"path_counts"`
	AICrawlerTypes      map[string]int64 `json:"ai_crawler_types"`
	AverageResponseTime time.Duration    `json:"average_response_time"`
	RequestsByHour      map[string]int64 `json:"requests_by_hour"`
}

// MetricsCollector holds the metrics with thread-safe access
type MetricsCollector struct {
	mutex         sync.RWMutex
	metrics       Metrics
	responseTimes []time.Duration // Keep last 1000, not serialized
}

// Analytics service configuration
const analyticsServiceURL = "https://httpbin.org/post"

// AnalyticsEntry represents analytics data to be sent to the analytics service
type AnalyticsEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "metrics" or "request"
	Data      Metrics   `json:"data,omitempty"`
	LogEntry  LogEntry  `json:"log_entry,omitempty"`
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: Metrics{
			StatusCodeCounts: make(map[int]int64),
			UserAgentCounts:  make(map[string]int64),
			PathCounts:       make(map[string]int64),
			AICrawlerTypes:   make(map[string]int64),
			RequestsByHour:   make(map[string]int64),
			StartTime:        time.Now(),
		},
		responseTimes: make([]time.Duration, 0, 1000),
	}
}

// RecordRequest updates metrics for a request
func (mc *MetricsCollector) RecordRequest(logEntry LogEntry, responseTime time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.TotalRequests++
	mc.metrics.LastRequestTime = logEntry.Timestamp
	mc.metrics.StatusCodeCounts[logEntry.StatusCode]++
	mc.metrics.PathCounts[logEntry.Path]++

	// Track requests by hour
	hourKey := logEntry.Timestamp.Format("2006-01-02-15")
	mc.metrics.RequestsByHour[hourKey]++

	// Track user agents (truncate to avoid memory issues)
	userAgent := logEntry.UserAgent
	if len(userAgent) > 100 {
		userAgent = userAgent[:100]
	}
	mc.metrics.UserAgentCounts[userAgent]++

	if logEntry.IsAICrawler {
		mc.metrics.AICrawlerRequests++
		// Track specific AI crawler types
		for _, crawler := range aiCrawlers {
			if strings.Contains(logEntry.UserAgent, crawler) {
				mc.metrics.AICrawlerTypes[crawler]++
				break
			}
		}
	} else {
		mc.metrics.ProxyRequests++
	}

	if logEntry.StatusCode >= 400 {
		mc.metrics.ErrorRequests++
	}

	// Track response times (keep only last 1000)
	mc.responseTimes = append(mc.responseTimes, responseTime)
	if len(mc.responseTimes) > 1000 {
		mc.responseTimes = mc.responseTimes[1:]
	}

	// Calculate average response time
	if len(mc.responseTimes) > 0 {
		var total time.Duration
		for _, rt := range mc.responseTimes {
			total += rt
		}
		mc.metrics.AverageResponseTime = total / time.Duration(len(mc.responseTimes))
	}
}

// GetMetrics returns a copy of current metrics
func (mc *MetricsCollector) GetMetrics() Metrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	// Create a copy
	metrics := Metrics{
		TotalRequests:       mc.metrics.TotalRequests,
		AICrawlerRequests:   mc.metrics.AICrawlerRequests,
		ProxyRequests:       mc.metrics.ProxyRequests,
		ErrorRequests:       mc.metrics.ErrorRequests,
		StartTime:           mc.metrics.StartTime,
		LastRequestTime:     mc.metrics.LastRequestTime,
		AverageResponseTime: mc.metrics.AverageResponseTime,
		StatusCodeCounts:    make(map[int]int64),
		UserAgentCounts:     make(map[string]int64),
		PathCounts:          make(map[string]int64),
		AICrawlerTypes:      make(map[string]int64),
		RequestsByHour:      make(map[string]int64),
	}

	// Copy maps
	for k, v := range mc.metrics.StatusCodeCounts {
		metrics.StatusCodeCounts[k] = v
	}
	for k, v := range mc.metrics.UserAgentCounts {
		metrics.UserAgentCounts[k] = v
	}
	for k, v := range mc.metrics.PathCounts {
		metrics.PathCounts[k] = v
	}
	for k, v := range mc.metrics.AICrawlerTypes {
		metrics.AICrawlerTypes[k] = v
	}
	for k, v := range mc.metrics.RequestsByHour {
		metrics.RequestsByHour[k] = v
	}

	return metrics
}

// forwardAnalytics sends analytics data to the configured analytics service
func forwardAnalytics(analyticsEntry AnalyticsEntry) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		analyticsData, err := json.Marshal(analyticsEntry)
		if err != nil {
			log.Printf("Failed to marshal analytics entry: %v", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", analyticsServiceURL, bytes.NewBuffer(analyticsData))
		if err != nil {
			log.Printf("Failed to create analytics request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Cloudflare-Worker-Analytics/1.0")

		client := &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: false,
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to forward analytics: %v", err)
			return
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Printf("Failed to close analytics response body: %v", closeErr)
			}
		}()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Analytics service returned non-success status: %d", resp.StatusCode)
		}
	}()
}

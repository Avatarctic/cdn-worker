package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// Log service configuration
const logServiceURL = "https://httpbin.org/post"

// LogEntry represents a log entry to be sent to the log service
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

// forwardLog sends log data to the configured log service
func forwardLog(logEntry LogEntry) {
	go func() {
		// Add timeout context
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		logData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", logServiceURL, bytes.NewBuffer(logData))
		if err != nil {
			log.Printf("Failed to create log request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Cloudflare-Worker-Logger/1.0")

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
			log.Printf("Failed to forward log: %v", err)
			return
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Printf("Failed to close log response body: %v", closeErr)
			}
		}()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Log service returned non-success status: %d", resp.StatusCode)
			// Read error response for debugging
			if body, readErr := io.ReadAll(resp.Body); readErr == nil {
				log.Printf("Log service error response: %s", string(body))
			}
		}
	}()
}

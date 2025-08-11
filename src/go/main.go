package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

// recoveryMiddleware wraps handlers with panic recovery
func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)

				// Log the panic
				logEntry := LogEntry{
					Timestamp:   time.Now(),
					Method:      r.Method,
					Path:        r.URL.Path,
					UserAgent:   r.Header.Get("User-Agent"),
					IsAICrawler: IsAICrawler(r.Header.Get("User-Agent")),
					StatusCode:  http.StatusInternalServerError,
					RemoteAddr:  r.RemoteAddr,
					Referer:     r.Header.Get("Referer"),
				}
				forwardLog(logEntry)

				// Return error response
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

// The URL of the actual website's origin server.
const originURL = "http://www.mysite.com/"

// This is the main Worker handler function, invoked for every request.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	userAgent := r.Header.Get("User-Agent")
	remoteAddr := r.RemoteAddr
	referer := r.Header.Get("Referer")

	// Prepare log entry
	logEntry := LogEntry{
		Timestamp:   startTime,
		Method:      r.Method,
		Path:        r.URL.Path,
		UserAgent:   userAgent,
		IsAICrawler: IsAICrawler(userAgent),
		RemoteAddr:  remoteAddr,
		Referer:     referer,
	}

	// 1. Check for AI crawlers.
	if IsAICrawler(userAgent) {
		// 2. If it's an AI crawler, serve the alternate response.
		logEntry.StatusCode = http.StatusOK
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		response := "<html><head><title>AI-Ready Content</title></head><body>This is content optimized for AI crawlers.</body></html>"
		if _, err := w.Write([]byte(response)); err != nil {
			log.Printf("Failed to write AI crawler response: %v", err)
			logEntry.StatusCode = http.StatusInternalServerError
		}

		// Log the AI crawler request
		log.Printf("AI Crawler detected: %s from %s accessing %s", userAgent, remoteAddr, r.URL.Path)
		forwardLog(logEntry)

		return
	}

	// 3. If it's NOT an AI crawler, act as a transparent proxy.

	// Create a new request to the origin, preserving the original method and URL path.
	originReq, err := http.NewRequest(r.Method, originURL+r.URL.Path, r.Body)
	if err != nil {
		logEntry.StatusCode = http.StatusInternalServerError
		log.Printf("Failed to create origin request: %v", err)
		forwardLog(logEntry)
		http.Error(w, "Failed to create origin request", http.StatusInternalServerError)
		return
	}

	// Copy all headers from the original request to the new origin request.
	originReq.Header = r.Header

	// Make the request to the origin server with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: false,
		},
	}
	originResp, err := client.Do(originReq)
	if err != nil {
		logEntry.StatusCode = http.StatusBadGateway
		log.Printf("Failed to fetch from origin: %v", err)
		forwardLog(logEntry)
		http.Error(w, "Failed to fetch from origin", http.StatusBadGateway)
		return
	}
	defer func() {
		if closeErr := originResp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close origin response body: %v", closeErr)
		}
	}()

	// Update log entry with origin response status
	logEntry.StatusCode = originResp.StatusCode

	// Log the proxy request
	log.Printf("Proxied request: %s %s -> %d", r.Method, r.URL.Path, originResp.StatusCode)
	forwardLog(logEntry)

	// 4. Copy the response from the origin back to the client.

	// Copy all headers from the origin response to the Worker response.
	for key, values := range originResp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set the status code.
	w.WriteHeader(originResp.StatusCode)

	// Copy the response body.
	if _, err := io.Copy(w, originResp.Body); err != nil {
		log.Printf("Failed to copy response body: %v (this may be expected if client disconnected)", err)
		// Note: We can't change status code here as headers are already sent
		// This error is often expected when clients disconnect early
	}
}

func main() {
	// For Cloudflare Workers deployment: Keep worker alive
	select {}
}

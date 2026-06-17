package handlers

import (
	"net/http"
)

// PingHandler provides a lightweight endpoint for latency and jitter measurement.
// It returns a 204 No Content response, minimizing payload and processing overhead.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	// CRITICAL: Prevent all layers of caching (Browser, CDN, Reverse Proxy)
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Explicitly signal the client to keep the TCP connection open
	// (Go HTTP/1.1 does this by default, but explicit is better for edge cases)
	w.Header().Set("Connection", "keep-alive")

	// Immediately flush the headers with a 204 status.
	// This tells the Go server to close the response writer without waiting for a body.
	w.WriteHeader(http.StatusNoContent)
}

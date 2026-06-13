package handlers

import (
	"io"
	"net/http"

	"github.com/ali-hasehmi/speedtest/internal/config"
	"github.com/ali-hasehmi/speedtest/logger"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxBytes := config.UploadMaxSize()

	// 1. FAST EXIT: Check the header if it exists.
	// This saves us from reading the body at all if the client is honest about being too large.
	if r.ContentLength > maxBytes {
		http.Error(w, "payload too large", http.StatusRequestEntityTooLarge)
		return
	}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// 2. STRICT ENFORCEMENT: Wrap the body.
	// This protects us if the client lied about Content-Length or used chunked encoding.
	lr := io.LimitReader(r.Body, maxBytes)

	// Consume the stream. If the client tries to send more than maxBytes,
	// lr will return EOF early, effectively truncating the read and saving our resources.
	_, err := io.Copy(io.Discard, lr)
	if err != nil {
		logger.Errorf("upload error: %v\n", err)
		http.Error(w, "upload failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

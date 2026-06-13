package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/ali-hasehmi/speedtest/internal/config"
	"github.com/ali-hasehmi/speedtest/internal/speedtest"
	"github.com/ali-hasehmi/speedtest/logger"
)

func TestMain(m *testing.M) {
	// Suppress any log messages
	logger.SetLevel(logger.NONE)

	// Initialize the buffer using the default from config.go
	speedtest.InitBuffer(config.DownloadBufferSize())

	os.Exit(m.Run())
}

func TestDownloadHandler(t *testing.T) {
	maxSize := config.DownloadMaxSize()

	tests := []struct {
		name           string
		querySize      string
		expectedStatus int
		expectedBytes  int64
	}{
		{
			name:           "Valid Size",
			querySize:      "1024",
			expectedStatus: http.StatusOK,
			expectedBytes:  1024,
		},
		{
			name:           "Invalid Size String",
			querySize:      "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBytes:  0,
		},
		{
			name:           "Zero Size",
			querySize:      "0",
			expectedStatus: http.StatusBadRequest,
			expectedBytes:  0,
		},
		{
			name:           "Negative Size",
			querySize:      "-500",
			expectedStatus: http.StatusBadRequest,
			expectedBytes:  0,
		},
		{
			name:           "Bigger than Maximum Size",
			querySize:      strconv.Itoa(int(maxSize + 1)),
			expectedStatus: http.StatusOK,
			expectedBytes:  maxSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/download?s="+tt.querySize, nil)
			rr := httptest.NewRecorder()

			DownloadHandler(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedStatus == http.StatusOK {
				if int64(rr.Body.Len()) != tt.expectedBytes {
					t.Errorf("handler returned wrong body size: got %v want %v", rr.Body.Len(), tt.expectedBytes)
				}
			}
		})
	}
}

func TestUploadHandler(t *testing.T) {
	maxUploadSize := config.UploadMaxSize()

	tests := []struct {
		name           string
		method         string
		body           []byte
		forceHeader    int64
		expectedStatus int
	}{
		{
			name:           "Valid POST Upload",
			method:         http.MethodPost,
			body:           make([]byte, 1024),
			forceHeader:    1024,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Method GET",
			method:         http.MethodGet,
			body:           nil,
			forceHeader:    0,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Payload Too Large (Header Fast Exit)",
			method:         http.MethodPost,
			body:           nil,
			forceHeader:    maxUploadSize + 1,
			expectedStatus: http.StatusRequestEntityTooLarge,
		},
		{
			name:           "Payload Too Large (Chunked/Hidden)",
			method:         http.MethodPost,
			body:           make([]byte, maxUploadSize+100),
			forceHeader:    -1,            // Simulate chunked encoding
			expectedStatus: http.StatusOK, // LimitReader truncates it safely
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyReader io.Reader
			if tt.body != nil {
				bodyReader = bytes.NewReader(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/api/upload", bodyReader)
			req.ContentLength = tt.forceHeader

			rr := httptest.NewRecorder()
			UploadHandler(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

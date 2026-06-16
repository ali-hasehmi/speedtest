package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ali-hasehmi/speedtest/internal/metadata"
)

func TestIPHandler_ContentNegotiation(t *testing.T) {
	tests := []struct {
		name           string
		remoteAddr     string
		queryString    string
		acceptHeader   string
		wantStatusCode int
		wantType       string
		mustContain    []string
		mustNotContain []string
	}{
		{
			name:           "Simple Mode - Pure Raw IP (Default)",
			remoteAddr:     "1.1.1.1:12345",
			queryString:    "",
			acceptHeader:   "*/*",
			wantStatusCode: http.StatusOK,
			wantType:       "text/plain; charset=utf-8",
			mustContain:    []string{"1.1.1.1"},
			mustNotContain: []string{"IP:", "Country:", "ASN:"},
		},
		{
			name:           "Simple Mode - Query present but not full mode",
			remoteAddr:     "8.8.8.8:54321",
			queryString:    "cache_buster=true",
			acceptHeader:   "*/*",
			wantStatusCode: http.StatusOK,
			wantType:       "text/plain; charset=utf-8",
			mustContain:    []string{"8.8.8.8"},
			mustNotContain: []string{"IP:", "ASN:"},
		},
		{
			name:           "Full Mode - Plaintext Format",
			remoteAddr:     "1.1.1.1:12345",
			queryString:    "mode=full",
			acceptHeader:   "text/plain",
			wantStatusCode: http.StatusOK,
			wantType:       "text/plain; charset=utf-8",
			mustContain:    []string{"IP: 1.1.1.1\n"},
		},
		{
			name:           "Full Mode - JSON via Query Parameter",
			remoteAddr:     "8.8.8.8:443",
			queryString:    "mode=full&format=json",
			acceptHeader:   "*/*",
			wantStatusCode: http.StatusOK,
			wantType:       "application/json",
			mustContain:    []string{"\"ip\":\"8.8.8.8\""},
		},
		{
			name:           "Full Mode - JSON via Accept Header",
			remoteAddr:     "1.1.1.1:80",
			queryString:    "mode=full",
			acceptHeader:   "application/json",
			wantStatusCode: http.StatusOK,
			wantType:       "application/json",
			mustContain:    []string{"\"ip\":\"1.1.1.1\""},
		},
		{
			name:           "Invalid RemoteAddr Error",
			remoteAddr:     "not-an-ip-address",
			queryString:    "",
			acceptHeader:   "*/*",
			wantStatusCode: http.StatusBadRequest,
			wantType:       "text/plain; charset=utf-8",
			mustContain:    []string{"invalid client ip address"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/api/ip"
			if tt.queryString != "" {
				url += "?" + tt.queryString
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.RemoteAddr = tt.remoteAddr
			if tt.acceptHeader != "" {
				req.Header.Set("Accept", tt.acceptHeader)
			}

			rec := httptest.NewRecorder()
			IPHandler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("Expected status code %d, got %d", tt.wantStatusCode, rec.Code)
			}

			gotType := rec.Header().Get("Content-Type")
			if gotType != tt.wantType {
				t.Errorf("Expected Content-Type '%s', got '%s'", tt.wantType, gotType)
			}

			body := rec.Body.String()

			for _, item := range tt.mustContain {
				if !strings.Contains(body, item) {
					t.Errorf("Expected body to contain %q, but it didn't.\nBody: %s", item, body)
				}
			}

			for _, item := range tt.mustNotContain {
				if strings.Contains(body, item) {
					t.Errorf("Expected body NOT to contain %q, but it did.\nBody: %s", item, body)
				}
			}

			if tt.wantType == "application/json" && rec.Code == http.StatusOK {
				var target metadata.Info
				if err := json.Unmarshal(rec.Body.Bytes(), &target); err != nil {
					t.Errorf("Handler returned malformed JSON object: %v", err)
				}
			}
		})
	}
}

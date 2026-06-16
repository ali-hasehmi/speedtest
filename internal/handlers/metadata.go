package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"github.com/ali-hasehmi/speedtest/internal/metadata"
	"github.com/ali-hasehmi/speedtest/logger"
)

// IPHandler manages both simple IP reflection and full metadata enrichment.
// Endpoints:
//   - GET /api/ip (Returns plain IP string, ultra-fast)
//   - GET /api/ip?mode=full (Returns plaintext metadata)
//   - GET /api/ip?mode=full&format=json (Returns JSON metadata)
func IPHandler(w http.ResponseWriter, r *http.Request) {
	ip, err := parseAddrOnly(r.RemoteAddr)
	if err != nil {
		http.Error(w, "invalid client ip address", http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")

	info := metadata.Info{
		IP: ip.String(),
	}

	// Optimization: Fast-path for the standard 'simple' text request.
	// If there are no query parameters at all, immediately stream the plain IP and exit.
	if r.URL.RawQuery == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, info.IP)
		return
	}

	// Content Negotiation
	query := r.URL.Query()
	mode := query.Get("mode")
	format := query.Get("format")

	fullMode := mode == "full"
	serveJSON := format == "json" || (format == "" && strings.Contains(r.Header.Get("Accept"), "application/json"))

	if fullMode {
		info = metadata.Lookup(ip)
	}

	// JSON format handling
	if serveJSON {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(info); err != nil {
			logger.Errorf("failed to encode metadata JSON: %v", err)
		}
		return
	}

	// Default fallback: Plain Text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if !fullMode {
		io.WriteString(w, info.IP)
		return
	}

	var sb strings.Builder
	sb.Grow(256) // Pre-allocate buffer capacity to avoid resize allocations

	sb.WriteString("IP: ")
	sb.WriteString(info.IP)
	sb.WriteByte('\n')

	if info.Country != "" {
		sb.WriteString("Country: ")
		sb.WriteString(info.Country)
		sb.WriteByte('\n')
	}
	if info.ASN != 0 {
		sb.WriteString("ASN: AS")
		sb.WriteString(strconv.FormatUint(uint64(info.ASN), 10))
		sb.WriteByte('\n')
	}
	if info.ISP != "" {
		sb.WriteString("ISP: ")
		sb.WriteString(info.ISP)
		sb.WriteByte('\n')
	}

	io.WriteString(w, sb.String())
}

// parseAddrOnly handles parsing strings like "192.0.2.1:54321" or "[2001:db8::1]:54321" safely
// into a standard allocation-free netip.Addr using standard standard library routines.
func parseAddrOnly(remoteAddr string) (netip.Addr, error) {
	// standard net/http RemoteAddr includes ports (e.g. "127.0.0.1:12345")
	addrPort, err := netip.ParseAddrPort(remoteAddr)
	if err == nil {
		return addrPort.Addr(), nil
	}

	// If chi middleware strips it down or we receive an address without a port:
	addr, err := netip.ParseAddr(remoteAddr)
	if err != nil {
		return netip.Addr{}, err
	}
	return addr, nil
}

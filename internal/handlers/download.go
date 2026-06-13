package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/ali-hasehmi/speedtest/internal/config"
	"github.com/ali-hasehmi/speedtest/internal/speedtest"
	"github.com/ali-hasehmi/speedtest/logger"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	maxBytes := config.DownloadMaxSize()
	sizeStr := r.URL.Query().Get("s")
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if errors.Is(err, strconv.ErrRange) {
		size = maxBytes
		err = nil
	}
	if err != nil || size <= 0 {
		logger.Errorf("invalid size: %s\n", err)
		http.Error(w, "invalid size", http.StatusBadRequest)
		return
	}
	if size > maxBytes {
		size = maxBytes
	}

	// Anti-cache headers for CDNs and reverse-proxies e.g. nginx
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	rr := speedtest.NewRepeatReader(size)
	io.Copy(w, rr)
}

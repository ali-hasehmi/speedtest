package config

import (
	"net/netip"
	"os"
	"strconv"
	"time"

	"github.com/ali-hasehmi/speedtest/logger"
)

func getEnvNetipAddr(key string, defaultVal netip.Addr) netip.Addr {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	val, err := netip.ParseAddr(v)
	if err != nil {
		logger.Fatalf("invalid %s: %v", key, err)
	}
	return val
}

func getEnvUint16(key string, defaultVal uint16) uint16 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	val, err := strconv.ParseUint(v, 10, 16)
	if err != nil {
		logger.Fatalf("invalid %s: %v", key, err)
	}
	return uint16(val)
}

func getEnvTimeDuration(key string, defaultVal time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	val, err := time.ParseDuration(v)
	if err != nil {
		logger.Fatalf("invalid %s: %v", key, err)
	}
	return val
}

func getEnvLogLevel(key string, defaultVal logger.LogLevel) logger.LogLevel {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	val, err := parseLogLevel(v)
	if err != nil {
		logger.Fatalf("invalid %s: %v", key, err)
	}
	return val
}

func getEnvByteSize(key string, defaultVal int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	val, err := parseByteSize(v)
	if err != nil {
		logger.Fatalf("invalid %s: %v", key, err)
	}
	return val
}

func getEnvString(key string, defaultVal string) string {
	v := os.Getenv(key)
	if v == "" {
		v = defaultVal
	}
	return v
}

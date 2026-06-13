package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ali-hasehmi/speedtest/logger"
)

func parseLogLevel(logLevelStr string) (logger.LogLevel, error) {

	// Warning, WARNING, warning are the same
	logLevelStr = strings.ToLower(logLevelStr)
	m := map[string]logger.LogLevel{
		"debug":   logger.DEBUG,
		"info":    logger.INFO,    // 1
		"warning": logger.WARNING, // 2
		"error":   logger.ERROR,   // 3
		"none":    logger.NONE,    // 4
	}
	loglevel, ok := m[logLevelStr]
	if !ok {
		return 0, errors.New("invalid loglevel value")
	}
	return loglevel, nil
}

// Parses values like 100MB to its equivalent byte count which is 100 * 1024 * 1024
func parseByteSize(s string) (int64, error) {

	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return 0, errors.New("empty size string")
	}

	// Find where numbers end and units begin
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if (c < '0' || c > '9') && c != '.' {
			break
		}
	}

	numStr := s[:i]
	unitStr := strings.TrimSpace(s[i:])

	// Handle case where it's just a number (default to bytes)
	// if unitStr == "" || unitStr == "b" {
	// 	val, err := strconv.ParseInt(numStr, 10, 64)
	// 	return val, err
	// }

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	var multiplier float64 = 1
	switch unitStr {
	case "", "b":
		multiplier = 1
	case "k", "kb":
		multiplier = 1000
	case "ki", "kib":
		multiplier = 1024
	case "m", "mb":
		multiplier = 1000 * 1000
	case "mi", "mib":
		multiplier = 1024 * 1024
	case "g", "gb":
		multiplier = 1000 * 1000 * 1000
	case "gi", "gib":
		multiplier = 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown unit: %s", unitStr)
	}

	return int64(val * multiplier), nil
}

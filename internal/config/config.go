package config

import (
	"errors"
	"net/netip"
	"os"
	"time"

	"github.com/ali-hasehmi/speedtest/logger"
	"github.com/joho/godotenv"
)

var (
	listenAddr netip.Addr = netip.IPv4Unspecified()
	listenPort uint16     = 8080

	readTimeout  = time.Duration(5 * time.Second)
	writeTimeout = time.Duration(10 * time.Second)
	idleTimeout  = time.Duration(30 * time.Second)

	downloadBufferSize int64 = 10 * 1024 * 1024  // Size of pre‑generated random buffer
	downloadMaxSize    int64 = 100 * 1024 * 1024 // Max allowed download size per request, 0 for none
	uploadMaxSize      int64 = 50 * 1024 * 1024  // Max allowed download size per request

	cityDBPath = "" // path to the City/Country MMDB file, If empty, city-level lookup is skipped entirely.
	asnDBPath  = "" // path to the ASN/ISP MMDB file. If empty, ASN/ISP lookup is skipped entirely.

	logLevel = logger.INFO
	logFile  = "" // path to write log data
)

// Load reads environment variables and overrides defaults
func Load(filenames ...string) error {

	// load env files
	err := godotenv.Load(filenames...)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	// get environments and override defaults if needed
	listenAddr = getEnvNetipAddr("LISTEN_ADDR", listenAddr)
	listenPort = getEnvUint16("LISTEN_PORT", listenPort)

	readTimeout = getEnvTimeDuration("READ_TIMEOUT", readTimeout)
	writeTimeout = getEnvTimeDuration("WRITE_TIMEOUT", writeTimeout)
	idleTimeout = getEnvTimeDuration("IDLE_TIMEOUT", idleTimeout)

	downloadBufferSize = getEnvByteSize("DOWNLOAD_BUFFER_SIZE", downloadBufferSize)
	downloadMaxSize = getEnvByteSize("DOWNLOAD_MAX_SIZE", downloadMaxSize)
	uploadMaxSize = getEnvByteSize("UPLOAD_MAX_SIZE", uploadMaxSize)

	cityDBPath = getEnvString("CITY_DB_PATH", cityDBPath)
	asnDBPath = getEnvString("ASN_DB_PATH", asnDBPath)

	logLevel = getEnvLogLevel("LOG_LEVEL", logLevel)
	logFile = getEnvString("LOG_FILE", logFile)

	return nil
}

// Getter functions
func ListenAddr() netip.Addr {
	return listenAddr
}
func ListenPort() uint16 {
	return listenPort
}
func ReadTimeout() time.Duration {
	return readTimeout
}
func WriteTimeout() time.Duration {
	return writeTimeout
}
func IdleTimeout() time.Duration {
	return idleTimeout
}
func DownloadBufferSize() int64 {
	return downloadBufferSize
}
func DownloadMaxSize() int64 {
	return downloadMaxSize
}
func UploadMaxSize() int64 {
	return uploadMaxSize
}
func LogLevel() logger.LogLevel {
	return logLevel
}
func LogFile() string {
	return logFile
}
func CityDBPath() string {
	return cityDBPath
}
func AsnDBPath() string {
	return asnDBPath
}

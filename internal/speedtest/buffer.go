package speedtest

import (
	"crypto/rand"

	"github.com/ali-hasehmi/speedtest/logger"
)

var (
	randomBuffer []byte
)

func InitBuffer(sizeBytes int64) {
	randomBuffer = make([]byte, sizeBytes)
	_, err := rand.Read(randomBuffer)
	if err != nil {
		logger.Fatal("failed to fill buffer with random data:", err)
	}
}

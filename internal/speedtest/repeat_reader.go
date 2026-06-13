package speedtest

import (
	"io"

	"github.com/ali-hasehmi/speedtest/logger"
)

type RepeatReader struct {
	pos  int
	left int64
}

func NewRepeatReader(size int64) *RepeatReader {
	return &RepeatReader{left: size}
}

func (r *RepeatReader) Read(p []byte) (n int, err error) {
	logger.Infof("repeatReader:Reader called\n")
	if r.left <= 0 {
		return 0, io.EOF
	}

	if r.pos >= len(randomBuffer) {
		r.pos = 0 // loop
	}

	toCopy := len(randomBuffer) - r.pos
	if int64(toCopy) > r.left {
		toCopy = int(r.left)
	}
	n = copy(p, randomBuffer[r.pos:r.pos+toCopy])
	r.pos += n
	r.left -= int64(n)
	logger.Infof("pos=%v, left=%v\n", r.pos, r.left)
	return n, nil
}

func (r *RepeatReader) WriteTo(w io.Writer) (n int64, err error) {
	for {
		if r.left <= 0 {
			return n, nil
		}
		if r.pos >= len(randomBuffer) {
			r.pos = 0
		}
		toCopy := len(randomBuffer) - r.pos
		if int64(toCopy) > r.left {
			toCopy = int(r.left)
		}
		nn, err := w.Write(randomBuffer[r.pos : r.pos+toCopy])
		if err != nil {
			logger.Error(err)
			return n, err
		}
		r.pos += nn
		r.left -= int64(nn)
		logger.Infof("left=%v, pos=%v\n", r.left, r.pos)
	}
}

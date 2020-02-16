package go_engine_io_parser

import (
	"errors"
	"net/http"
	"time"
)

// Checker is function to check request.
type Checker func(*http.Request) (http.Header, error)

// ErrInvalidFrame is returned when writing invalid frame type.
var ErrInvalidFrame = errors.New("invalid frame type")

// ErrInvalidPacket is returned when writing invalid packet type.
var ErrInvalidPacket = errors.New("invalid packet type")

var chars = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_")

// Timestamp returns a string based on different nano time.
func Timestamp() string {
	now := time.Now().UnixNano()
	ret := make([]byte, 0, 16)
	for now > 0 {
		ret = append(ret, chars[int(now%int64(len(chars)))])
		now = now / int64(len(chars))
	}
	return string(ret)
}

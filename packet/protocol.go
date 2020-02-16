// Package packet is codec of packet for connection which supports framing
package packet

import (
	"github.com/go-engine.io-parser"
	"io"
)

// FrameReader is the reader which supports framing
type FrameReader interface {
	NextReader() (FrameType, io.ReadCloser, error)
}

// FrameWriter is the writer which supports framing
type FrameWriter interface {
	NextWriter(typ FrameType) (io.WriteCloser, error)
}

// NewEncoder creates a packet encoder which writes to w
func NewEncoder(writer FrameWriter) FrameWriter {
	return go_engine_io_parser.newEncoder(writer)
}

// NewDecoder creates a packet decoder which reads from reader
func NewDecoder(reader FrameReader) FrameReader {
	return go_engine_io_parser.newDecoder(reader)
}

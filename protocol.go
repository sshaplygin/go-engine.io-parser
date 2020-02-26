package go_engine_io_parser

import (
	"github.com/go-engine.io-parser/frame"
	"io"
)

// FrameReader is the reader which supports framing
type FrameReader interface {
	NextReader() (frame.FrameType, io.ReadCloser, error)
}

// FrameWriter is the writer which supports framing
type FrameWriter interface {
	NextWriter(typ frame.FrameType) (io.WriteCloser, error)
}

// NewEncoder creates a packet encoder which writes to w
func NewEncoder(writer FrameWriter) FrameWriter {
	return newEncoder(writer)
}

// NewDecoder creates a packet decoder which reads from reader
func NewDecoder(reader FrameReader) FrameReader {
	return newDecoder(reader)
}

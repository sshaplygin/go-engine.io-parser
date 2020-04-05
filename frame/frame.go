package frame

import (
	"github.com/go-engine.io-parser/packet"
	"io"
)

// FrameType is the type of frames.
type FrameType byte

const (
	// FrameString identifies a string frame.
	FrameString FrameType = iota
	// FrameBinary identifies a binary frame.
	FrameBinary
)

// ByteToFrameType converts a byte to FrameType.
func ByteToFrameType(b byte) FrameType {
	return FrameType(b)
}

// Byte returns type in byte.
func (t FrameType) Byte() byte {
	return byte(t)
}

// FrameReader reads a frame. It need be closed before next reading.
type FrameReader interface {
	NextReader() (FrameType, packet.PacketType, io.ReadCloser, error)
}

// FrameWriter writes a frame. It need be closed before next writing.
type FrameWriter interface {
	NextWriter(ft FrameType, pt packet.PacketType) (io.WriteCloser, error)
}

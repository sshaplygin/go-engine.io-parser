package go_engine_io_parser

// FrameType is the type of frames.
type FrameType byte

const (
	// FrameString identifies a string frame.
	FrameString FrameType = iota
	// FrameBinary identifies a binary frame.
	FrameBinary
)

// Byte returns type in byte.
func (t FrameType) Byte() byte {
	return byte(t)
}

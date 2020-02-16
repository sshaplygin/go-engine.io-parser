package go_engine_io_parser

import (
	"github.com/go-engine.io-parser/packet"
	"io"
)

type decoder struct {
	frameReader packet.FrameReader
}

func newDecoder(frameReader packet.FrameReader) *decoder {
	return &decoder{
		frameReader: frameReader,
	}
}

func (dec *decoder) NextReader() (FrameType, PacketType, io.ReadCloser, error) {
	ft, r, err := dec.frameReader.NextReader()
	if err != nil {
		return 0, 0, nil, err
	}
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		r.Close()
		return 0, 0, nil, err
	}
	return ft, ByteToPacketType(b[0], ft), r, nil
}

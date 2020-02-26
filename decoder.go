package go_engine_io_parser

import (
	"github.com/mrfoe7/go-engine.io-parser/frame"
	"github.com/mrfoe7/go-engine.io-parser/packet"
	"io"
)

type decoder struct {
	frameReader parser.FrameReader
}

func newDecoder(frameReader parser.FrameReader) *decoder {
	return &decoder{
		frameReader: frameReader,
	}
}

func (dec *decoder) NextReader() (frame.FrameType, packet.PacketType, io.ReadCloser, error) {
	ft, r, err := dec.frameReader.NextReader()
	if err != nil {
		return 0, 0, nil, err
	}
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		err := r.Close()
		//todo:
		return 0, 0, nil, err
	}
	return ft, packet.ByteToPacketType(b[0], ft), r, nil
}

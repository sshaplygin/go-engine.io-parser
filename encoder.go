package go_engine_io_parser

import (
	"github.com/go-engine.io-parser/packet"
	"io"
)

type encoder struct {
	w packet.FrameWriter
}

func newEncoder(w packet.FrameWriter) *encoder {
	return &encoder{
		w: w,
	}
}

func (e *encoder) NextWriter(ft FrameType, pt PacketType) (io.WriteCloser, error) {
	w, err := e.w.NextWriter(ft)
	if err != nil {
		return nil, err
	}
	var b [1]byte
	if ft == FrameString {
		b[0] = pt.StringByte()
	} else {
		b[0] = pt.BinaryByte()
	}
	if _, err := w.Write(b[:]); err != nil {
		w.Close()
		return nil, err
	}
	return w, nil
}

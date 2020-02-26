package go_engine_io_parser

import (
	"github.com/mrfoe7/go-engine.io-parser/frame"
	"github.com/mrfoe7/go-engine.io-parser/packet"
	"io"
)

type encoder struct {
	w frame.FrameWriter
}

func newEncoder(w frame.FrameWriter) *encoder {
	return &encoder{
		w: w,
	}
}

func (e *encoder) NextWriter(ft frame.FrameType, pt packet.PacketType) (io.WriteCloser, error) {
	w, err := e.w.NextWriter(ft)
	if err != nil {
		return nil, err
	}
	var b [1]byte
	if ft == frame.FrameString {
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

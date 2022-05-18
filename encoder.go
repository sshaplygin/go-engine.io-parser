package go_engine_io_parser

import (
	"bytes"
	"encoding/base64"
	"io"
)

// Writer writes a frame. It need be Closed before next writing.
type Writer interface {
	NextWriter(ft FrameType, pt Type) (io.WriteCloser, error)
}

type writerFeeder interface {
	getWriter() (io.Writer, error)
	putWriter(error) error
}

type encoder struct {
	supportBinary bool
	feeder        writerFeeder

	ft         FrameType
	pt         Type
	header     bytes.Buffer
	frameCache bytes.Buffer
	b64Writer  io.WriteCloser
	rawWriter  io.Writer
}

func (e *encoder) Noop() []byte {
	if e.supportBinary {
		return []byte{0x00, 0x01, 0xff, '6'}
	}
	return []byte("1:6")
}

func (e *encoder) NextWriter(ft FrameType, pt Type) (io.WriteCloser, error) {
	w, err := e.feeder.getWriter()
	if err != nil {
		return nil, err
	}
	e.rawWriter = w

	e.ft = ft
	e.pt = pt
	e.frameCache.Reset()

	if !e.supportBinary && ft == FrameBinary {
		e.b64Writer = base64.NewEncoder(base64.StdEncoding, &e.frameCache)
	} else {
		e.b64Writer = nil
	}
	return e, nil
}

func (e *encoder) Write(p []byte) (int, error) {
	if e.b64Writer != nil {
		return e.b64Writer.Write(p)
	}
	return e.frameCache.Write(p)
}

func (e *encoder) Close() error {
	if e.b64Writer != nil {
		e.b64Writer.Close()
	}

	var writeHeader func() error
	if e.supportBinary {
		writeHeader = e.writeBinaryHeader
	} else {
		if e.ft == FrameBinary {
			writeHeader = e.writeB64Header
		} else {
			writeHeader = e.writeTextHeader
		}
	}

	e.header.Reset()
	err := writeHeader()
	if err == nil {
		_, err = e.header.WriteTo(e.rawWriter)
	}
	if err == nil {
		_, err = e.frameCache.WriteTo(e.rawWriter)
	}
	err = e.feeder.putWriter(err)
	return err
}

func (e *encoder) writeTextHeader() error {
	l := int64(e.frameCache.Len() + 1) // length for packet type
	err := writeTextLen(l, &e.header)
	if err == nil {
		err = e.header.WriteByte(e.pt.StringByte())
	}
	return err
}

func (e *encoder) writeB64Header() error {
	l := int64(e.frameCache.Len() + 2) // length for 'b' and packet type
	err := writeTextLen(l, &e.header)
	if err == nil {
		err = e.header.WriteByte(BinarySymbol)
	}
	if err == nil {
		err = e.header.WriteByte(e.pt.StringByte())
	}
	return err
}

func (e *encoder) writeBinaryHeader() error {
	l := int64(e.frameCache.Len() + 1) // length for packet type
	b := e.pt.StringByte()
	if e.ft == FrameBinary {
		b = e.pt.Byte()
	}
	err := e.header.WriteByte(e.ft.Byte())
	if err == nil {
		err = writeBinaryLen(l, &e.header)
	}
	if err == nil {
		err = e.header.WriteByte(b)
	}
	return err
}

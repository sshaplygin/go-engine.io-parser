package packet

import (
	"bytes"
	"github.com/googollee/go-engine.io/base"
	"io"
	"io/ioutil"
)

type fakeConnReader struct {
	frames []Frame
}

func newFakeConnReader(frames []Frame) *fakeConnReader {
	return &fakeConnReader{
		frames: frames,
	}
}

func (reader *fakeConnReader) NextReader() (FrameType, io.ReadCloser, error) {
	if len(reader.frames) == 0 {
		return FrameString, nil, io.EOF
	}
	frame := reader.frames[0]
	reader.frames = reader.frames[1:]
	return frame.fType, ioutil.NopCloser(bytes.NewReader(f.data)), nil
}

type fakeFrame struct {
	writer *fakeConnWriter
	fType  FrameType
	data   *bytes.Buffer
}

func newFakeFrame(writer *fakeConnWriter, fType FrameType) *fakeFrame {
	return &fakeFrame{
		writer: writer,
		fType:  fType,
		data:   bytes.NewBuffer(nil),
	}
}

func (w *fakeFrame) Write(p []byte) (int, error) {
	return w.data.Write(p)
}

func (w *fakeFrame) Read(p []byte) (int, error) {
	return w.data.Read(p)
}

func (w *fakeFrame) Close() error {
	if w.w == nil {
		return nil
	}
	w.w.frames = append(w.w.frames, Frame{
		typ:  w.typ,
		data: w.data.Bytes(),
	})
	return nil
}

type fakeConnWriter struct {
	frames []Frame
}

func newFakeConnWriter() *fakeConnWriter {
	return &fakeConnWriter{}
}

func (w *fakeConnWriter) NextWriter(typ base.FrameType) (io.WriteCloser, error) {
	return newFakeFrame(w, typ), nil
}

type fakeOneFrameConst struct {
	b byte
}

func (c *fakeOneFrameConst) Read(p []byte) (int, error) {
	p[0] = c.b
	return 1, nil
}

type fakeConstReader struct {
	ft base.FrameType
	r  *fakeOneFrameConst
}

func newFakeConstReader() *fakeConstReader {
	return &fakeConstReader{
		ft: base.FrameString,
		r: &fakeOneFrameConst{
			b: base.MESSAGE.StringByte(),
		},
	}
}

func (r *fakeConstReader) NextReader() (base.FrameType, io.ReadCloser, error) {
	ft := r.ft
	switch ft {
	case base.FrameBinary:
		r.ft = base.FrameString
		r.r.b = base.MESSAGE.StringByte()
	case base.FrameString:
		r.ft = base.FrameBinary
		r.r.b = base.MESSAGE.BinaryByte()
	}
	return ft, ioutil.NopCloser(r.r), nil
}

type fakeOneFrameDiscarder struct{}

func (d fakeOneFrameDiscarder) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d fakeOneFrameDiscarder) Close() error {
	return nil
}

type fakeDiscardWriter struct{}

func (w *fakeDiscardWriter) NextWriter(typ base.FrameType) (io.WriteCloser, error) {
	return fakeOneFrameDiscarder{}, nil
}

package go_engine_io_parser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mrfoe7/go-engine.io-parser/frame"
	"github.com/mrfoe7/go-engine.io-parser/packet"
)

var (
	writeErr = errors.New("write error")
)

type fakeWriterFeeder struct {
	w           io.Writer
	returnError error
	passingErr  error
}

func (f *fakeWriterFeeder) getWriter() (io.Writer, error) {
	return f.w, f.returnError
}

func (f *fakeWriterFeeder) putWriter(err error) error {
	f.passingErr = err
	return f.returnError
}

func TestEncoder(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	f := &fakeWriterFeeder{
		w: buf,
	}

	for _, test := range tests {
		buf.Reset()
		e := encoder{
			supportBinary: test.supportBinary,
			feeder:        f,
		}

		for _, test := range test.packets {
			fw, err := e.NextWriter(test.ft, test.pt)
			require.Nil(t, err)
			_, err = fw.Write(test.data)
			require.Nil(t, err)
			err = fw.Close()
			require.Nil(t, err)
		}

		assert.Equal(t, test.data, buf.Bytes())
	}
}

func TestEncoderBeginError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	f := &fakeWriterFeeder{
		w: buf,
	}
	e := encoder{
		supportBinary: true,
		feeder:        f,
	}

	buf.Reset()
	targetErr := newOperationError("payload", errPaused)
	f.returnError = targetErr

	_, err := e.NextWriter(frame.FrameBinary, packet.OPEN)
	assert.Equal(t, targetErr, err)
}

type errorWrite struct {
	err error
}

func (f *errorWrite) Write(p []byte) (int, error) {
	return 0, f.err
}

func TestEncoderEndError(t *testing.T) {
	f := &fakeWriterFeeder{
		w: &errorWrite{
			err: writeErr,
		},
	}
	e := encoder{
		supportBinary: true,
		feeder:        f,
	}

	targetErr := errors.New("error")

	fw, err := e.NextWriter(frame.FrameBinary, packet.OPEN)
	require.Nil(t, err)
	f.returnError = targetErr
	err = fw.Close()
	assert.Equal(t, targetErr, err)
	assert.Equal(t, f.passingErr, writeErr)
}

func TestEncoderNOOP(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		supportBinary bool
		data          []byte
	}{
		{true, []byte{0x00, 0x01, 0xff, '6'}},
		{false, []byte("1:6")},
	}

	for _, test := range tests {
		e := encoder{
			supportBinary: test.supportBinary,
		}
		assert.Equal(test.data, e.NOOP())
	}

	// NOOP should be thread-safe
	var wg sync.WaitGroup
	max := 100
	wg.Add(100)
	for i := 0; i < max; i++ {
		go func(i int) {
			defer wg.Done()
			e := encoder{
				supportBinary: i&0x1 == 0,
			}
			e.NOOP()
		}(i)
	}
	wg.Wait()
}

func BenchmarkStringEncoder(b *testing.B) {
	must := require.New(b)
	packets := []Packet{
		{frame.FrameString, packet.OPEN, []byte{}},
		{frame.FrameString, packet.MESSAGE, []byte("你好\n")},
		{frame.FrameString, packet.PING, []byte("probe")},
	}
	e := encoder{
		supportBinary: false,
		feeder: &fakeWriterFeeder{
			w: ioutil.Discard,
		},
	}

	// warm up for memory allocation
	for _, p := range packets {
		f, err := e.NextWriter(p.ft, p.pt)
		must.Nil(err)
		_, err = f.Write(p.data)
		must.Nil(err)
		err = f.Close()
		must.Nil(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range packets {
			f, _ := e.NextWriter(p.ft, p.pt)
			f.Write(p.data)
			f.Close()
		}
	}
}

func BenchmarkB64Encoder(b *testing.B) {
	must := require.New(b)
	packets := []Packet{
		{frame.FrameBinary, packet.OPEN, []byte{}},
		{frame.FrameBinary, packet.MESSAGE, []byte("你好\n")},
		{frame.FrameBinary, packet.PING, []byte("probe")},
	}
	e := encoder{
		supportBinary: false,
		feeder: &fakeWriterFeeder{
			w: ioutil.Discard,
		},
	}

	// warm up for memory allocation
	for _, p := range packets {
		f, err := e.NextWriter(p.ft, p.pt)
		must.Nil(err)
		_, err = f.Write(p.data)
		must.Nil(err)
		err = f.Close()
		must.Nil(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range packets {
			f, _ := e.NextWriter(p.ft, p.pt)
			f.Write(p.data)
			f.Close()
		}
	}
}

func BenchmarkBinaryEncoder(b *testing.B) {
	must := require.New(b)
	packets := []Packet{
		{frame.FrameString, packet.OPEN, []byte{}},
		{frame.FrameBinary, packet.MESSAGE, []byte("你好\n")},
		{frame.FrameString, packet.PING, []byte("probe")},
	}
	e := encoder{
		supportBinary: true,
		feeder: &fakeWriterFeeder{
			w: ioutil.Discard,
		},
	}

	// warm up for memory allocation
	for _, p := range packets {
		f, err := e.NextWriter(p.ft, p.pt)
		must.Nil(err)
		_, err = f.Write(p.data)
		must.Nil(err)
		err = f.Close()
		must.Nil(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range packets {
			f, _ := e.NextWriter(p.ft, p.pt)
			f.Write(p.data)
			f.Close()
		}
	}
}

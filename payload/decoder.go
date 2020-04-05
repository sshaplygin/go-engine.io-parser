package payload

import (
	"bufio"
	"encoding/base64"
	"io"
	"io/ioutil"

	"github.com/mrfoe7/go-engine.io-parser/units"
)

type byteReader interface {
	ReadByte() (byte, error)

	io.Reader
}

type readerFeeder interface {
	getReader() (io.Reader, bool, error)
	putReader(error) error
}

type decoder struct {
	feeder readerFeeder

	ft            units.FrameType
	pt            units.PacketType
	supportBinary bool
	rawReader     byteReader
	limitReader   io.LimitedReader
	b64Reader     io.Reader
}

func (d *decoder) NextReader() (units.FrameType, units.PacketType, io.ReadCloser, error) {
	if d.rawReader == nil {
		r, supportBinary, err := d.feeder.getReader()
		if err != nil {
			return units.FrameString, units.OPEN, nil, err
		}
		br, ok := r.(byteReader)
		if !ok {
			br = bufio.NewReader(r)
		}
		if err := d.setNextReader(br, supportBinary); err != nil {
			return units.FrameString, units.OPEN, nil, d.sendError(err)
		}
	}

	return d.ft, d.pt, d, nil
}

func (d *decoder) Read(p []byte) (int, error) {
	if d.b64Reader != nil {
		return d.b64Reader.Read(p)
	}
	return d.limitReader.Read(p)
}

func (d *decoder) Close() error {
	if _, err := io.Copy(ioutil.Discard, d); err != nil {
		return d.sendError(err)
	}
	err := d.setNextReader(d.rawReader, d.supportBinary)
	if err != nil {
		if err != io.EOF {
			return d.sendError(err)
		}
		d.rawReader = nil
		d.limitReader.R = nil
		d.limitReader.N = 0
		d.b64Reader = nil
		err = d.sendError(nil)
	}
	return err
}

func (d *decoder) setNextReader(r byteReader, supportBinary bool) error {
	var read func(byteReader) (units.FrameType, units.PacketType, int64, error)
	if supportBinary {
		read = d.binaryRead
	} else {
		read = d.textRead
	}

	ft, pt, l, err := read(r)
	if err != nil {
		return err
	}

	d.ft = ft
	d.pt = pt
	d.rawReader = r
	d.limitReader.R = r
	d.limitReader.N = l
	d.supportBinary = supportBinary
	if !supportBinary && ft == units.FrameBinary {
		d.b64Reader = base64.NewDecoder(base64.StdEncoding, &d.limitReader)
	} else {
		d.b64Reader = nil
	}
	return nil
}

func (d *decoder) sendError(err error) error {
	if e := d.feeder.putReader(err); e != nil {
		return e
	}
	return err
}

func (d *decoder) textRead(r byteReader) (units.FrameType, units.PacketType, int64, error) {
	l, err := readTextLen(r)
	if err != nil {
		return units.FrameString, units.OPEN, 0, err
	}

	ft := units.FrameString
	b, err := r.ReadByte()
	if err != nil {
		return units.FrameString, units.OPEN, 0, err
	}
	l--

	if b == 'b' {
		ft = units.FrameBinary
		b, err = r.ReadByte()
		if err != nil {
			return units.FrameString, units.OPEN, 0, err
		}
		l--
	}

	pt := units.ByteToPacketType(b, units.FrameString)
	return ft, pt, l, nil
}

func (d *decoder) binaryRead(r byteReader) (units.FrameType, units.PacketType, int64, error) {
	b, err := r.ReadByte()
	if err != nil {
		return units.FrameString, units.OPEN, 0, err
	}
	if b > 1 {
		return units.FrameString, units.OPEN, 0, errInvalidPayload
	}
	ft := units.ByteToFrameType(b)

	l, err := readBinaryLen(r)
	if err != nil {
		return units.FrameString, units.OPEN, 0, err
	}

	b, err = r.ReadByte()
	if err != nil {
		return units.FrameString, units.OPEN, 0, err
	}
	pt := units.ByteToPacketType(b, ft)
	l--

	return ft, pt, l, nil
}

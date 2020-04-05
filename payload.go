package go_engine_io_parser

import (
	"io"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/mrfoe7/go-engine.io-parser/frame"
	"github.com/mrfoe7/go-engine.io-parser/packet"
)

const (
	operationRead    = "read"
	operationWrite   = "write"
	operationPayload = "payload"
)

type readArg struct {
	r             io.Reader
	supportBinary bool
}

// Payload does encode or decode to payload protocol.
type Payload struct {
	closeOnce sync.Once
	err       atomic.Value

	pauser *pauser

	close      chan struct{}
	readerChan chan readArg
	readError  chan error
	writerChan chan io.Writer
	writeError chan error

	feeding      int64
	readDeadline atomic.Value
	decoder      decoder

	flushing      int64
	writeDeadline atomic.Value
	encoder       encoder
}

// New returns a new payload.
func NewPayload(supportBinary bool) *Payload {
	ret := &Payload{
		close:      make(chan struct{}),
		pauser:     newPauser(),
		readerChan: make(chan readArg),
		readError:  make(chan error),
		writerChan: make(chan io.Writer),
		writeError: make(chan error),
	}
	ret.readDeadline.Store(time.Time{})
	ret.decoder.feeder = ret
	ret.writeDeadline.Store(time.Time{})
	ret.encoder.supportBinary = supportBinary
	ret.encoder.feeder = ret
	return ret
}

// FeedIn feeds in a new reader for NextReader.
// Multi-FeedIn needs be called sync.
//
// If Close called when FeedIn, it returns io.EOF.
// If have Pause-ed when FeedIn, it returns ErrPaused.
// If NextReader has timeout, it returns ErrTimeout.
// If read error while FeedIn, it returns read error.
func (p *Payload) FeedIn(r io.Reader, supportBinary bool) error {
	select {
	case <-p.close:
		return p.load()
	default:
	}

	if !atomic.CompareAndSwapInt64(&p.feeding, 0, 1) {
		return newOperationError(operationRead, errOverlap)
	}
	defer atomic.StoreInt64(&p.feeding, 0)
	if ok := p.pauser.Working(); !ok {
		return newOperationError(operationPayload, errPaused)
	}
	defer p.pauser.Done()

	for {
		after, ok := p.readTimeout()
		if !ok {
			return p.Store("read", errTimeout)
		}
		select {
		case <-p.close:
			return p.load()
		case <-after:
			// it may changed during wait, need check again
			continue
		case p.readerChan <- readArg{
			r:             r,
			supportBinary: supportBinary,
		}:
		}
		break
	}

	for {
		after, ok := p.readTimeout()
		if !ok {
			return p.Store("read", errTimeout)
		}
		select {
		case <-after:
			// it may changed during wait, need check again
			continue
		case err := <-p.readError:
			return p.Store(operationRead, err)
		}
	}
}

// FlushOut write data from NextWriter.
// FlushOut needs be called sync.
//
// If Close called when Flushout,  it return io.EOF.
// If Pause called when Flushout, it flushs out a NOOP message and return
// nil.
// If NextWriter has timeout, it returns ErrTimeout.
// If write error while FlushOut, it returns write error.
func (p *Payload) FlushOut(w io.Writer) error {
	select {
	case <-p.close:
		return p.load()
	default:
	}
	if !atomic.CompareAndSwapInt64(&p.flushing, 0, 1) {
		return newOperationError(operationWrite, errOverlap)
	}
	defer atomic.StoreInt64(&p.flushing, 0)

	if ok := p.pauser.Working(); !ok {
		_, err := w.Write(p.encoder.NOOP())
		return err
	}
	defer p.pauser.Done()

	for {
		after, ok := p.writeTimeout()
		if !ok {
			return p.Store("write", errTimeout)
		}
		select {
		case <-p.close:
			return p.load()
		case <-after:
			continue
		case <-p.pauser.PausingTrigger():
			_, err := w.Write(p.encoder.NOOP())
			return err
		case p.writerChan <- w:
		}
		break
	}

	for {
		after, ok := p.writeTimeout()
		if !ok {
			return p.Store(operationWrite, errTimeout)
		}
		select {
		case <-after:
		// it may changed during wait, need check again
		case err := <-p.writeError:
			return p.Store(operationWrite, err)
		}
	}
}

// NextReader returns a reader for next frame.
// NextReader and SetReadDeadline needs be called sync.
//
// If Close called when NextReader,  it return io.EOF.
// Pause doesn't effect to NextReader. NextReader should wait till resumed
// and next FeedIn.
func (p *Payload) NextReader() (frame.FrameType, packet.PacketType, io.ReadCloser, error) {
	ft, pt, r, err := p.decoder.NextReader()
	return ft, pt, r, err
}

// SetReadDeadline sets next reader deadline.
// NextReader and SetReadDeadline needs be called sync.
// NextReader will wait a FeedIn call, then it returns ReadCloser which
// decodes packet from FeedIn's Reader.
//
// If Close called when SetReadDeadline,  it return io.EOF.
// If beyond the time set by SetReadDeadline, it returns ErrTimeout.
// Pause doesn't effect to SetReadDeadline.
func (p *Payload) SetReadDeadline(t time.Time) error {
	p.readDeadline.Store(t)
	return nil
}

// NextWriter returns a writer for next frame.
// NextWriter and SetWriterDeadline needs be called sync.
// NextWriter will wait a FlushOut call, then it returns WriteCloser which
// encode package to FlushOut's Writer.
//
// If Close called when NextWriter,  it returns io.EOF.
// If beyond the time set by SetWriteDeadline, it returns ErrTimeout.
// If Pause called when NextWriter, it returns ErrPaused.
func (p *Payload) NextWriter(ft frame.FrameType, pt packet.PacketType) (io.WriteCloser, error) {
	return p.encoder.NextWriter(ft, pt)
}

// SetWriteDeadline sets next writer deadline.
// NextWriter and SetWriteDeadline needs be called sync.
//
// If Close called when SetWriteDeadline,  it return io.EOF.
// Pause doesn't effect to SetWriteDeadline.
func (p *Payload) SetWriteDeadline(t time.Time) error {
	p.writeDeadline.Store(t)
	return nil
}

// Pause pauses the payload. It will wait all reader and writer closed which
// created from NextReader or NextWriter.
// It can call in multi-goroutine.
func (p *Payload) Pause() {
	p.pauser.Pause()
}

// Resume resumes the payload.
// It can call in multi-goroutine.
func (p *Payload) Resume() {
	p.pauser.Resume()
}

// Close closes the payload.
// It can call in multi-goroutine.
func (p *Payload) Close() error {
	p.closeOnce.Do(func() {
		close(p.close)
	})
	return nil
}

// Store stores a error in payload, and block all other request.
func (p *Payload) Store(operation string, err error) error {
	old := p.err.Load()
	//ref:
	if old == nil {
		if err == io.EOF {
			return err
		}
		op := newOperationError(operation, err)
		p.err.Store(op)
		return op
	}
	return old.(error)
}

func (p *Payload) readTimeout() (<-chan time.Time, bool) {
	deadline := p.readDeadline.Load().(time.Time)
	wait := deadline.Sub(time.Now())
	if deadline.IsZero() {
		// wait for every
		wait = math.MaxInt64
	}
	if wait <= 0 {
		return nil, false
	}
	return time.After(wait), true
}

func (p *Payload) writeTimeout() (<-chan time.Time, bool) {
	deadline := p.writeDeadline.Load().(time.Time)
	wait := deadline.Sub(time.Now())
	if deadline.IsZero() {
		// wait for every
		wait = math.MaxInt64
	}
	if wait <= 0 {
		return nil, false
	}
	return time.After(wait), true
}

func (p *Payload) getReader() (io.Reader, bool, error) {
	select {
	case <-p.close:
		return nil, false, p.load()
	default:
	}

	if ok := p.pauser.Working(); !ok {
		return nil, false, newOperationError(operationPayload, errPaused)
	}
	p.pauser.Done()

	for {
		after, ok := p.readTimeout()
		if !ok {
			return nil, false, p.Store(operationRead, errTimeout)
		}
		select {
		case <-p.close:
			return nil, false, p.load()
		case <-p.pauser.PausedTrigger():
			return nil, false, newOperationError(operationPayload, errPaused)
		case <-after:
			continue
		case arg := <-p.readerChan:
			return arg.r, arg.supportBinary, nil
		}
	}
}

func (p *Payload) putReader(err error) error {
	select {
	case <-p.close:
		return p.load()
	default:
	}
	for {
		after, ok := p.readTimeout()
		if !ok {
			return p.Store(operationRead, errTimeout)
		}
		select {
		case <-p.close:
			return p.load()
		case <-after:
			continue
		case p.readError <- err:
		}
		return nil
	}
}

func (p *Payload) getWriter() (io.Writer, error) {
	select {
	case <-p.close:
		return nil, p.load()
	default:
	}

	if ok := p.pauser.Working(); !ok {
		return nil, newOperationError(operationPayload, errPaused)
	}
	p.pauser.Done()

	for {
		after, ok := p.writeTimeout()
		if !ok {
			return nil, p.Store("write", errTimeout)
		}
		select {
		case <-p.close:
			return nil, p.load()
		case <-p.pauser.PausedTrigger():
			return nil, newOperationError(operationPayload, errPaused)
		case <-after:
			continue
		case w := <-p.writerChan:
			return w, nil
		}
	}
}

func (p *Payload) putWriter(err error) error {
	select {
	case <-p.close:
		return p.load()
	default:
	}
	for {
		after, ok := p.writeTimeout()
		if !ok {
			return p.Store(operationWrite, errTimeout)
		}
		ret := p.Store(operationWrite, err)
		select {
		case <-p.close:
			return p.load()
		case <-after:
			continue
		case p.writeError <- err:
			return ret
		}
	}
}

func (p *Payload) load() error {
	ret := p.err.Load()
	if ret == nil {
		return io.EOF
	}
	return ret.(error)
}

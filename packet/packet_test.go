package packet

import (
	"github.com/go-engine.io-parser/frame"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketType(t *testing.T) {
	at := assert.New(t)
	tests := []struct {
		b         byte
		frameType frame.FrameType
		typ       PacketType
		strbyte   byte
		binbyte   byte
		str       string
	}{
		{0, frame.FrameBinary, OPEN, '0', 0, "open"},
		{1, frame.FrameBinary, CLOSE, '1', 1, "close"},
		{2, frame.FrameBinary, PING, '2', 2, "ping"},
		{3, frame.FrameBinary, PONG, '3', 3, "pong"},
		{4, frame.FrameBinary, MESSAGE, '4', 4, "message"},
		{5, frame.FrameBinary, UPGRADE, '5', 5, "upgrade"},
		{6, frame.FrameBinary, NOOP, '6', 6, "noop"},

		{'0', frame.FrameString, OPEN, '0', 0, "open"},
		{'1', frame.FrameString, CLOSE, '1', 1, "close"},
		{'2', frame.FrameString, PING, '2', 2, "ping"},
		{'3', frame.FrameString, PONG, '3', 3, "pong"},
		{'4', frame.FrameString, MESSAGE, '4', 4, "message"},
		{'5', frame.FrameString, UPGRADE, '5', 5, "upgrade"},
		{'6', frame.FrameString, NOOP, '6', 6, "noop"},
	}

	for _, test := range tests {
		typ := ByteToPacketType(test.b, test.frameType)
		at.Equal(test.typ, typ)
		at.Equal(test.strbyte, typ.StringByte())
		at.Equal(test.binbyte, typ.BinaryByte())
		at.Equal(test.str, typ.String())
		at.Equal(test.str, PacketType(typ).String())
	}
}

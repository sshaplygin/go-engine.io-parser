package go_engine_io_parser

import (
	"github.com/mrfoe7/go-engine.io-parser/frame"
	"github.com/mrfoe7/go-engine.io-parser/packet"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TesPacketTypeCase struct {
	b          byte
	frameType  frame.FrameType
	packetType packet.PacketType

	strByte byte
	binByte byte
	str     string
}

func Test_PacketType(t *testing.T) {
	tests := []TesPacketTypeCase{
		{0, frame.FrameBinary, packet.OPEN, '0', 0, "open"},
		{1, frame.FrameBinary, packet.CLOSE, '1', 1, "close"},
		{2, frame.FrameBinary, packet.PING, '2', 2, "ping"},
		{3, frame.FrameBinary, packet.PONG, '3', 3, "pong"},
		{4, frame.FrameBinary, packet.MESSAGE, '4', 4, "message"},
		{5, frame.FrameBinary, packet.UPGRADE, '5', 5, "upgrade"},
		{6, frame.FrameBinary, packet.NOOP, '6', 6, "noop"},

		{'0', frame.FrameString, packet.OPEN, '0', 0, "open"},
		{'1', frame.FrameString, packet.CLOSE, '1', 1, "close"},
		{'2', frame.FrameString, packet.PING, '2', 2, "ping"},
		{'3', frame.FrameString, packet.PONG, '3', 3, "pong"},
		{'4', frame.FrameString, packet.MESSAGE, '4', 4, "message"},
		{'5', frame.FrameString, packet.UPGRADE, '5', 5, "upgrade"},
		{'6', frame.FrameString, packet.NOOP, '6', 6, "noop"},
	}

	for _, test := range tests {
		typ := byteToPacketType(test.b, test.frameType)
		assert.Equal(t, test.packetType, typ)
		assert.Equal(t, test.strByte, typ.StringByte())
		assert.Equal(t, test.binByte, typ.BinaryByte())
		assert.Equal(t, test.str, packet.PacketType(typ).String())
		assert.Equal(t, test.str, typ.String())
	}
}

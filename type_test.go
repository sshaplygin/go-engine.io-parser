package go_engine_io_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TesPacketTypeCase struct {
	b          byte
	frameType  FrameType
	packetType Type

	strByte byte
	binByte byte
	str     string
}

func Test_PacketType(t *testing.T) {
	tests := []TesPacketTypeCase{
		{0, FrameBinary, Open, '0', 0, "Open"},
		{1, FrameBinary, CLOSE, '1', 1, "close"},
		{2, FrameBinary, Ping, '2', 2, "Ping"},
		{3, FrameBinary, PONG, '3', 3, "pong"},
		{4, FrameBinary, Message, '4', 4, "Message"},
		{5, FrameBinary, UPGRADE, '5', 5, "upgrade"},
		{6, FrameBinary, NOOP, '6', 6, "noop"},

		{'0', FrameString, Open, '0', 0, "Open"},
		{'1', FrameString, CLOSE, '1', 1, "close"},
		{'2', FrameString, Ping, '2', 2, "Ping"},
		{'3', FrameString, PONG, '3', 3, "pong"},
		{'4', FrameString, Message, '4', 4, "Message"},
		{'5', FrameString, UPGRADE, '5', 5, "upgrade"},
		{'6', FrameString, NOOP, '6', 6, "noop"},
	}

	for _, test := range tests {
		typ := byteToPacketType(test.b, test.frameType)
		assert.Equal(t, test.packetType, typ)
		assert.Equal(t, test.strByte, typ.StringByte())
		assert.Equal(t, test.binByte, typ.BinaryByte())
		assert.Equal(t, test.str, Type(typ).String())
		assert.Equal(t, test.str, typ.String())
	}
}

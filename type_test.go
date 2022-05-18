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
		{frameType: FrameBinary, packetType: Open, strByte: '0', str: "Open"},
		{b: 1, frameType: FrameBinary, packetType: Close, strByte: '1', binByte: 1, str: "Close"},
		{b: 2, frameType: FrameBinary, packetType: Ping, strByte: '2', binByte: 2, str: "Ping"},
		{b: 3, frameType: FrameBinary, packetType: Pong, strByte: '3', binByte: 3, str: "Pong"},
		{b: 4, frameType: FrameBinary, packetType: Message, strByte: '4', binByte: 4, str: "Message"},
		{b: 5, frameType: FrameBinary, packetType: Upgrade, strByte: '5', binByte: 5, str: "Upgrade"},
		{b: 6, frameType: FrameBinary, packetType: Noop, strByte: '6', binByte: 6, str: "Noop"},

		{b: '0', frameType: FrameString, packetType: Open, strByte: '0', str: "Open"},
		{b: '1', frameType: FrameString, packetType: Close, strByte: '1', binByte: 1, str: "Close"},
		{b: '2', frameType: FrameString, packetType: Ping, strByte: '2', binByte: 2, str: "Ping"},
		{b: '3', frameType: FrameString, packetType: Pong, strByte: '3', binByte: 3, str: "Pong"},
		{b: '4', frameType: FrameString, packetType: Message, strByte: '4', binByte: 4, str: "Message"},
		{b: '5', frameType: FrameString, packetType: Upgrade, strByte: '5', binByte: 5, str: "Upgrade"},
		{b: '6', frameType: FrameString, packetType: Noop, strByte: '6', binByte: 6, str: "Noop"},
	}

	for _, test := range tests {
		typ := byteToPacketType(test.b, test.frameType)
		assert.Equal(t, test.packetType, typ)
		assert.Equal(t, test.strByte, typ.StringByte())
		assert.Equal(t, test.binByte, typ.Byte())
		assert.Equal(t, test.str, typ.String())
		assert.Equal(t, test.str, typ.String())
	}
}

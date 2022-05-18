package go_engine_io_parser

type Packet struct {
	ft   FrameType
	pt   Type
	data []byte
}

var tests = []struct {
	supportBinary bool
	data          []byte
	packets       []Packet
}{
	{supportBinary: true,
		data: []byte{0x00, 0x01, 0xff, '0'},
		packets: []Packet{
			{
				ft:   FrameString,
				pt:   Open,
				data: []byte{},
			},
		},
	},
	{
		supportBinary: true,
		data:          []byte{0x00, 0x01, 0x03, 0xff, '4', 'h', 'e', 'l', 'l', 'o', ' ', 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd},
		packets: []Packet{
			{
				ft: FrameString, pt: Message, data: []byte("hello 你好"),
			},
		},
	},
	{true, []byte{0x01, 0x01, 0x03, 0xff, 0x04, 'h', 'e', 'l', 'l', 'o', ' ', 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd}, []Packet{
		{FrameBinary, Message, []byte("hello 你好")},
	}},
	{true, []byte{
		0x01, 0x07, 0xff, 0x04, 'h', 'e', 'l', 'l', 'o', '\n',
		0x00, 0x08, 0xff, '4', 0xe4, 0xbd, 0xa0, 0xe5, 0xa5, 0xbd, '\n',
		0x00, 0x06, 0xff, '2', 'p', 'r', 'o', 'b', 'e',
	}, []Packet{
		{FrameBinary, Message, []byte("hello\n")},
		{FrameString, Message, []byte("你好\n")},
		{FrameString, Ping, []byte("probe")},
	}},

	{false, []byte("1:0"), []Packet{
		{FrameString, Open, []byte{}},
	}},
	{false, []byte("13:4hello 你好"), []Packet{
		{FrameString, Message, []byte("hello 你好")},
	}},
	{false, []byte("18:b4aGVsbG8g5L2g5aW9"), []Packet{
		{FrameBinary, Message, []byte("hello 你好")},
	}},
	{false, []byte("10:b4aGVsbG8K8:4你好\n6:2probe"), []Packet{
		{FrameBinary, Message, []byte("hello\n")},
		{FrameString, Message, []byte("你好\n")},
		{FrameString, Ping, []byte("probe")},
	}},
}

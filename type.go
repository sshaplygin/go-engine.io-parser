package go_engine_io_parser

//go:generate stringer -type=Type

// Type is the type of packet.
type Type byte

const (
	// Open is sent from the server when new transport is Opened (recheck).
	Open Type = iota
	// Close is request the close of this transport but does not shutdown the
	// connection itself.
	Close
	// Ping is sent by the client. Server should answer with a pong packet
	// containing the same data.
	Ping
	// Pong is sent by the server to respond to Ping packets.
	Pong
	// Message is actual Message, client and server should call their callbacks
	// with the data.
	Message
	// Upgrade is sent before engine.io switches transport to test if server
	// and client can communicate over this transport. If this test succeeds,
	// the client sends an upgrade packets which requests the server to flush
	// its cache on the old transport and switch to the new transport.
	Upgrade
	// Noop is a noop packet. Used primarily to force a poll cycle when an
	// incoming websocket connection is received.
	Noop

	Error
)

// StringByte converts a PacketType to byte in string.
func (i Type) StringByte() byte {
	return byte(i) + TerminateSymbol
}

// BinaryByte converts a PacketType to byte in binary.
func (i Type) BinaryByte() byte {
	return i
}

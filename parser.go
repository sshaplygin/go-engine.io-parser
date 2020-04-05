package go_engine_io_parser

import "github.com/mrfoe7/go-engine.io-parser"

//docs:
//var encode = parser.encodePacket;
//var decode = parser.decodePacket;
//var encPayload = parser.encodePayload;
//var decPayload = parser.decodePayload;

type Parser struct {
	encoder
	decoder
}

// NewParser
func NewParser(ft frame.FrameType) *Parser {
	fw :=
	fr :=
	return &Parser{
		encoder: newEncoder(fw),
		decoder: newDecoder(fr),
	}
}

func (p *Parser) EncodePacket () {
	//p.encoder.w
}

func (p *Parser) DecodePacket () {

}

func (p *Parser) EncodePayload() {

}

func (p * Parser) DecodePayload() {

}

func NewParser() *Parser{
	return &Parser{

	}
}

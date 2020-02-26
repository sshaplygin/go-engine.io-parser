package go_engine_io_parser

//docs:
//var encode = parser.encodePacket;
//var decode = parser.decodePacket;
//var encPayload = parser.encodePayload;
//var decPayload = parser.decodePayload;

type Parser struct {
	encoder
	decoder
}

func NewParser() *Parser {
	return
}

func (p *Parser) EncodePacket () {

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

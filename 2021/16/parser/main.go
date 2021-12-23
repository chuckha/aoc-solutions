package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
packets -> packet packets
packets -> nil?
packet -> version litType litValue
packet -> version opType opValue
version -> 000,001,010,011,100,101,110,111
litValue -> litValHeader litData litValue
litValue -> litValHeader litLastData
litValHeader -> litData, lastLitData
litData -> 0000,0001,...,1111
lastLitData -> 0000,0001,...,1111
opValue -> numberOfPackets packets
opValue -> numberOfBytes packets
type -> opType, litType
opType -> 000,001,010,011,101,110,111
litType -> 100

packet : version type {literal packet|op packet}
op packet: lengthID
lengthID : 0, 1
*/

// read x number of packets

// read a list of packets
//  read a single packet
// read the version
// read the type
//   if it's a lit, read v,t,lit value
//   if it's an op, read v,t,length id type (0/1)
//    	if it's type 0, read x bytes
//      if it's type 1, read x packets

// OP(1:2){OP(0:22){}}
//                                      22
// 011 000 1 00000000010 000 000 0 000000000010110 (000 100 01010 101 100 01011) 001000100000000010000100011000111000110100

// []*P{[]*P{}}
/*
P{
	Data: 00000000000000000101100001000101010110001011001000100000000010000100011000111000110100
	Parsed: []*P{
		P: {
			Data: 000 000 0 000000000010110 0001000101010110001011
			Parsed: []*P{
				P: {
					Data: 0001000101010110001011
					Parsed: []*P {
						{
							Data: 00010001010
							Lit: 1010
						},
						{
							Data: 10110001011
							Lit: 1011
						}
					}
				}
			}
		},
	},
}

*/

const (
	version          = "version"
	typeID           = "typeID"
	literalValue     = "literalValue"
	fifteenBitNumber = "fifteenBitNumber"
	elevenBitNumber  = "elevenBitNumber"
	lengthTypeID     = "lengthTypeID"
	subpackets       = "subpackets"

	litValueTypeID = 4
)

func binToDec(in string) int {
	out, _ := strconv.ParseInt(in, 2, 64)
	return int(out)
}

/*
* 620080001611562C8802118E34 represents an operator packet (version 3) which
contains two sub-packets; each sub-packet is an operator packet that contains
two literal values. This packet has a version sum of 12.
*/

func main() {
	lines := internal.ReadInput()
	input := hexToBin(lines[0])
	tokens := make(chan item)
	l := &lex{
		input:  input,
		tokens: tokens,
	}
	go l.run()
	tkns := make([]item, 0)
	for token := range tokens {
		tkns = append(tkns, token)
	}
	prsr := &parser{tkns, 0}
	packets := prsr.parsePackets()
	fmt.Println(packets[0])
	fmt.Println(packets[0].value())
}

type parser struct {
	tokens []item
	pos    int
}

func (p *parser) next() (item, bool) {
	if p.pos >= len(p.tokens) {
		return item{}, false
	}
	out := p.tokens[p.pos]
	p.pos++
	return out, true
}
func (p *parser) peek() (item, bool) {
	if p.pos > len(p.tokens)-1 {
		return item{}, false
	}
	return p.tokens[p.pos+1], true
}

func (p *parser) parsePackets() []*packet {
	if _, ok := p.peek(); !ok {
		return nil
	}
	out := []*packet{}
	cur := &packet{}
	version, ok := p.next()
	if !ok {
		return out
	}
	vint := binToDec(version.value)
	cur.Version = vint
	packetType, ok := p.next()
	if !ok {
		panic("weird place to end my dude")
	}
	ptint := binToDec(packetType.value)
	cur.PacketType = ptint
	if packetType.value == "100" {
		literalValue, ok := p.next()
		if !ok {
			panic("ended in the middle of a litearl value packet")
		}
		cur.LitValue = litBinVal(literalValue.value)
		out = append(out, cur)
		out = append(out, p.parsePackets()...)
		return out
	}
	lengthTypeID, ok := p.next()
	if !ok {
		panic("ended in the middle of a length type id")
	}
	cur.LengthTypeID = binToDec(lengthTypeID.value)
	if lengthTypeID.value == "0" {
		subpacketLength, ok := p.next()
		if !ok {
			panic("ended in the middle of a subpacket length")
		}
		cur.SubpacketLength = binToDec(subpacketLength.value)
		cur.Subpackets = p.parseBits(cur.SubpacketLength)
		out = append(out, cur)
		out = append(out, p.parsePackets()...)
		return out
	}
	numberOfPackets, ok := p.next()
	if !ok {
		panic("ended in the middle of number of packets")
	}
	cur.NumberOfPackets = binToDec(numberOfPackets.value)
	cur.Subpackets = p.parseNumPackets(cur.NumberOfPackets)
	out = append(out, cur)
	out = append(out, p.parsePackets()...)
	return out
}

func (p *parser) getPacket() *packet {
	cur := &packet{}
	version, ok := p.next()
	if !ok {
		return nil
	}
	vint := binToDec(version.value)
	cur.Version = vint
	cur.bits += version.bits
	packetType, ok := p.next()
	if !ok {
		panic("ended during packet type")
	}
	ptint := binToDec(packetType.value)
	cur.PacketType = ptint
	cur.bits += packetType.bits
	if packetType.value == "100" {
		literalValue, ok := p.next()
		if !ok {
			panic("ended during literal value")
		}
		cur.LitValue = litBinVal(literalValue.value)
		cur.bits += literalValue.bits
		return cur
	}
	lengthTypeID, ok := p.next()
	if !ok {
		panic("ended during length type id")
	}
	cur.LengthTypeID = binToDec(lengthTypeID.value)
	cur.bits += lengthTypeID.bits
	if lengthTypeID.value == "0" {
		subpacketLength, ok := p.next()
		if !ok {
			panic("ended during subpacket length")
		}
		cur.SubpacketLength = binToDec(subpacketLength.value)
		cur.bits += subpacketLength.bits
		cur.Subpackets = p.parseBits(cur.SubpacketLength)
		return cur
	}
	numberOfPackets, ok := p.next()
	if !ok {
		panic("ended during number of packets")
	}
	cur.NumberOfPackets = binToDec(numberOfPackets.value)
	cur.bits += numberOfPackets.bits
	cur.Subpackets = p.parseNumPackets(cur.NumberOfPackets)
	return cur
}

func (p *parser) parseBits(numBits int) []*packet {
	out := []*packet{}
	for numBits > 0 {
		pkt := p.getPacket()
		if pkt == nil {
			return out
		}
		numBits -= pkt.bits
		out = append(out, pkt)
	}
	return out
}

func (p *parser) parseNumPackets(num int) []*packet {
	out := make([]*packet, 0)
	for i := 0; i < num; i++ {
		out = append(out, p.getPacket())
	}
	return out
}

type packet struct {
	Version    int
	PacketType int
	bits       int

	LitValue int

	LengthTypeID int
	// 15 bit number // 0 type id
	SubpacketLength int
	// 11 bit number // 1 type id
	NumberOfPackets int
	Subpackets      []*packet
}

func (p *packet) value() int {
	if p == nil {
		return 0
	}
	if p.PacketType == 4 {
		return p.LitValue
	}
	if p.PacketType == 0 {
		sum := 0
		for _, sp := range p.Subpackets {
			sum += sp.value()
		}
		return sum
	}
	if p.PacketType == 1 {
		prod := 1
		for _, sp := range p.Subpackets {
			prod *= sp.value()
		}
		return prod
	}
	if p.PacketType == 2 {
		min := 99999999999
		for _, sp := range p.Subpackets {
			val := sp.value()
			if val < min {
				min = val
			}
		}
		return min
	}
	if p.PacketType == 3 {
		max := 0
		for _, sp := range p.Subpackets {
			val := sp.value()
			if val > max {
				max = val
			}
		}
		return max
	}
	if p.PacketType == 5 {
		if p.Subpackets[0].value() > p.Subpackets[1].value() {
			return 1
		}
		return 0
	}
	if p.PacketType == 6 {
		if p.Subpackets[0].value() < p.Subpackets[1].value() {
			return 1
		}
		return 0
	}
	if p.PacketType == 7 {
		if p.Subpackets[0].value() == p.Subpackets[1].value() {
			return 1
		}
		return 0
	}
	panic("yoikes")
}

func (p *packet) String() string {
	out, _ := json.MarshalIndent(p, "", "    ")
	return string(out)
	// if p.packetType == 4 {
	// 	return fmt.Sprintf("Literal Packet: %d", p.litValue)
	// }
	// out := fmt.Sprintf("Operator packet:\n")
	// for _, pkt := range p.subpackets {
	// 	out += fmt.Sprintf("\t%v\n", pkt)
	// }
	// return out
}
func (p *packet) versionSum() int {
	if p == nil {
		return 0
	}
	if len(p.Subpackets) == 0 {
		return p.Version
	}
	sum := 0
	for _, sp := range p.Subpackets {
		sum += sp.versionSum()
	}
	return p.Version + sum
}

func hexToBin(in string) string {
	out := ""
	t := translation()
	for _, c := range in {
		out += t[string(c)]
	}
	return out
}

func translation() map[string]string {
	return map[string]string{
		"0": "0000",
		"1": "0001",
		"2": "0010",
		"3": "0011",
		"4": "0100",
		"5": "0101",
		"6": "0110",
		"7": "0111",
		"8": "1000",
		"9": "1001",
		"A": "1010",
		"B": "1011",
		"C": "1100",
		"D": "1101",
		"E": "1110",
		"F": "1111",
	}
}

type item struct {
	kind  string
	value string
	bits  int
}

func (i item) String() string {
	return fmt.Sprintf("item: kind: %v, value: %v, bits: %v", i.kind, i.value, i.bits)
}

type lex struct {
	input  string
	pos    int
	start  int
	tokens chan item
}

func (l *lex) run() {
	for state := versionFn; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lex) next() string {
	if l.pos >= len(l.input) {
		return "EOF"
	}
	bit := string(l.input[l.pos])
	l.pos++
	return bit
}

func (l *lex) peek() string {
	c := l.next()
	l.pos -= 1
	return c
}
func (l *lex) ignore() {
	l.start = l.pos
}
func (l *lex) send(kind string) {
	l.tokens <- item{
		kind:  kind,
		value: l.input[l.start:l.pos],
		bits:  l.pos - l.start,
	}
	l.start = l.pos
}

func (l *lex) errorf(error string) stateFn {
	l.tokens <- item{
		kind:  "error",
		value: error,
	}
	return nil
}

type stateFn func(*lex) stateFn

// 001 010 1 0110001011001000100000000010000100011000111000110100
func versionFn(l *lex) stateFn {
	for i := 0; i < 3; i++ {
		if l.next() == "EOF" {
			return l.errorf("found an early end in a version")
		}
	}
	l.send(version)
	return packetType
}

func packetType(l *lex) stateFn {
	for i := 0; i < 3; i++ {
		if l.next() == "EOF" {
			return l.errorf("bad input, end came after a version")
		}
	}
	val := l.input[l.start:l.pos]
	l.send(typeID)
	switch val {
	case "100":
		return litValue
	default:
		return operatorValue
	}
}

func litValue(l *lex) stateFn {
	// read 5 at a time until you find a 5 that starts with 0
	for {
		stop := false
		header := l.next()
		if header == "0" || header == "EOF" {
			stop = true
		}
		for j := 0; j < 4; j++ {
			l.next()
		}
		if stop {
			break
		}
	}
	l.send(literalValue)
	// some weird 0s at the end
	if len(l.input[l.pos:]) < 8 {
		return nil
	}
	// start over
	return versionFn
}

func operatorValue(l *lex) stateFn {
	c := l.next()
	if c == "EOF" {
		return l.errorf("should be length-type, got end")
	}
	l.send(lengthTypeID)
	if c == "1" {
		return elevenBitParse
	}
	return fifteenBitParse
}

func fifteenBitParse(l *lex) stateFn {
	for i := 0; i < 15; i++ {
		if l.next() == "EOF" {
			return l.errorf(fmt.Sprintf("found end when expecting a %d bit number", 15))
		}
	}
	l.send(fifteenBitNumber)
	return versionFn
}

func elevenBitParse(l *lex) stateFn {
	for i := 0; i < 11; i++ {
		if l.next() == "EOF" {
			return l.errorf(fmt.Sprintf("found end when expecting a %d bit number", 11))
		}
	}
	l.send(elevenBitNumber)
	return versionFn
}

func litBinVal(in string) int {
	out := ""
	i := 0
	for len(in) >= 5 {
		out += in[i*5+1 : (i+1)*5]
		in = in[(i+1)*5:]
	}
	return binToDec(out)
}

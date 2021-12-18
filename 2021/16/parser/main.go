package main

import (
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

//                                      22
// 011 000 1 00000000010 000 000 0 000000000010110 0001000101010110001011001000100000000010000100011000111000110100

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
	fmt.Println(hexToBin(lines[0]))
	l := &lex{
		input: hexToBin(lines[0]),
	}
	tlps := l.run()
	parseSubpacket(tlps[0])
	fmt.Println(tlps)
}

func nextPacket(data string) string {
	packetType := data[2:5]
	if packetType == "010" {
		return "literal"
	}
	return "op"
}

func parseSubpacket(tlp *packet) {
	fmt.Println(tlp)
	fmt.Println("parsing new data", tlp.subpacketData)
	l := &lex{input: tlp.subpacketData}
	tlps := l.run()
	fmt.Println("me", tlps)
	if len(tlps) == 0 {
		fmt.Println("Failed to parse a packet?")
		panic("bad news")
	}
	tlp.subpackets = tlps
	for _, tlp := range tlp.subpackets {
		if tlp.subpacketData != "" {
			parseSubpacket(tlp)
		}
	}
}

type packet struct {
	version int
	typeID  int

	// literal packet only
	literalValue int

	// operator packet only
	lengthTypeID   int
	fifteenBitData int
	elevenBitData  int
	subpacketData  string
	subpackets     []*packet
}

func (p *packet) String() string {
	if p.typeID == 4 {
		return fmt.Sprintf("LVP: %d", p.literalValue)
	}
	if p.typeID != 4 {
		return fmt.Sprintf("OP (%s): {%v}", p.subpacketData, p.subpackets)
	}
	return ""
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
	input   string
	pos     int
	start   int
	packet  *packet
	packets []*packet
}

func (l *lex) run() []*packet {
	for state := versionFn; state != nil; {
		state = state(l)
	}
	return l.packets
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
func (l *lex) value() string {
	return l.input[l.start:l.pos]
}

func (l *lex) errorf(error string) stateFn {
	fmt.Println(error)
	return nil
}

type stateFn func(*lex) stateFn

// 001 010 1 0110001011001000100000000010000100011000111000110100
func versionFn(l *lex) stateFn {
	fmt.Println("parsing version", l.input[l.pos:])
	l.packet = &packet{}
	for i := 0; i < 3; i++ {
		if l.next() == "EOF" {
			return l.errorf("found an early end in a version")
		}
	}
	v, _ := strconv.ParseInt(l.value(), 2, 64)
	l.packet.version = int(v)
	l.ignore()
	//	l.send(version)
	return packetType
}

func packetType(l *lex) stateFn {
	fmt.Println("parsing packet type")
	for i := 0; i < 3; i++ {
		if l.next() == "EOF" {
			return l.errorf("bad input, end came after a version")
		}
	}
	pact, _ := strconv.ParseInt(l.value(), 2, 64)
	l.packet.typeID = int(pact)
	l.ignore()
	//	l.send(typeID)
	switch pact {
	case litValueTypeID:
		return litValue
	default:
		return operatorValue
	}
}

func litValue(l *lex) stateFn {
	fmt.Println("parsing lit val")
	// read 5 at a time until you find a 5 that starts with 0
	num := ""
	for {
		stop := false
		header := l.next()
		if header == "0" || header == "EOF" {
			stop = true
		}
		l.ignore()
		for j := 0; j < 4; j++ {
			l.next()
		}
		num += l.value()
		if stop {
			break
		}
	}
	litVal, _ := strconv.ParseInt(num, 2, 64)
	l.packet.literalValue = int(litVal)
	l.ignore()
	l.packets = append(l.packets, l.packet)
	if len(l.input[l.pos:]) < 8 {
		return nil
	}
	return versionFn
}

func operatorValue(l *lex) stateFn {
	fmt.Println("parsing op value")
	fmt.Println(l.input[l.start:l.pos], l.start, l.pos)
	if l.next() == "EOF" {
		return l.errorf("should be length-type, got end")
	}
	//	l.send(lengthTypeID)
	fmt.Println("value now", l.value())
	if l.value() == "1" {
		l.ignore()
		return elevenBitParse
	}
	l.ignore()
	return fifteenBitParse
}

func fifteenBitParse(l *lex) stateFn {
	fmt.Println("parsing fifteen bit parse")
	for i := 0; i < 15; i++ {
		if l.next() == "EOF" {
			return l.errorf(fmt.Sprintf("found end when expecting a %d bit number", 15))
		}
	}
	//	l.send(fifteenBitNumber)
	btz, _ := strconv.ParseInt(l.value(), 2, 64)
	l.packet.fifteenBitData = int(btz)
	fmt.Println(l.packet, l.packet.fifteenBitData)
	l.ignore()
	return parseBits(int(btz))
}

func elevenBitParse(l *lex) stateFn {
	fmt.Println("parsing 11 bit parse")
	for i := 0; i < 11; i++ {
		if l.next() == "EOF" {
			return l.errorf(fmt.Sprintf("found end when expecting a %d bit number", 11))
		}
	}
	//	l.send(elevenBitNumber)
	numSubPacket, _ := strconv.ParseInt(l.value(), 2, 64)
	l.packet.elevenBitData = int(numSubPacket)
	l.ignore()
	return parseSubPackets(int(numSubPacket))
}

func parseBits(n int) stateFn {
	fmt.Println("parse bits")
	return func(l *lex) stateFn {
		for i := 0; i < n; i++ {
			l.next()
		}
		l.packet.subpacketData = l.value()
		l.ignore()
		l.packets = append(l.packets, l.packet)
		fmt.Println("parsed bits", l.packet.subpacketData)
		return versionFn
	}
}

func parseSubPackets(n int) stateFn {
	fmt.Println("parsing sub packets")
	return func(l *lex) stateFn {
		fmt.Println("parsing", n, "packets")
		for i := 0; i < n; i++ {
			// read six numbers for the header
			for j := 0; j < 6; j++ {
				l.next()
			}
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
		}
		l.packet.subpacketData = l.value()
		l.ignore()
		l.packets = append(l.packets, l.packet)
		if len(l.input[l.pos:]) < 8 {
			return nil
		}
		return versionFn
	}
}

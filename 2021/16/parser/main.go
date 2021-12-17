package main

import (
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

const (
	version          = "version"
	typeID           = "typeID"
	literalValue     = "literalValue"
	fifteenBitNumber = "fifteenBitNumber"
	elevenBitNumber  = "elevenBitNumber"
	lengthTypeID     = "lengthTypeID"
	subpackets       = "subpackets"
)

/*
* 620080001611562C8802118E34 represents an operator packet (version 3) which
contains two sub-packets; each sub-packet is an operator packet that contains
two literal values. This packet has a version sum of 12.
*/

func main() {
	lines := internal.ReadInput()
	fmt.Println(hexToBin(lines[0]))
	tokens := make(chan item)
	l := &lex{
		input:  hexToBin(lines[0]),
		tokens: tokens,
	}
	//	versionsum := 0
	go l.run() // terminates after its done parsing
	for token := range tokens {

		fmt.Println(token)
	}
}

type packet struct {
	version int
	typeID  int

	// literal packet only
	literalValue int

	// operator packet only
	lengthTypeID     int
	fifteenBitnumber int
	elevenBitNumber  int
	size             int
	subpackets       []*packet
}

func (p *packet) String() string {
	if p.typeID == 4 {
		return fmt.Sprintf("LVP: %d", p.literalValue)
	}
	if p.typeID != 4 {
		return fmt.Sprintf("OP: {}")
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
	pt := ""
	for i := 0; i < 3; i++ {
		c := l.next()
		if c == "EOF" {
			return l.errorf("bad input, end came after a version")
		}
		pt += c
	}
	l.send(typeID)
	switch pt {
	case "100":
		return litValue
	default:
		return operatorValue
	}
}

func litValue(l *lex) stateFn {
	// is it the last bit in a literal value?
	c := l.next()
	if c == "EOF" {
		return l.errorf("invalid end found in the middle of a literal value")
	}
	if c == "0" {
		return lastLiteralValue
	}
	if c != "1" {
		return l.errorf("invalid character in literal value header")
	}
	return notLastLiteralValue
}

func lastLiteralValue(l *lex) stateFn {
	for i := 0; i < 4; i++ {
		c := l.next()
		if c == "EOF" {
			return l.errorf("bad end in what should be last literal value")
		}
	}
	l.send(literalValue)
	if l.peek() == "EOF" {
		return nil
	}
	// slurp the rest of the unnecessary bits
	for (l.pos-l.start)%4 != 0 {
		l.next()
	}
	// read one more char because this will be off by 1 above
	if l.peek() == "EOF" {
		return nil
	}
	l.ignore()
	if len(l.input)-l.pos < 8 {
		return nil
	}
	return versionFn
}

func notLastLiteralValue(l *lex) stateFn {
	for i := 0; i < 4; i++ {
		c := l.next()
		if c == "EOF" {
			return l.errorf("bad end in what should be last literal value")
		}
	}
	return litValue
}

func operatorValue(l *lex) stateFn {
	n := l.next()
	if n == "EOF" {
		return l.errorf("should be length-type, got end")
	}
	l.send(lengthTypeID)
	if n == "1" {
		return elevenBitParse
	}
	return fifteenBitParse
}

func fifteenBitParse(l *lex) stateFn {
	num := ""
	for i := 0; i < 15; i++ {
		c := l.next()
		if c == "EOF" {
			return l.errorf(fmt.Sprintf("found end when expecting a %d bit number", 15))
		}
		num += c
	}
	l.send(fifteenBitNumber)
	btz, _ := strconv.ParseInt(num, 2, 64)
	return parseBits(int(btz))
}

func elevenBitParse(l *lex) stateFn {
	num := ""
	for i := 0; i < 11; i++ {
		c := l.next()
		if c == "EOF" {
			return l.errorf(fmt.Sprintf("found end when expecting a %d bit number", 11))
		}
		num += c
	}
	l.send(elevenBitNumber)
	numSubPacket, _ := strconv.ParseInt(num, 2, 64)
	return parseSubPackets(int(numSubPacket))
}

func parseBits(n int) stateFn {
	return func(l *lex) stateFn {
		for i := 0; i < n; i++ {
			l.next()
		}
		l.send(subpackets)
		return nil
	}
}

func parseSubPackets(n int) stateFn {
	return func(l *lex) stateFn {
		for i := 0; i < n; i++ {
			// read six numbers for the header
			for j := 0; j < 6; j++ {
				l.next()
			}
			// read 5 at a time until you find a 5 that starts with 0
			for {
				stop := false
				header := l.next()
				if header == "0" {
					stop = true
				}
				for j := 0; j < 4; j++ {

				}
				if stop {
					break
				}
			}
		}
		l.send(subpackets)
		return nil
	}
}

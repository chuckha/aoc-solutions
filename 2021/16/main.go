package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
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
		input:      hexToBin(lines[0]),
		tokens:     tokens,
		subpackets: -1,
	}
	versionsum := 0
	go l.run() // terminates after its done parsing
	for tok := range tokens {
		fmt.Println(tok)
		if tok.kind == "version" {
			a, _ := strconv.ParseInt(tok.value, 2, 64)
			versionsum += int(a)
		}
		if tok.kind == "error" {
			break
		}
	}
	fmt.Println(versionsum)
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
}

type lex struct {
	input  string
	pos    int
	start  int
	tokens chan item

	subpackets int
}

func (l *lex) reset() {
	l.input = ""
	l.pos = 0
	l.start = 0
	l.tokens = nil
}
func (l *lex) run() {
	for state := version; state != nil; {
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

type packet struct {
	version    string
	packetType string
}

type stateFn func(*lex) stateFn

func version(l *lex) error {
	for i := 0; i < 3; i++ {
		if l.next() == "EOF" {
			return errors.New("found an early end in a version")
		}
	}
	l.send("version")
	pt := ""
	for i := 0; i < 3; i++ {
		c := l.next()
		if c == "EOF" {
			return errors.New("bad input, end came after a version")
		}
		pt += c
	}
	l.send("packet type")
	switch pt {
	case "100":
		for {
			c := l.next()
			if c == "EOF" {
				return errors.New("invalid end found in the middle of a literal value")
			}
			if c == "0" {
				for i := 0; i < 4; i++ {
					c := l.next()
					if c == "EOF" {
						return errors.New("bad end in what should be last literal value")
					}
				}
				l.send("literal value")
				continue
			}
			if c != "1" {
				return errors.New("invalid character in literal value header")
			}
			for i := 0; i < 4; i++ {
				c := l.next()
				if c == "EOF" {
					return errors.New("bad end in what should be last literal value")
				}
			}
		}
		l.send("literal value")
	default:
		n := l.next()
		if n == "EOF" {
			return errors.New("should be length-type, got end")
		}
		l.send("length type id")
		if n == "1" {
			num := ""
			for i := 0; i < 11; i++ {
				c := l.next()
				if c == "EOF" {
					return errors.New(fmt.Sprintf("found end when expecting a %d bit number", 11))
				}
				num += c
			}
			l.send(fmt.Sprintf("%d bit number", 11))
			//	fmt.Println(l.input)
			//	fmt.Println(l.input[:l.pos])
			// don't' actualy care what the number is??
			// need to make it parse x packets and then ...??????
			return nil
		}
		num := ""
		for i := 0; i < 15; i++ {
			c := l.next()
			if c == "EOF" {
				return errors.New(fmt.Sprintf("found end when expecting a %d bit number", 15))
			}
			num += c
		}
		l.send(fmt.Sprintf("%d bit number", 15))
		//	toParse, _ := strconv.ParseInt(num, 2, 64)
		return nil
	}
	return nil
}

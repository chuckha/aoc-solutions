package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	input := lines[0]
	fmt.Println(breakup(1, input))
}

func breakup(multiplier int, in string) int {
	if strings.Index(in, "(") == -1 {
		return len(in) * multiplier
	}
	unMultiplied := 0
	start := strings.Index(in, "(")
	if start > 0 {
		unMultiplied = len(in[:start])
	}
	end := strings.Index(in, ")")
	marker := in[start+1 : end]
	parts := strings.Split(marker, "x")
	charlen, _ := strconv.Atoi(parts[0])
	times, _ := strconv.Atoi(parts[1])
	data := in[end+1 : end+1+charlen]
	tail := in[end+1+charlen:]
	additions := 0
	if len(tail) > 0 {
		additions = breakup(multiplier, tail)
	}
	return additions + breakup(multiplier*times, data) + unMultiplied
}

/*
(27x12)(20x12)(13x14)(7x10)(1x12)A 20160 * 12 = 241920
(20x12)(13x14)(7x10)(1x12)A 1680 * 12 = 20160
(13x14)(7x10)(1x12)A = 120 * 14 = 1680
(7x10)(1x12)A = 12 * 10 = 120
(1x12)A = 12

12 * 12 * 14 * 10 * 12 * 1

(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN
(25x3)(3x3)ABC(2x3)XY(5x2)PQRST + X + (18x9)(3x2)TWO(5x7)SEVEN
(3x9)ABC(2x9)XY(5x6)PQRST + X + (3x18)TWO(5x63)SEVEN
27 + 18 + 30 + 1 + 54 + 315



multiplier
	data



(8x2)(3x2)TWO
(3x2)TWO(3x2)TWO
TWOTWOTWOTWO
*/

func buildRepeater(in string) {

}

type lex struct {
	items []string
	input string
	pos   int
}

func (l *lex) run() {
	for state := dataState; state != nil; {
		state = state(l)
	}
}
func (l *lex) next() string {
	if l.pos == len(l.input) {
		fmt.Println(len(strings.Join(l.items, "")))
		os.Exit(0)
	}
	out := string(l.input[l.pos])
	l.pos++
	return out
}
func (l *lex) add(item string) {
	if item == "" {
		return
	}
	l.items = append(l.items, item)
}

type state func(*lex) state

func dataState(l *lex) state {
	out := ""
	for {
		switch r := l.next(); {
		case r == "(":
			l.add(out)
			return markerState
		default:
			out += r
		}
	}
}

func markerState(l *lex) state {
	charCountS := ""
	repeatS := ""
	repeat := false
	for {
		if !repeat {
			switch r := l.next(); {
			case isNumeric(r):
				charCountS += r
			case r == "x":
				repeat = true
			}
		}
		if repeat {
			switch r := l.next(); {
			case isNumeric(r):
				repeatS += r
			case r == ")":
				charCount, _ := strconv.Atoi(charCountS)
				repeatTimes, _ := strconv.Atoi(repeatS)
				data := ""
				for i := 0; i < charCount; i++ {
					data += l.next()
				}
				//				fmt.Printf("adding %d copies of %q\n", repeatTimes, data)
				l.add(strings.Repeat(data, repeatTimes))
				return dataState
			}
		}
	}
}

func isNumeric(in string) bool {
	return strings.Contains("0123456789", in)
}

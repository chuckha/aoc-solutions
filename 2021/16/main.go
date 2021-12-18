package main

import (
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	input := hexToBin(lines[0])
	// read the header/type
	version := input[0:3]
	packetType := input[3:6]
	if packetType == "100" {
		endOfPacket, data := findEndOfLitearlPacket(input)
		fmt.Println(input[:endOfPacket], data)
	}
	// operator packet
	if packetType != "100" {

	}
	fmt.Println(version, packetType)
}

type packet struct {
	version int
	typeID  int

	literalValue int
}

func makeLiteralPacket(data string) *packet {
	version := binToDec(data[0:3])
	typeID := binToDec(data[3:6])
	_, literalValue := findEndOfLitearlPacket(data)
	return &packet{
		version:      version,
		typeID:       typeID,
		literalValue: literalValue,
	}
}

func binToDec(in string) int {
	out, _ := strconv.ParseInt(in, 2, 64)
	return int(out)
}

func findEndOfLitearlPacket(data string) (int, int) {
	actual := ""
	start := 0
	pos := 6 // ignore the header
	for {
		stop := false
		if string(data[pos]) == "0" {
			stop = true
		}
		pos++
		start = pos
		pos += 4 // read the actual data
		actual += data[start:pos]
		if stop {
			break
		}
	}
	outdata, _ := strconv.ParseInt(actual, 2, 64)
	return pos, int(outdata)
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

package main

import (
	"fmt"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	input := hexToBin(lines[0])
	pkt := packets(input)[0]
	fmt.Println(pkt)
	fmt.Println(pkt.value())
}

func packets(in string) []packet {
	//	fmt.Println("in", in)
	out := []packet{}
	for {
		if len(in) < 8 {
			return out
		}
		switch in[3:6] {
		case "100":
			endOfPacket, _ := findEndOfLitearlPacket(in)
			out = append(out, makeLiteralPacket(in[:endOfPacket]))
			//			fmt.Println("literal", in[:endOfPacket])
			in = in[endOfPacket:]
		default:
			switch in[6] {
			case '0':
				//				fmt.Println("parsing operator packet of length id 0")
				fifteenBitNumber := readFifteenBitNumber(in)
				//				fmt.Println("fifteenBitNumber", fifteenBitNumber)
				pkt := make15OperatorPacket(in[:22])
				//				fmt.Println("subpackets: ", in[22:22+fifteenBitNumber])
				pkt.subpackets = packets(in[22 : 22+fifteenBitNumber])
				out = append(out, pkt)
				in = in[22+fifteenBitNumber:] // 3 + 3 + 1= 7 + 15 = 22
			case '1':
				//				fmt.Println("parsing operator packet of length id 1")
				elevenBitNumber := readElevenBitNumber(in)
				//				fmt.Println("elevenBitnumber", elevenBitNumber)
				pkt := make11OperatorPacket(in)
				remainingPackets := packets(in[3+3+1+11:])
				// /				fmt.Println(remainingPackets)
				pkt.subpackets = remainingPackets[:elevenBitNumber]
				out = append(out, pkt)
				out = append(out, remainingPackets[elevenBitNumber:]...)
				return out
			default:
				panic("wat")
			}
		}
	}
}

type packet struct {
	version int
	typeID  int

	literalValue int

	lengthTypeID         int
	elevenBitNumber      int
	fifteenBitNumber     int
	literalSubpacketData string
	subpackets           []packet
}

func (p packet) versionSum() int {
	sum := p.version
	for _, sp := range p.subpackets {
		sum += sp.versionSum()
	}
	return sum
}

func (p packet) value() int {
	//	fmt.Println("type id", p.typeID)
	switch p.typeID {
	case 0:
		return p.subPacketSum()
	case 1:
		return p.subPacketProduct()
	case 2:
		return p.spMinimum()
	case 3:
		return p.spMaximum()
	case 4:
		return p.literalValue
	case 5:
		return p.gt()
	case 6:
		return p.lt()
	case 7:
		return p.eq()
	default:
		panic("oh shiiii")
	}
}

func (p packet) subPacketSum() int {
	sum := 0
	for _, sp := range p.subpackets {
		sum += sp.value()
	}
	return sum
}

func (p packet) subPacketProduct() int {
	prod := 1
	for _, sp := range p.subpackets {
		prod *= sp.value()
	}
	return prod
}

func (p packet) spMinimum() int {
	min := 99999999999
	for _, sp := range p.subpackets {
		val := sp.value()
		if val < min {
			min = val
		}
	}
	return min
}

func (p packet) gt() int {
	if p.subpackets[0].value() > p.subpackets[1].value() {
		return 1
	}
	return 0
}

func (p packet) lt() int {
	if p.subpackets[0].value() < p.subpackets[1].value() {
		return 1
	}
	return 0
}
func (p packet) eq() int {
	if p.subpackets[0].value() == p.subpackets[1].value() {
		return 1
	}
	return 0
}

func (p packet) spMaximum() int {
	max := 0
	for _, sp := range p.subpackets {
		val := sp.value()
		if val > max {
			max = val
		}
	}
	return max
}

func (p packet) String() string {
	if p.typeID == 4 {
		return fmt.Sprintf("LIT: %d", p.literalValue)
	}
	if p.typeID != 4 {
		return fmt.Sprintf("OP: %v", p.subpackets)
	}
	return "eh?"
}

func readFifteenBitNumber(in string) int {
	fifteenbitNumber := in[8 : 8+14] // 0 to 14 == 15
	return binToDec(fifteenbitNumber)
}

func readElevenBitNumber(in string) int {
	elevenBitNumber := in[8 : 8+10] // 0 to 10 == 11
	return binToDec(elevenBitNumber)
}

func make11OperatorPacket(data string) packet {
	version := binToDec(data[0:3])
	typeID := binToDec(data[3:6])
	lengthTypeID := binToDec(string(data[7]))
	elevenBitNumber := binToDec(data[8 : 8+10])
	return packet{
		version:         version,
		typeID:          typeID,
		lengthTypeID:    lengthTypeID,
		elevenBitNumber: elevenBitNumber,
	}
}

func make15OperatorPacket(data string) packet {
	version := binToDec(data[0:3])
	typeID := binToDec(data[3:6])
	lengthTypeID := binToDec(string(data[7]))
	fifteenBitNumber := binToDec(data[8:22])
	return packet{
		version: version,
		typeID:  typeID,
		// splice off the 15 bit number
		lengthTypeID:     lengthTypeID,
		fifteenBitNumber: fifteenBitNumber,
	}
}

func makeLiteralPacket(data string) packet {
	version := binToDec(data[0:3])
	typeID := binToDec(data[3:6])
	_, literalValue := findEndOfLitearlPacket(data)
	return packet{
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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	fmt.Println(part2(0, lines) * part2a(0, lines))

}

func part2(pos int, lines []string) int64 {
	if len(lines) == 1 {
		return lineToDec(lines[0])
	}
	mcb := findMostCommonBitInPosition(pos, lines)
	return part2(pos+1, filter(pos, mcb, lines))
}

func part2a(pos int, lines []string) int64 {
	if len(lines) == 1 {
		return lineToDec(lines[0])
	}
	lcb := findLeastCommonBitInPosition(pos, lines)
	return part2a(pos+1, filter(pos, lcb, lines))
}

func filter(pos int, bit string, lines []string) []string {
	out := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if string(line[pos]) == bit {
			out = append(out, line)
		}
	}
	return out
}
func findLeastCommonBitInPosition(pos int, lines []string) string {
	ones := 0
	zeroes := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[pos] == '1' {
			ones++
			continue
		}
		zeroes++
	}
	if zeroes > ones {
		return "1"
	}
	return "0"
}

func findMostCommonBitInPosition(pos int, lines []string) string {
	ones := 0
	zeroes := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[pos] == '1' {
			ones++
			continue
		}
		zeroes++
	}
	if zeroes > ones {
		return "0"
	}
	return "1"
}

func lineToDec(in string) int64 {
	dec, _ := strconv.ParseInt(in, 2, 64)
	return dec
}

// part 1
func iterate(lines []string) {
	counts := map[int]int{}
	totalLines := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		totalLines++
		for i, c := range line {
			if c == '1' {
				counts[i]++
			}
		}
	}
	gamma := make([]string, len(lines[0]))
	epsilon := make([]string, len(lines[0]))
	for pos, val := range counts {
		if val > totalLines/2 {
			gamma[pos] = "1"
			epsilon[pos] = "0"
			continue
		}
		gamma[pos] = "0"
		epsilon[pos] = "1"
	}
	// part 1 computation
	gam := strings.Join(gamma, "")
	eps := strings.Join(epsilon, "")

	gammaDec, _ := strconv.ParseInt(gam, 2, 64)
	epsilonDec, _ := strconv.ParseInt(eps, 2, 64)
	fmt.Println(gammaDec * epsilonDec)

}

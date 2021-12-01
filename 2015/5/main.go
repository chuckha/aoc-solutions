package main

import (
	"bufio"
	"fmt"
	"os"
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
	count := 0
	for _, line := range lines {
		if isReallyNice(line) {
			count++
		}
	}
	fmt.Println(count)
}

type pair struct {
	pos int
	val string
}

func isReallyNice(in string) bool {
	deets := map[string]int{}
	var prev rune
	var prevprev rune
	hasRepeatedLetterExactlyOneAway := false
	for i, c := range in {
		if i > 0 {
			val, ok := deets[string(prev)+string(c)]
			if ok {
				deets[string(prev)+string(c)] = val - i
			} else {
				deets[string(prev)+string(c)] = i
			}
		}
		if prevprev == c {
			hasRepeatedLetterExactlyOneAway = true
		}
		prevprev = prev
		prev = c
	}
	hasTwoLettersThatRepeatButDoNotOverlap := false
	for _, v := range deets {
		if v < -1 {
			hasTwoLettersThatRepeatButDoNotOverlap = true
		}
	}

	return hasTwoLettersThatRepeatButDoNotOverlap && hasRepeatedLetterExactlyOneAway
}

// part 1

func isNice(in string) bool {
	vowelCount := 0
	hasDouble := false
	var prev rune
	for _, c := range in {
		if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
			vowelCount++
		}
		if c == prev {
			hasDouble = true
		}
		if prev == 'a' && c == 'b' {
			return false
		}
		if prev == 'c' && c == 'd' {
			return false
		}
		if prev == 'p' && c == 'q' {
			return false
		}
		if prev == 'x' && c == 'y' {
			return false
		}
		prev = c
	}
	//	fmt.Fprintln(os.Stderr, "vowel count", vowelCount)
	//	fmt.Fprintln(os.Stderr, "hasDouble", hasDouble)
	return vowelCount >= 3 && hasDouble
}

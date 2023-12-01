package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()
	count := 0
	count2 := 0
	for _, line := range input {
		sec := newSection(line)
		if completeOverlap(sec) {
			count++
		}
		//		fmt.Println(sec, noOverlap(sec))
		if !noOverlap(sec) {
			count2++
		}
	}
	fmt.Println(count, count2)
}

func noOverlap(pair []section) bool {
	if pair[0].high < pair[1].low {
		return true
	}
	if pair[1].high < pair[0].low {
		return true
	}
	return false
}

func completeOverlap(pair []section) bool {
	if pair[0].low <= pair[1].low && pair[0].high >= pair[1].high {
		return true
	}
	if pair[1].low <= pair[0].low && pair[1].high >= pair[0].high {
		return true
	}
	return false
}

type section struct {
	low  int
	high int
}

func newSection(item string) []section {
	out := make([]section, 2)
	for i, part := range strings.Split(item, ",") {
		nums := strings.Split(part, "-")
		l, _ := strconv.Atoi(nums[0])
		r, _ := strconv.Atoi(nums[1])
		out[i].low = l
		out[i].high = r
	}
	return out
}

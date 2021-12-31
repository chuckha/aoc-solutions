package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()

	spans := []span{}
	for _, line := range input {
		spans = append(spans, newSpan(line))
	}
	sort.Sort(spanSlice(spans))
	for _, span := range spans {
		fmt.Println(span)
	}
	fmt.Println("---------------------------")
	for i := 0; i < len(spans)-1; i++ {
		//		fmt.Println("combining", i, spans[i], spans[i+1])
		news, ok := spans[i].combine(spans[i+1])
		if ok {
			fmt.Println(spans[i], spans[i+1], news)
			tail := append([]span{news}, spans[i+2:]...)
			head := spans[:i]
			spans = append(head, tail...)
			//			fmt.Println(spans)
			i = -1
		}
	}
	validIPs := []int{}
	count := 0
	for i := 0; i < len(spans)-1; i++ {
		fmt.Printf("%d - %d = %d\n", spans[i+1].low, spans[i].high, spans[i+1].low-spans[i].high-1)
		validIPs = append(validIPs, spans[i+1].low-1)
		count += spans[i+1].low - spans[i].high - 1
	}
	fmt.Println(count)
	fmt.Println(spans)
	for _, ip := range validIPs {
		fmt.Println(ip)
	}
	fmt.Println(len(validIPs))
}

type span struct {
	low  int
	high int
}

func (s span) combine(s2 span) (span, bool) {
	// they do connect bcause inclusive
	if s2.low-s.high == 1 {
		return span{low: s.low, high: s2.high}, true
	}
	if s.low <= s2.low && s.high <= s2.high && s.high >= s2.low {
		return span{
			low:  s.low,
			high: s2.high,
		}, true
	}
	if s.low <= s2.low && s.high >= s2.high {
		return s, true
	}
	return span{}, false
}

type spanSlice []span

func (s spanSlice) Len() int {
	return len(s)
}
func (s spanSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s spanSlice) Less(i, j int) bool {
	return s[i].low < s[j].low
}

func newSpan(line string) span {
	w := strings.Split(line, "-")
	low, _ := strconv.Atoi(w[0])
	high, _ := strconv.Atoi(w[1])
	return span{
		low:  low,
		high: high,
	}
}

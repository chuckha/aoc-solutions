package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

// 650 * 10000
// 6500000 signal length
// 5977737 is the offset

var replications = 10000
var inputSize = 0
var applySignalCount = 100

func main() {
	in := input.GetInput(2019, 16)[0]
	input := stringToIntSlice(in)
	input = repeat(input, replications)
	off := getOffset(input)
	out := input[off:]
	for i := 0; i < applySignalCount; i++ {
		tt := time.Now()
		wholeSum := sum(out)
		out = apply6(out, wholeSum)
		oneItr := time.Now().Sub(tt)
		fmt.Println("one finished in", oneItr, "expected time remaining", time.Duration(99-i)*oneItr)
	}
	fmt.Println("OFFSET", off)
	fmt.Println(out[:8])

}

func apply6(nums []int, wholeSum int) []int {
	out := make([]int, len(nums))
	out[0] = wholeSum % 10
	for i := 1; i < len(nums); i++ {
		wholeSum = (wholeSum - nums[i-1])
		out[i] = wholeSum % 10
	}
	return out
}

// func addit(signal, indexes []int) int {
// 	sum := 0
// 	for _, i := range indexes {
// 		sum += signal[i]
// 	}
// 	return sum
// }

// func subit(signal, indexes []int) int {
// 	sum := 0
// 	for _, i := range indexes {
// 		sum -= signal[i]
// 	}
// 	return sum
// }

// func apply5(signal string, repititions int) []int {
// 	out := make([]int, len(signal))
// 	nums := stringToIntSlice(signal)
// 	// make a pattern and apply it to a signal to get a single number
// 	lookup := map[key]int{}
// 	cacheHits := 0
// 	for i := 0; i < len(out); i++ {
// 		fmt.Println("pattern key:", i)
// 		patterns := patternsAsKey(len(signal), i+1)
// 		add, sub := patternsAsKey2(len(signal), i+1)
// 		num := addit(nums, add) + subit(nums, sub)

// 		// pl := make([][]int, len(patterns))
// 		// for j := range patterns {
// 		// 	pl[j] = stringToIntSlice(patterns[j])
// 		// }
// 		//		fmt.Println(patterns)
// 		r := repititions
// 		c := 0
// 		for _, pattern := range patterns {
// 			k := key{
// 				signal:  signal,
// 				pattern: pattern,
// 			}
// 			if _, ok := lookup[k]; !ok {
// 				lookup[k] = applyOnce5(nums, stringToIntSlice(pattern))
// 				for _, c := range nums {
// 					fmt.Printf("%02d", c)
// 				}
// 				fmt.Println()
// 				for _, c := range stringToIntSlice(pattern) {
// 					fmt.Printf("%02d", c)
// 				}
// 				fmt.Println()
// 			} else {
// 				cacheHits++
// 			}
// 			// if i have more patterns than repitions, just use the first n pattners
// 			// if i have more repititions than patterns,
// 			out[i] = ((out[i] + lookup[k]) * r) % 10
// 			c++
// 			if c == r {
// 				break
// 			}
// 		}
// 	}
// 	//	fmt.Println("cache hits", cacheHits)
// 	return out
// }

type key struct {
	signal  string
	pattern string
}

func apply4(signal string, repititions int) []int {
	out := make([]int, len(signal))
	nums := stringToIntSlice(signal)
	// make a pattern and apply it to a signal to get a single number
	lookup := map[key]int{}
	cacheHits := 0
	for i := 0; i < len(out); i++ {
		patterns := patternsAsKey(len(signal), i+1)
		pl := make([][]int, len(patterns))
		for j := range patterns {
			pl[j] = stringToIntSlice(patterns[j])
		}
		//		fmt.Println(patterns)
		r := repititions
		for r > 0 {
			for q, pattern := range patterns {
				k := key{
					signal:  signal,
					pattern: pattern,
				}
				if _, ok := lookup[k]; !ok {
					lookup[k] = applyOnce5(nums, pl[q])
				} else {
					cacheHits++
				}
				// if i have more patterns than repitions, just use the first n pattners
				// if i have more repititions than patterns,
				out[i] = (out[i] + lookup[k]) % 10
				r--
				if r == 0 {
					break
				}
			}
		}
	}
	//	fmt.Println("cache hits", cacheHits)
	return out
}

func applyOnce5(signal, pattern []int) int {
	out := 0
	for i := 0; i < len(signal); i++ {
		if pattern[i] == 0 {
		}
		out += signal[i] * pattern[i]
	}
	return abs(out) % 10
}

func intSliceToString(p []int) string {
	var out strings.Builder
	for _, d := range p {
		out.WriteString(fmt.Sprintf("%d", d))
	}
	return out.String()
}

func stringToIntSlice(in string) []int {
	out := []int{}
	idx := 0
	for i := 0; i < len(in); i++ {
		if string(in[i]) != "-" {
			x, _ := strconv.Atoi(string(in[i]))
			out = append(out, x)
			idx++
			continue
		}
		i++
		x, _ := strconv.Atoi(string(in[i]))
		out = append(out, -x)
		idx++
	}
	return out
}

func pattern(size, n int) [][]int {
	basePattern := []int{0, 1, 0, -1}
	first := []int{}
	out := [][]int{}
	bpIdx := 0
	init := true
	toGrab := n
	for {
		// make the next pattern
		next := make([]int, size)
		// start the index at 0
		idx := 0
		// while the index is still within bounds (size of the current pattern)
		for idx < len(next) {
			if init {
				toGrab -= 1 // but the first time, get one less
			}
			// grab the same item toGrab times
			for toGrab > 0 {
				next[idx] = basePattern[bpIdx]
				idx++
				toGrab--
				if idx >= size {
					break
				}
			}
			if toGrab == 0 {
				bpIdx = (bpIdx + 1) % 4
				toGrab = n
			}
			if init {
				init = false
				toGrab = n
			}
		}
		if internal.EqualSlice(next, first) {
			return out
		}
		out = append(out, next)
		if len(out) == 1 {
			first = next
		}
	}
}

// return a list of indexes to add and a list of indexes to subtract
func patternsAsKey2(size, n int) ([]int, []int) {
	basePattern := []int{0, 1, 0, -1}
	first := ""
	add := []int{}
	sub := []int{}
	bpIdx := 0
	init := true
	toGrab := n
	signalIdx := 0

	for {
		// make the next pattern
		next := make([]int, size)
		// start the index at 0
		idx := 0
		// while the index is still within bounds (size of the current pattern)
		for idx < size {
			//			fmt.Println("loop2", toGrab, idx, len(next), size, init, n)
			if init {
				toGrab -= 1 // but the first time, get one less
			}
			// grab the same item toGrab times
			for toGrab > 0 {
				bp := basePattern[bpIdx]
				if bp == 1 {
					add = append(add, signalIdx%size)
				}
				if bp == -1 {
					sub = append(sub, signalIdx%size)
				}
				signalIdx++
				next[idx] = basePattern[bpIdx]
				idx++
				toGrab--
				if idx >= size {
					break
				}
			}
			if toGrab == 0 {
				bpIdx = (bpIdx + 1) % 4
				toGrab = n
			}
			if init {
				init = false
				toGrab = n
			}
		}
		nxt := intSliceToString(next)
		if nxt == first {
			return add, sub
		}
		if first == "" {
			first = nxt
		}
	}
}

func patternsAsKey(size, n int) []string {
	basePattern := []int{0, 1, 0, -1}
	first := ""
	out := []string{}
	bpIdx := 0
	init := true
	toGrab := n
	for {
		// make the next pattern
		next := make([]int, size)
		// start the index at 0
		idx := 0
		// while the index is still within bounds (size of the current pattern)
		for idx < len(next) {
			//			fmt.Println("loop2", toGrab, idx, len(next), size, init, n)
			if init {
				toGrab -= 1 // but the first time, get one less
			}
			// grab the same item toGrab times
			for toGrab > 0 {
				next[idx] = basePattern[bpIdx]
				idx++
				toGrab--
				if idx >= size {
					break
				}
			}
			if toGrab == 0 {
				bpIdx = (bpIdx + 1) % 4
				toGrab = n
			}
			if init {
				init = false
				toGrab = n
			}
		}
		nxt := intSliceToString(next)
		if nxt == first {
			return out
		}
		out = append(out, nxt)
		if len(out) == 1 {
			first = nxt
		}
	}
}

// signal is just og
// n will generate the pattern for us until it repeats
// then we apply the pattern to OG 10000 times
func applyOnce4(signal []int, n int) int {
	out := 0
	size := len(signal)
	twoN := n + n
	i := 0
	end := i + n
	for {
		out += sum(signal[i:min(end, size)])
		i += twoN
		end = i + n
		if i >= size {
			break
		}
		out -= sum(signal[i:min(end, size)])
		i += twoN
		end = i + n
		if i >= size {
			break
		}
	}

	return abs(out) % 10
}

func apply3(og, signal []int) []int {
	// make a pattern and apply it to a signal to get a single number
	out := make([]int, len(signal))
	// in the case of 32 % i == 0
	for i := 1; i <= len(out)/2; i++ {
		patternRepeat := len(og) % i
		tt := time.Now()
		out[i-1] = (applyOnce3(og, i) * replications / patternRepeat) % 10
		fmt.Println("replication shortcut takes", time.Now().Sub(tt))
	}
	for i := (len(out) / 2) + 1; i <= len(out); i++ {
		tt := time.Now()
		out[i-1] = sum(signal[i-1:]) % 10
		fmt.Println("last half takes", time.Now().Sub(tt))
	}
	return out
}
func applyOnce3(signal []int, n int) int {
	out := 0
	size := len(signal)
	twoN := n + n
	i := 0
	end := i + n
	for {
		out += sum(signal[i:min(end, size)])
		i += twoN
		end = i + n
		if i >= size {
			break
		}
		out -= sum(signal[i:min(end, size)])
		i += twoN
		end = i + n
		if i >= size {
			break
		}
	}

	return abs(out)
}

func old(nums []int) {
	inputSize = len(nums)
	offset := getOffset(nums)
	nums = repeat(nums, replications)
	// we can be smarter about the repeat part
	for i := 0; i < applySignalCount; i++ {
		tt := time.Now()
		nums = apply2(nums)
		oneItr := time.Now().Sub(tt)
		fmt.Println("one finished in", oneItr, "expected time remaining", time.Duration(99-i)*oneItr)
	}
	if len(nums) > offset {
		fmt.Println(nums[offset:8])
	} else {
		fmt.Println(nums[:8])
	}
}

func getOffset(in []int) int {
	offset := make([]int, 7)
	copy(offset, in)
	out := ""
	for _, v := range offset {
		out += fmt.Sprintf("%d", v)
	}
	o, _ := strconv.Atoi(out)
	return o
}

// 12345678 repeated 3 times should be [2 4 4 2 6 0 4 8 7 1 8 8 2 7 1 4 6 5 3 0 6 1 5 8]

func repeat(in []int, n int) []int {
	out := make([]int, len(in)*n)
	for j := 0; j < n; j++ {
		for i := 0; i < len(in); i++ {
			out[i+j*len(in)] = in[i]
		}
	}
	return out
}

func apply2(signal []int) []int {
	// make a pattern and apply it to a signal to get a single number
	out := make([]int, len(signal))
	// in the case of 32 % i == 0
	for i := 1; i <= len(out)/2; i++ {
		if i%inputSize == 0 {
			out[i-1] = applyRepeat(signal[i-1:], i)
			continue
		}
		out[i-1] = applyOnce2(signal[i-1:], i)
	}
	for i := (len(out) / 2) + 1; i <= len(out); i++ {
		out[i-1] = sum(signal[i-1:]) % 10
	}
	return out
}

func applyRepeat(signal []int, n int) int {
	firstPattern := len(signal) / replications
	return applyOnce2(signal[:firstPattern], n)
}

func applyOnce2(signal []int, n int) int {
	out := 0
	size := len(signal)
	twoN := n + n
	i := 0
	end := i + n
	for {
		out += sum(signal[i:min(end, size)])
		i += twoN
		end = i + n
		if i >= size {
			break
		}
		out -= sum(signal[i:min(end, size)])
		i += twoN
		end = i + n
		if i >= size {
			break
		}
	}

	return abs(out) % 10
}

func sum(in []int) int {
	out := 0
	for i := 0; i < len(in); i++ {
		out = out + in[i]
	}
	return out
}
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

/* if i % len(in) ==0, use repeat
12345678 12345678 12345678 12345678 12345678 12345678 0
10101010 10101010 10101010 10101010 10101010 10101010 1 %32 == 1
01100110 01100110 01100110 01100110 01100110 01100110 2 %32 == 0
00111000 11100011 10001110 00111000 11100011 10001110 3 %32 == 3
00011110 00011110 00011110 00011110 00011110 00011110 4 %32 ==
00001111 10000011 11100000 11111000 00111110 00001111 5
00000111 11100000 01111110 00000011 11110000 00111111 6
00000011 11111000 00001111 11100000 00111111 10000000 7
00000001 11111110 00000001 11111110 00000001 11111110 8
00000000 11111111 10000000 00111111 11100000 00001111 9
00000000 01111111 11100000 00000111 11111110 00000000 10
00000000 00111111 11111000 00000000 11111111 11100000 11
00000000 00011111 11111110 00000000 00011111 11111110 12
00000000 00001111 11111111 10000000 00000011 11111111 13
00000000 00000111 11111111 11100000 00000000 01111111 14
00000000 00000011 11111111 11111000 00000000 00001111 15
00000000 00000001 11111111 11111110 00000000 00000001 16
00000000 00000000 11111111 11111111 10000000 00000000 17
00000000 00000000 01111111 11111111 11100000 00000000 18
00000000 00000000 00111111 11111111 11111000 00000000 19
00000000 00000000 00011111 11111111 11111110 00000000 20
00000000 00000000 00001111 11111111 11111111 10000000 21
00000000 00000000 00000111 11111111 11111111 11100000 22
00000000 00000000 00000011 11111111 11111111 11111000 23
00000000 00000000 00000001 11111111 11111111 11111110 24
00000000 00000000 00000000 11111111 11111111 11111111 25
00000000 00000000 00000000 01111111 11111111 11111111 26
00000000 00000000 00000000 00111111 11111111 11111111 27

*/

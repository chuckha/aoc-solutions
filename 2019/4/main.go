package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()[0]
	eps := strings.Split(input, "-")
	//	fmt.Println(eps)
	//	fmt.Println(increasing([]byte("10")))
	// fmt.Println(doubleGroups([]byte("112224")))
	// os.Exit(0)
	cur := []byte(eps[0])
	fin := []byte(eps[1])
	c := 0
	for !eq(cur, fin) {
		if increasing(cur) && doubleGroups(cur) && double(cur) {
			fmt.Println(string(cur))
			c++
		}
		cur = inc(cur)
	}
	fmt.Println(c)
}

func eq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if b[i] != c {
			return false
		}
	}
	return true
}

// 0 == 48
// 9 == 57
func inc(in []byte) []byte {
	if len(in) == 1 {
		if in[0] == 57 {
			return []byte{49, 48}
		}
		return []byte{in[0] + 1}
	}
	if in[len(in)-1] >= 57 {
		return append(inc(in[:len(in)-1]), 48)
	}
	return append(in[:len(in)-1], in[len(in)-1]+1)
}

func increasing(in []byte) bool {
	for i := 0; i < len(in)-1; i++ {
		if in[i] > in[i+1] {
			return false
		}
	}
	return true
}

// 1528 too low (part 1?)
// part 2:
// 1504 is not the right answer
// 1737 is not the right answer
// 2129 too high
// 2514 too high

func doubleGroups(in []byte) bool {
	///	fmt.Println(string(in))
	doubleFound := false
	// count how many in a row are the same
	for i := 0; i < len(in); i++ {
		j := i + 1
		item := in[i]
		count := 1
		for j < len(in) && in[j] == item {
			//			fmt.Println(i, j, count, string(item))
			count++
			j++
		}
		if count == 2 {
			doubleFound = true
		}
		i = j - 1 // the loop adds one more
	}
	return doubleFound
	// if the first number is the same as the second number
	// if the third number is the same as the first and second number
	// if the third number is different than the fourth number, fail
	// if the first number and the second number are the same, read the next number
	//	if it' sthe same as the first two, read the fourth number
	// if it' sthe same read the fifth number
	// if it's the same read the sixth number
	// if it's the same we're done success
}

func double(in []byte) bool {
	for i := 0; i < len(in)-1; i++ {
		for j := i + 1; j < len(in); j++ {
			if in[i] == in[j] {
				return true
			}
		}
	}
	return false
}

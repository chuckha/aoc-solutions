package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

	length := 0
	for _, line := range lines {
		box := newBox(line)
		length += ribbon(box)
	}
	fmt.Println(length)
	os.Exit(0)
	// part 1
	total := 0
	for _, line := range lines {
		box := newBox(line)
		total += area(box)
	}
	fmt.Println(total)
	// 2 * l * w + 2 * w * h + 2 * h * l = sa
	// 2 ( l * w + w * h + h * l ) = sa
	// 2 ( w ( l + h) + h * l)

}

func newBox(in string) []int {
	dims := strings.Split(in, "x")
	out := make([]int, 3)

	for i, d := range dims {
		j, _ := strconv.Atoi(d)
		out[i] = j
	}
	sort.Ints(out)
	return out
}

func ribbon(in []int) int {
	return in[0]*2 + in[1]*2 + in[0]*in[1]*in[2]
}

func area(in []int) int {
	return 2*(in[0]*in[1]+in[0]*in[2]+in[1]*in[2]) + in[0]*in[1]
}

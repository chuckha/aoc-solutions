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
	nums := make([]int, len(lines))
	for i, line := range lines {
		d, _ := strconv.Atoi(strings.TrimSpace(line))
		nums[i] = d
	}

	sums := []int{}
	for i := 2; i < len(nums); i++ {
		sums = append(sums, nums[i]+nums[i-1]+nums[i-2])
	}

	// part 1
	o := 0
	for i := 1; i < len(sums); i++ {
		if sums[i]-sums[i-1] > 0 {
			o++
		}
	}
	fmt.Println(o)
}

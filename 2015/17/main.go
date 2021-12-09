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
		data := scanner.Text()
		data = strings.TrimSpace(data)
		if data == "" {
			continue
		}
		lines = append(lines, data)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	numbers := []int{}
	for _, line := range lines {
		n, _ := strconv.Atoi(line)
		numbers = append(numbers, n)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbers)))
	// [47 46 36 36 32 32 31 30 28 26 19 15 11 11 5 3 3 3 1 1]
	// [47, 46, 36, 32, 31, 30, 28, 26, 19, 15, 11, 5, 3, 1]
	// [ 1,  1,  2,  2,  1,  1,  1,  1,  1,  1,  2, 1, 3, 2]
	_, min := minboxes(150, []int{}, numbers, noop)
	answersWithFilter, _ := minboxes(150, []int{}, numbers, lengthFilter(min))
	fmt.Println(answersWithFilter)
	//fmt.Println(memoizedCount(6, []int{36, 32, 32, 31, 30, 28, 26, 19, 11, 11, 5, 3, 3, 3, 1, 1}))
	//	10 6 [36 32 32 31 30 28 26 19 11 11 5 3 3 3 1 1]
	//	fmt.Println(ways(25, []int{}, []int{20, 15, 10, 5, 5}))
}

func noop([]int) bool {
	return true
}
func lengthFilter(size int) func([]int) bool {
	return func(in []int) bool {
		return len(in) == size
	}
}

func minboxes(s int, f, c []int, filter func([]int) bool) (int, int) {
	min := 1000
	var ways func(int, []int, []int) int
	ways = func(sum int, fixed, coins []int) int {
		if sum == 0 {
			if len(fixed) < min {
				min = len(fixed)
			}
			if filter(fixed) {
				return 1
			}
			return 0
		}
		if sum < 0 {
			return 0
		}
		count := 0
		for i, coin := range coins {
			removed := make([]int, len(coins))
			copy(removed, coins)
			removed = removed[i+1:]
			added := make([]int, len(fixed))
			copy(added, fixed)
			added = append(added, coin)
			count += ways(sum-coin, added, removed)
		}
		return count
	}
	return ways(s, f, c), min
}

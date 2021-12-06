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

	timers := strings.Split(lines[0], ",")
	sum := 0
	for _, timer := range timers {
		t, _ := strconv.Atoi(timer)
		sum += genMemo(t, 256)
	}
	fmt.Println(sum)
	// one fish over 80 days will generate x new fish
	// number, days remaining (3, 18) -> [(8, 14), (8, 7), (8, 0)]
	//                        (4, 18) -> [(8, 13), (8, 6)]
	//                        (3, 18) -> same as above
	//                        (1, 18) -> [(8, 16), (8, 9), (8, 2)]
	//                        (8, 16) -> [(8, 6)]
	//						  (8, 7) ->  [()]

}

var daysRequiredToSpawn = 7

type key struct{ a, b int }

func genMemo(t, r int) int {
	memo := map[key]int{}
	var generate func(int, int) int
	generate = func(timer, remaining int) int {
		if out, ok := memo[key{timer, remaining}]; ok {
			return out
		}
		if remaining <= timer {
			return 1
		}

		sum := 1
		if timer < remaining {
			remaining = remaining - timer - 1
			for remaining >= 0 {
				sumA := generate(8, remaining)
				memo[key{timer, remaining}] = sumA
				sum += sumA
				remaining -= daysRequiredToSpawn
			}
		}
		memo[key{timer, remaining}] = sum
		return sum
	}
	return generate(t, r)
}

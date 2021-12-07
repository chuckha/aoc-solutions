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

	ns := strings.Split(lines[0], ",")
	nums := make([]int, len(ns))
	for i, n := range ns {
		nums[i], _ = strconv.Atoi(n)
	}

	lowest := -1
	highest := 0
	for _, num := range nums {
		if num > highest {
			highest = num
		}
		if num < lowest {
			lowest = num
		}
	}

	leastFuel := 9999999999
	for i := lowest; i < highest; i++ {
		fuelBurnt := totalDistance(nums, i)
		if fuelBurnt < leastFuel {
			leastFuel = fuelBurnt
		}
	}
	fmt.Println(leastFuel)
}

func totalDistance(crabs []int, alignment int) int {
	totalFuel := 0
	for _, crab := range crabs {
		if crab > alignment {
			n := crab - alignment
			totalFuel += ((n * (n + 1)) / 2)
		}
		if crab < alignment {
			n := alignment - crab
			totalFuel += ((n * (n + 1)) / 2)
		}
	}
	return totalFuel
}

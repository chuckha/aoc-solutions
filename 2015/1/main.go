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

	floor := 0
	count := 0
	for _, c := range lines[0] {
		count++
		if c == '(' {
			floor++
		}
		if c == ')' {
			floor--
		}
		if floor < 0 {
			fmt.Println(count)
			os.Exit(0)
		}

	}
	fmt.Println(floor)
}

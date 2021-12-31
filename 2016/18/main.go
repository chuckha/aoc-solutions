package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	firstRow := internal.ReadInput()[0]
	size := 40
	size = 400000 // part 2
	rows := []string{firstRow}
	for len(rows) < size {
		rows = append(rows, nextRow(rows[len(rows)-1]))
	}
	fmt.Println(count(rows))
}
func count(in []string) int {
	count := 0
	for _, row := range in {
		for _, c := range row {
			if c == '.' {
				count++
			}
		}
	}
	return count
}

func nextRow(prev string) string {
	var out strings.Builder
	for i := 0; i < len(prev); i++ {
		// default everything safe
		left := "."
		center := "."
		right := "."
		if i > 0 {
			left = string(prev[i-1])
		}
		center = string(prev[i])
		if i < len(prev)-1 {
			right = string(prev[i+1])
		}
		above := left + center + right
		switch above {
		case "...":
			out.WriteString(".")
		case "^..":
			out.WriteString("^")
		case ".^.":
			out.WriteString(".")
		case "..^":
			out.WriteString("^")
		case "^^.":
			out.WriteString("^")
		case "^.^":
			out.WriteString(".")
		case ".^^":
			out.WriteString("^")
		case "^^^":
			out.WriteString(".")

		}
	}
	return out.String()
}

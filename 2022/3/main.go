package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()
	sum := 0
	for i := 0; i < len(input); i += 3 {
		sum += val(findSimilar(input[i], input[i+1], input[i+2]))
	}
	// Part 1
	// for _, in := range input {
	// 	left, right := split(in)
	// 	_, _, overlap := compartments(left, right)
	// 	for k := range overlap {
	// 		sum += val(k)
	// 	}
	// }
	fmt.Println(sum)
}

func findSimilar(a, b, c string) string {
	countsa := map[string]int{}
	for _, v := range a {
		countsa[string(v)]++
	}
	countsb := map[string]int{}
	for _, v := range b {
		countsb[string(v)]++
	}
	countsc := map[string]int{}
	for _, v := range c {
		countsc[string(v)]++
	}
	setCount := map[string]int{}
	for k := range countsa {
		setCount[k]++
	}
	for k := range countsb {
		setCount[k]++
	}
	for k := range countsc {
		setCount[k]++
	}
	fmt.Println(setCount)
	for s, c := range setCount {
		if c == 3 {
			return s
		}
	}
	panic("no badge")
}

func split(line string) (string, string) {
	return line[:len(line)/2], line[len(line)/2:]
}

func compartments(left, right string) (map[string]struct{}, map[string]struct{}, map[string]struct{}) {
	overlap := map[string]struct{}{}
	lc := map[string]struct{}{}
	for _, t := range left {
		lc[string(t)] = struct{}{}
	}
	rc := map[string]struct{}{}
	for _, t := range right {
		rc[string(t)] = struct{}{}
		if _, ok := lc[string(t)]; ok {
			overlap[string(t)] = struct{}{}
		}
	}
	return lc, rc, overlap
}

func val(l string) int {
	alpha := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return strings.Index(alpha, l) + 1
}

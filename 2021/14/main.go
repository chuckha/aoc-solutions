package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	rules := map[string]string{}
	for _, line := range lines[1:] {
		parts := strings.Split(line, " -> ")
		rules[parts[0]] = parts[1]
	}
	start := lines[0]
	pairCounts := countPairs(rules, start)
	letterCounts := countLetters(start)
	for i := 0; i < 40; i++ {
		///		fmt.Println(countLetters(start), len(start))
		//		fmt.Println(countPairs(rules, start))
		pairCounts, letterCounts = nextRound(pairCounts, letterCounts, rules)
		///		start = ruleApplications(rules, pairSplit(start))
	}
	//	fmt.Println(letterCounts)
	fmt.Println(mostAndLeastv2(letterCounts))
}

func nextRound(inputPairs, inputLetters map[string]int, rules map[string]string) (map[string]int, map[string]int) {
	pairCount := map[string]int{}
	letterCount := map[string]int{}
	for k, v := range inputLetters {
		letterCount[k] = v
	}
	for k, v := range inputPairs {
		gen := rules[k]
		letterCount[gen] += v
		left := string(k[0]) + gen
		right := gen + string(k[1])
		pairCount[left] += v
		pairCount[right] += v
	}
	return pairCount, letterCount
}

func countPairs(rules map[string]string, in string) map[string]int {
	pairs := pairSplit(in)
	out := map[string]int{}
	// for k := range rules {
	// 	out[k] = 0
	// }
	for _, pair := range pairs {
		out[pair]++
	}
	return out
}

func countLetters(in string) map[string]int {
	counts := map[string]int{}
	for _, c := range in {
		counts[string(c)]++
	}
	return counts
}

func mostAndLeastv2(letterCounts map[string]int) int {
	max := 0
	min := 999999999999999999
	for _, v := range letterCounts {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	//	fmt.Println("max", max, "min", min)
	return max - min
}

func mostAndLeast(in string) int {
	counts := map[string]int{}
	for _, c := range in {
		counts[string(c)]++
	}
	max := 0
	min := 9999
	for _, v := range counts {
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max - min
}

func pairSplit(in string) []string {
	out := make([]string, len(in)-1)
	for i := 0; i < len(in)-1; i++ {
		out[i] = in[i : i+2]
	}
	return out
}

func ruleApplications(rules map[string]string, pairs []string) string {
	last := string(pairs[len(pairs)-1][1])
	for i, pair := range pairs {
		if insert, ok := rules[pair]; ok {
			pairs[i] = string(pair[0]) + insert
		}
	}
	pairs[len(pairs)-1] += last
	return strings.Join(pairs, "")
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*

Sprinkles: capacity 2, durability 0, flavor -2, texture 0, calories 3
Butterscotch: capacity 0, durability 5, flavor -3, texture 0, calories 3
Chocolate: capacity 0, durability 0, flavor 5, texture -1, calories 8
Candy: capacity 0, durability -1, flavor 0, texture 5, calories 8

score = x1 * x2 * x3 * x4
x1 = 2S            (S)
x2 = 5B - K        (B, K)
x3 = -2S - 3B + 5C (S, B, C) (B, C, X1)
x4 = -C + 5K       (C, K) -> (C, B, X2)
500 = 3S + SB + 8C + 8K

descition variables:
	S, B, C, and K
objective function:
	score = 2S * (5B - K) * (-2S - 3B + 5C) * (-C + 5K)
constraints:
	S + B + C + K = 100
	S, B, C K > 0
	x1, x2, x3, x4 > 0
*/

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
	ingredients := make([]*ingredient, len(lines))
	for i, line := range lines {
		words := strings.Split(line, " ")
		name := strings.TrimSuffix(words[0], ":")
		capacity, _ := strconv.Atoi(strings.TrimSuffix(words[2], ","))
		durability, _ := strconv.Atoi(strings.TrimSuffix(words[4], ","))
		flavor, _ := strconv.Atoi(strings.TrimSuffix(words[6], ","))
		texture, _ := strconv.Atoi(strings.TrimSuffix(words[8], ","))
		calories, _ := strconv.Atoi(words[10])
		ingredients[i] = newIng(name, capacity, durability, flavor, texture, calories)
	}
	//	fmt.Println(exampleScore([]int{44, 56}))
	//	fmt.Println(score(ingredients))
	highest := 0
	for _, combo := range division(100) {
		score := fastScore(combo)
		if score > highest {
			highest = score
		}
	}
	fmt.Println(highest)
}

func exampleScore(in []int) int {
	return (-in[0] + in[1]*2) * (in[0]*-2 + in[1]*3) * (in[0]*6 + in[1]*-2) * (in[0]*3 - in[1])
}

// Sprinkles: capacity 2, durability 0, flavor -2, texture 0, calories 3
// Butterscotch: capacity 0, durability 5, flavor -3, texture 0, calories 3
// Chocolate: capacity 0, durability 0, flavor 5, texture -1, calories 8
// Candy: capacity 0, durability -1, flavor 0, texture 5, calories 8
// 500 = 3S + 3B + 8C + 8K

func fastScore(in []int) int {
	cal := (3*in[0] + 3*in[1] + 8*in[2] + 8*in[3])
	if cal != 500 {
		return 0
	}
	a := (2 * in[0])
	b := (5*in[1] - in[3])
	c := (-2*in[0] - 3*in[1] + 5*in[2])
	d := (-in[2] + 5*in[3])
	if a < 0 || b < 0 || c < 0 || d < 0 {
		return 0
	}
	return a * b * c * d
}

// divide num across a list (of 4)
func division(num int) [][]int {
	out := [][]int{}
	for i := 0; i <= 100; i++ {
		for j := 0; j <= 100; j++ {
			if i+j > 100 {
				break
			}
			for k := 0; k <= 100; k++ {
				if i+j+k > 100 {
					break
				}
				for l := 0; l <= 100; l++ {
					if l+k+j+i != 100 {
						continue
					}
					out = append(out, []int{i, j, k, l})
				}
			}
		}
	}
	return out
}

/*
100 0 0 0
99 0 0 1
99 0 1 0
99 1 0 0
98 0 1 1
98 1 0 1
98 1 1 0
97 1 1 1
96 1 1 2
96 1 2 1
96 2 1 1
95 1 2 2
95 2 1 2
95 2 2 1
94 2 2 2

*/

func setCounts(counts []int, in []*ingredient) {
	for i, count := range counts {
		in[i].count = count
	}
}

func score(in []*ingredient) int {
	var cap, dur, fla, tex int
	for _, ingredient := range in {
		cap += ingredient.cap * ingredient.count
		dur += ingredient.dur * ingredient.count
		fla += ingredient.fla * ingredient.count
		tex += ingredient.tex * ingredient.count
	}
	if cap < 0 || dur < 0 || fla < 0 || tex < 0 {
		return 0
	}
	fmt.Printf("%d * %d * %d * %d\n", cap, dur, fla, tex)
	return cap * dur * fla * tex
}

type ingredient struct {
	name  string
	cap   int
	dur   int
	fla   int
	tex   int
	cal   int
	count int
}

func newIng(name string, cap, dur, fla, tex, cal int) *ingredient {
	return &ingredient{
		name:  name,
		cap:   cap,
		dur:   dur,
		fla:   fla,
		tex:   tex,
		cal:   cal,
		count: 0,
	}
}

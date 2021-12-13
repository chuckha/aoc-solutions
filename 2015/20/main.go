package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	input, _ := strconv.Atoi(lines[0])

	for i := 1; ; i++ {
		//		fmt.Printf("House %d got %d presents\n", i, sum(factors(i))*10)
		if sum(factors(i))*11 >= input {
			fmt.Println(i)
			break
		}
	}

}

/*
8 -> 1, 2, 4

9 -> 1, 3, 9
8 -> 1, 2, 4, 8
7 -> 1, 7
6 -> 1, 2, 3, 6
20 -> 1, 2, 4, 5, 10

100
1 -> 1-50
2 -> 2->100
3 -> 3->150
4 -> 4->200



if factor * 50 < number

// all factors (that are not itself) are < sqrt

*/
func sum(n []int) int {
	out := 0
	for _, j := range n {
		out += j
	}
	return out
}

func factors(i int) []int {
	out := []int{}
	sqrt := int(math.Sqrt(float64(i)))
	for j := 1; j <= sqrt; j++ {
		if i%j == 0 {
			if j*50 >= i {
				out = append(out, j)
			}
			opposite := i / j
			if opposite != j {
				if opposite*50 >= i {
					out = append(out, opposite)
				}
			}
		}
	}
	return out
}

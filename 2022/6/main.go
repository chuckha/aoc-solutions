package main

import (
	"fmt"
	"os"

	"github.com/chuckha/aoc-solutions/internal"
)

const packetSize = 14

func main() {
	input := internal.ReadInput()[0]

	//O(n*m^2)
	// n == size of input; m == packet size
	for i := 0; i < len(input); {
		skip := overlap(input[i : i+packetSize])
		if skip == 0 {
			fmt.Println(i + packetSize)
			os.Exit(0)
		}
		fmt.Println("skipping", skip)
		i += skip
	}

	// O(n * m)
	// n == input size; m == number of unique items in input set
	// counts := signal{}
	// for i := 0; i < packetSize; i++ {
	// 	counts.add(input[i])
	// }
	// for i := packetSize; i < len(input); i++ {
	// 	counts.remove(input[i-packetSize])
	// 	counts.add(input[i])
	// 	if counts.done() {
	// 		fmt.Println(i + 1)
	// 		return
	// 	}
	// }
}

// return the left index that matched so we can skip that many in the next iteration
// e.g. annbcde -> annb -> nbcd is the next possible one
func overlap(items string) int {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i] == items[j] {
				return i + 1
			}
		}
	}
	return 0
}

// is first == second
// is first == third
// is first == fourth
// is second == third
// is second == fourth
// is third == fourth

type signal map[byte]int

// O(m) (m == unique characters in input)
func (s signal) done() bool {
	count := 0
	for _, v := range s {
		if v == 1 {
			count++
		}
	}
	return count == packetSize
}

func (s signal) add(b byte) {
	s[b]++
}
func (s signal) remove(b byte) {
	s[b]--
}

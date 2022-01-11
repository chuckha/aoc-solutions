package main

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	polymer := []byte(lines[0])
	fmt.Println(string(polymer))
	shortest := math.MaxInt
	shortestpolymer := []byte{}
	for i := 'a'; i < 'z'; i++ {
		fmt.Println(string(i))
		np := bytes.ReplaceAll(polymer, []byte(string(i)), []byte(""))
		np = bytes.ReplaceAll(np, bytes.ToUpper([]byte(string(i))), []byte(""))
		for i := 0; i < len(np)-1; i++ {
			if react(string(np[i]), string(np[i+1])) {
				np = append(np[:i], np[i+2:]...)
				//			fmt.Println(string(polymer))
				i = -1
			}
		}
		if len(np) < shortest {
			shortest = len(np)
			shortestpolymer = np
		}
	}
	fmt.Println(shortest)
	fmt.Println(string(shortestpolymer))

	// for i := 0; i < len(polymer)-1; i++ {
	// 	if react(string(polymer[i]), string(polymer[i+1])) {
	// 		polymer = append(polymer[:i], polymer[i+2:]...)
	// 		//			fmt.Println(string(polymer))
	// 		i = -1
	// 	}
	// }
	// fmt.Println(string(polymer))
	// fmt.Println(len(polymer))
}

func react(l, r string) bool {
	if l == r {
		return false
	}
	if strings.ToUpper(l) == r {
		return true
	}
	if strings.ToUpper(r) == l {
		return true
	}
	return false
}

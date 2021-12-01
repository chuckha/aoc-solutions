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

	// part 2
	location := &loc{}
	roboSanta := &loc{}

	temp := location
	count := map[[2]int]int{*location: 2}
	for i, c := range lines[0] {
		if i%2 == 1 {
			temp = roboSanta
		} else {
			temp = location
		}

		switch c {
		case '^':
			temp.n()
		case '>':
			temp.e()
		case 'v':
			temp.s()
		case '<':
			temp.w()
		}
		count[*temp]++
		//		fmt.Println(location)
	}
	fmt.Println(len(count))

	// part 1

	// location := &loc{}
	// count := map[[2]int]int{*location: 1}
	// for _, c := range lines[0] {
	// 	switch c {
	// 	case '^':
	// 		location.n()
	// 	case '>':
	// 		location.e()
	// 	case 'v':
	// 		location.s()
	// 	case '<':
	// 		location.w()
	// 	}
	// 	count[*location]++
	// 	//		fmt.Println(location)
	// }
	// fmt.Println(len(count))
}

type loc [2]int

func (l *loc) n() {
	l[1]++
}
func (l *loc) s() {
	l[1]--
}
func (l *loc) e() {
	l[0]++
}
func (l *loc) w() {
	l[0]--
}

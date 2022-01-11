package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	maxin := internal.ReadInput()[0]

	// max, _ := strconv.Atoi(maxin)
	//	necessary := max + 10
	r := &recipes{}
	// bootstrap
	r.add(3)
	r.add(7)
	track := &suffix{
		expected: maxin,
		current:  "",
	}
	for !track.done() {
		// if r.count > 20179083 {
		// 	panic("missed it somehow")
		// }
		// if r.count >= 20179081+6 {
		// 	fmt.Println(r.end.Prev.Prev.Prev.Prev.Prev.Prev.Data)
		// 	fmt.Println(r.end.Prev.Prev.Prev.Prev.Prev.Data)
		// 	fmt.Println(r.end.Prev.Prev.Prev.Prev.Data)
		// 	fmt.Println(r.end.Prev.Prev.Prev.Data)
		// 	fmt.Println(r.end.Prev.Prev.Data)
		// 	fmt.Println(r.end.Prev.Data)
		// 	fmt.Println(r.end.Data)
		// 	break
		// }

		//		fmt.Println(r.root.PrintFormat(r.elf1, r.elf2))
		nextScores := r.brew()
		for _, s := range nextScores {
			r.add(s)
			track.add(s)
			if track.done() {
				break
			}
			track.shorten()
		}
		r.moveElves()
	}
	fmt.Println(r.count - len(track.expected))
	// start := r.root.Forward(max)
	// for i := 0; i < 10; i++ {
	// 	fmt.Print(start.Data)
	// 	start = start.Next
	// }
	// fmt.Println()
}
func done(desired, have string) bool {
	return desired == have
}

type suffix struct {
	current  string
	expected string
}

func (s *suffix) done() bool {
	return s.current == s.expected
}
func (s *suffix) add(num int) {
	s.current += fmt.Sprintf("%d", num)
}
func (s *suffix) shorten() {
	for !strings.HasPrefix(s.expected, s.current) {
		s.current = s.current[1:]
	}
}

type recipes struct {
	// root since we never remove anything
	count int
	root  *internal.CircularLinkedList[int]
	end   *internal.CircularLinkedList[int]
	elf1  *internal.CircularLinkedList[int]
	elf2  *internal.CircularLinkedList[int]
}

func (r *recipes) brew() []int {
	value := r.elf1.Data + r.elf2.Data
	if value < 10 {
		return []int{value}
	}
	return []int{value / 10, value % 10}
}

func (r *recipes) moveElves() {
	dist1 := r.elf1.Data + 1
	r.elf1 = r.elf1.Forward(dist1)
	dist2 := r.elf2.Data + 1
	r.elf2 = r.elf2.Forward(dist2)
}

func (r *recipes) add(data int) {
	if r.count == 0 {
		l := internal.NewCircularLinkedList(data)
		r.root = l
		r.elf1 = l
		r.end = l
		r.count++
		return
	}
	if r.count == 1 {
		l := r.root.InsertAfter(data)
		r.elf2 = l
		r.count++
		r.end = l
		return
	}
	r.end = r.end.InsertAfter(data)
	r.count++
}

//  20179081
// 296773567

// 20179082 too high

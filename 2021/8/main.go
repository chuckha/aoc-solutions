package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

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
	// fmt.Println(" 000 ")
	// fmt.Println("1   2")
	// fmt.Println("1   2")
	// fmt.Println(" 333 ")
	// fmt.Println("4   5")
	// fmt.Println("4   5")
	// fmt.Println(" 666 ")
	displays := []*display{}
	sum := 0
	for _, line := range lines {
		sections := strings.Split(line, " | ")
		d := newDisplay()
		for _, word := range strings.Split(sections[0], " ") {
			d.info(strings.Split(word, ""))
		}
		// do something with actual output numbers
		// for _, word := range strings.Split(sections[1], " ") {
		// 	if len(word) == 2 || len(word) == 4 || len(word) == 3 || len(word) == 7 {
		// 		count++
		// 	}
		// }
		displays = append(displays, d)
		sum += d.outputValue(strings.Split(sections[1], " "))
	}
	fmt.Println(sum)
}

// SSD = seven segment display
type display struct {
	id       int
	segments map[string]*segment
}

func newDisplay() *display {
	letters := []string{"a", "b", "c", "d", "e", "f", "g"}
	segs := map[string]*segment{}
	for i, l := range letters {
		segs[string(l)] = newSegment(letters[i])
	}
	return &display{
		segments: segs,
	}
}

// in is "cfagb bag cbgfd fbagc"
func (d *display) outputValue(in []string) int {
	digits := []int{}
	for _, word := range in {
		num := []int{}
		for _, letter := range word {
			seg := d.segments[string(letter)]
			num = append(num, seg.locaysh)
		}
		sort.Sort(sort.IntSlice(num))
		digits = append(digits, translate(num))
	}
	return digits[0]*1000 + digits[1]*100 + digits[2]*10 + digits[3]
}

func (d *display) String() string {
	out := []string{}
	for _, seg := range d.segments {
		out = append(out, seg.String())
	}
	return strings.Join(out, "\n")
}

func (d *display) groupSegments(letters []string) ([]*segment, []*segment) {
	inGroup := make([]*segment, 0)
	outGroup := make([]*segment, 0)
	for name, seg := range d.segments {
		if internal.Search(name, letters) >= 0 {
			inGroup = append(inGroup, seg)
			continue
		}
		outGroup = append(outGroup, seg)
	}
	return inGroup, outGroup
}

// info takes something like 'ab'
func (d *display) info(in []string) {
	// if its two letters, it must be trying to display a 1
	if len(in) == 2 {
		inGroup, outGroup := d.groupSegments(in)
		for _, seg := range inGroup {
			seg.removePoss(0, 1, 3, 4, 6)
		}
		for _, seg := range outGroup {
			seg.removePoss(2, 5)
		}
	}

	if len(in) == 3 {
		inGroup, outGroup := d.groupSegments(in)
		for _, seg := range inGroup {
			seg.removePoss(1, 3, 4, 6)
		}
		for _, seg := range outGroup {
			seg.removePoss(0, 2, 5)
		}
	}

	if len(in) == 4 {
		inGroup, outGroup := d.groupSegments(in)
		for _, seg := range inGroup {
			seg.removePoss(0, 4, 6)
		}
		for _, seg := range outGroup {
			seg.removePoss(1, 2, 3, 5)
		}
	}

	if len(in) == 5 {
		_, outGroup := d.groupSegments(in)
		for _, seg := range outGroup {
			seg.removePoss(0, 3, 6)
		}
	}

	if len(in) == 6 {
		_, outGroup := d.groupSegments(in)
		for _, seg := range outGroup {
			seg.removePoss(0, 5, 6)
		}
	}

	for _, seg := range d.segments {
		if seg.locaysh != -1 {
			_, outGroup := d.groupSegments([]string{seg.letter})
			for _, s := range outGroup {
				s.removePoss(seg.locaysh)
			}
		}
	}

}

/*
 000
1   2
1   2
 333
4   5
4   5
 666
*/

func translate(in []int) int {
	if internal.EqualSlice(in, []int{0, 1, 2, 4, 5, 6}) {
		return 0
	}
	if internal.EqualSlice(in, []int{2, 5}) {
		return 1
	}
	if internal.EqualSlice(in, []int{0, 2, 3, 4, 6}) {
		return 2
	}
	if internal.EqualSlice(in, []int{0, 2, 3, 5, 6}) {
		return 3
	}
	if internal.EqualSlice(in, []int{1, 2, 3, 5}) {
		return 4
	}
	if internal.EqualSlice(in, []int{0, 1, 3, 5, 6}) {
		return 5
	}
	if internal.EqualSlice(in, []int{0, 1, 3, 4, 5, 6}) {
		return 6
	}
	if internal.EqualSlice(in, []int{0, 2, 5}) {
		return 7
	}
	if internal.EqualSlice(in, []int{0, 1, 2, 3, 4, 5, 6}) {
		return 8
	}
	if internal.EqualSlice(in, []int{0, 1, 2, 3, 5, 6}) {
		return 9
	}
	panic(fmt.Sprintf("bad input %v", in))
}

type segment struct {
	letter            string
	locaysh           int
	possibleLocations map[int]struct{}
}

func newSegment(l string) *segment {
	return &segment{
		letter:  l,
		locaysh: -1,
		possibleLocations: map[int]struct{}{
			0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {},
		},
	}
}

func (s *segment) removePoss(ints ...int) {
	for _, j := range ints {
		delete(s.possibleLocations, j)
	}
	if len(s.possibleLocations) == 1 {
		for k := range s.possibleLocations {
			s.locaysh = k
		}
	}
}

func (s *segment) String() string {
	return fmt.Sprintf("%s: %v", s.letter, s.possibleLocations)
}

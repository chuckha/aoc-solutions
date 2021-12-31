package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	input := "abcdefgh"
	//	input = "abcde"
	input = "decab"
	input = "fbgdceah"
	data := []string{}
	for _, c := range input {
		data = append(data, string(c))
	}

	reverse := true

	if reverse {
		lines = rev(lines)
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	instructions := []instruction{}
	for _, line := range lines {
		words := strings.Split(line, " ")
		// swap letter e with letter h
		if words[0] == "swap" && words[1] == "letter" {
			instructions = append(instructions, swapLetter{words[2], words[5]})
			continue
		}
		// move position 6 to position 3
		if words[0] == "move" && words[1] == "position" {
			from, _ := strconv.Atoi(words[2])
			to, _ := strconv.Atoi(words[5])
			instructions = append(instructions, movePosition{from, to, reverse})
			continue
		}
		// reverse positions 1 through 6
		if words[0] == "reverse" && words[1] == "positions" {
			from, _ := strconv.Atoi(words[2])
			to, _ := strconv.Atoi(words[4])
			instructions = append(instructions, reversePositions{from, to}) // there is no reverse for this op
			continue
		}
		// rotate right 6 steps
		if words[0] == "rotate" && words[1] == "right" {
			num, _ := strconv.Atoi(words[2])
			instructions = append(instructions, rotateRight{num, reverse})
			continue
		}
		// rotate based on position of letter e
		if words[0] == "rotate" && words[1] == "based" {
			instructions = append(instructions, rotateBasedOnLetter{words[6], reverse})
			continue
		}
		// swap position 7 with position 2
		if words[0] == "swap" && words[1] == "position" {
			from, _ := strconv.Atoi(words[2])
			to, _ := strconv.Atoi(words[5])
			instructions = append(instructions, swapPositions{from, to})
			continue
		}
		// rotate left 0 steps
		if words[0] == "rotate" && words[1] == "left" {
			num, _ := strconv.Atoi(words[2])
			instructions = append(instructions, rotateLeft{num, reverse})
			continue
		}
		panic("handle this line: " + line)
	}
	for _, inst := range instructions {
		fmt.Println(inst)
		//		fmt.Println("before", data)
		inst.run(data)
		//fmt.Println("after", data)
		fmt.Println(strings.Join(data, ""))
	}
}

type instruction interface {
	run([]string)
	String() string
}
type rotateLeft struct {
	num     int
	reverse bool
}

func (r rotateLeft) run(in []string) {
	if r.reverse {
		r.runReverse(in)
		return
	}
	out := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		x := i - r.num
		if x < 0 {
			x = len(in) + x
		}
		out[x] = in[i]
	}
	copy(in, out)
}
func (r rotateLeft) runReverse(in []string) {
	rr := rotateRight{r.num, false}
	rr.run(in)
}
func (r rotateLeft) String() string {
	return fmt.Sprintf("rotating left %d places", r.num)
}

type swapPositions struct {
	from, to int
}

func (s swapPositions) run(in []string) {
	in[s.from], in[s.to] = in[s.to], in[s.from]
}
func (s swapPositions) String() string {
	return fmt.Sprintf("swapping positions from %d to %d", s.from, s.to)
}

type rotateBasedOnLetter struct {
	letter  string
	reverse bool
}

func (r rotateBasedOnLetter) run(in []string) {
	if r.reverse {
		r.runReverse(in)
		return
	}
	// rotate it once, then rotate by the index, then once more if idx >= 4
	idx := internal.Search(r.letter, in)
	rotate := rotateRight{1, false}
	rotate.run(in)
	rotate = rotateRight{idx, false}
	rotate.run(in)
	if idx >= 4 {
		rotate = rotateRight{1, false}
		rotate.run(in)
	}
}
func (r rotateBasedOnLetter) runReverse(in []string) {
	data := make([]string, len(in))
	copy(data, in)
	for {
		// rotate in once left
		// run rotate based on letter
		// see if it equals the in
		rl := rotateLeft{1, false}
		rl.run(data)
		tryit := make([]string, len(data))
		copy(tryit, data)
		rb := rotateBasedOnLetter{letter: r.letter}
		rb.run(tryit)
		if internal.EqualSlice(tryit, in) {
			copy(in, data)
			return
		}
	}
}
func (r rotateBasedOnLetter) String() string {
	return fmt.Sprintf("rotating based on letter %q", r.letter)
}

type rotateRight struct {
	num     int
	reverse bool
}

func (r rotateRight) run(in []string) {
	if r.reverse {
		r.runReverse(in)
		return
	}
	out := make([]string, len(in))
	for i := 0; i < len(in); i++ {
		out[(i+r.num)%len(in)] = in[i]
	}
	copy(in, out)
}
func (r rotateRight) runReverse(in []string) {
	rl := rotateLeft{r.num, false}
	rl.run(in)
}
func (r rotateRight) String() string {
	return fmt.Sprintf("rotating right %d places", r.num)
}

type reversePositions struct {
	from, to int
}

func (r reversePositions) run(in []string) {
	i := r.from
	j := r.to
	for i <= j {
		in[i], in[j] = in[j], in[i]
		i++
		j--
	}
}
func (r reversePositions) String() string {
	return fmt.Sprintf("reversing from %d to %d", r.from, r.to)
}

type swapLetter struct {
	first, second string
}

func (s swapLetter) run(in []string) {
	idx1 := internal.Search(s.first, in)
	idx2 := internal.Search(s.second, in)
	in[idx1], in[idx2] = in[idx2], in[idx1]
}

func (s swapLetter) String() string {
	return fmt.Sprintf("swap %s with %s", s.first, s.second)
}

// move position 6 to position 3

type movePosition struct {
	from, to int
	reverse  bool
}

func (m movePosition) run(in []string) {
	if m.reverse {
		m.runReverse(in)
		return
	}
	letter := in[m.from]
	without := []string{}
	for i := 0; i < len(in); i++ {
		if i == m.from {
			continue
		}
		without = append(without, in[i])
	}
	out := []string{}
	for i := 0; i < len(in); i++ {
		if i < m.to {
			out = append(out, without[i])
		}
		if i == m.to {
			out = append(out, letter)
		}
		if i > m.to {
			out = append(out, without[i-1])
		}
	}
	copy(in, out)
}
func (m movePosition) runReverse(in []string) {
	mp := movePosition{m.to, m.from, false}
	mp.run(in)
}
func (m movePosition) String() string {
	return fmt.Sprintf("moving position %d to position %d", m.from, m.to)
}

// swap position 4 with position 0
// swap letter d with letter b
// reverse positions 0 through 4
// rotate left 1
// move position 1 to position 4
// move position 3 to position 0
// rotate based on position of letter b
// rotate based on position of letter d

func rev(in []string) []string {
	out := make([]string, len(in))
	j := 0
	for i := len(in) - 1; i >= 0; i-- {
		out[j] = in[i]
		j++
	}
	return out
}

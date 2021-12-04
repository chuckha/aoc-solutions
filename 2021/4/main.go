package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		data := scanner.Text()
		data = strings.Replace(data, "  ", " ", -1)
		data = strings.TrimSpace(data)
		lines = append(lines, data)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	inputs := lines[0]
	lines = lines[2:]
	boards := []*board{}
	for i := 0; (i * 6) < len(lines); i++ {
		boards = append(boards, newBoard(lines[i*6:(i*6)+5]))
	}
	ins := parseInputNumbers(inputs)
	numwon := 0
	for _, in := range ins {
		for _, b := range boards {
			if b.won() {
				continue
			}
			b.call(in)
			if b.won() {
				numwon++
			}
			if numwon == len(boards) {
				fmt.Println(b.score(in))
				return
			}
		}
	}
}

func parseInputNumbers(line string) []int {
	nums := strings.Split(line, ",")
	out := make([]int, len(nums))
	for i, n := range nums {
		d, _ := strconv.Atoi(n)
		out[i] = d
	}
	return out
}

// position is i,j
type cell struct {
	position [2]int
	number   int
	called   bool
}

func newCell(pos [2]int, num int) *cell {
	return &cell{
		position: pos,
		number:   num,
		called:   false,
	}
}

type board struct {
	numbers   map[int]*cell
	positions [5][5]*cell
}

func newBoard(in []string) *board {
	b := &board{
		numbers:   map[int]*cell{},
		positions: [5][5]*cell{},
	}
	for i, line := range in {
		nums := strings.Split(line, " ")
		for j, num := range nums {
			if strings.TrimSpace(num) == "" {
				continue
			}
			c, _ := strconv.Atoi(num)
			cell := newCell([2]int{i, j}, c)
			b.numbers[c] = cell
			b.positions[i][j] = cell
		}
	}
	return b
}

func (b *board) String() string {
	out := ""
	called := ""
	for _, row := range b.positions {
		rs := []string{}
		c := []string{}
		for _, cell := range row {
			rs = append(rs, fmt.Sprintf("%v", cell.number))
			c = append(c, fmt.Sprintf("%v", cell.called))
		}
		out += strings.Join(rs, " ") + "\n"
		called += strings.Join(c, " ") + "\n"
	}
	return out + "\n" + called
}

func (b *board) call(i int) {
	cell, ok := b.numbers[i]
	if !ok {
		return
	}
	cell.called = true
}

func (b *board) won() bool {
	return b.acrossVictories() || b.downVictories()
}

func (b *board) acrossVictories() bool {
	for i := 0; i < 5; i++ {
		all := true
		for _, cell := range b.positions[i] {
			all = all && cell.called
			if !all {
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}

func (b *board) downVictories() bool {
	for j := 0; j < 5; j++ {
		all := true
		for i := 0; i < 5; i++ {
			all = all && b.positions[i][j].called
			if !all {
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}

func (b *board) score(winningNum int) int {
	sum := 0
	for _, n := range b.numbers {
		if n.called {
			continue
		}
		sum += n.number
	}
	return sum * winningNum
}

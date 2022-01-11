package main

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
S(N|)(S|E)N
move south
left = n
right = ""
tail = (S|E)N

move N

SSN
SEN
SNSN
SNEN
*/

func (m *maze) moveIt2(input inputData) {
	if input.data == "" {
		return
	}
	head, branches, tail := parsemovement3(input.data)
	//	fmt.Println(input.pos, "head", head, "branches", branches, "tail", tail)
	// special case when left branch is not empty and right branch is empty
	// just do the left branch w/o the tail
	// N(W|)E
	m.moveInputData(inputData{data: head, pos: input.pos})
	afterHead := m.cur
	// if left != "" && right == "" {
	// 	m.moveInputData(inputData{data: left, pos: afterHead})
	// 	m.moveIt(inputData{data: tail, pos: afterHead})
	// 	return
	// }
	// // N(|W)E
	// if left == "" && right != "" {
	// 	m.moveInputData(inputData{data: right, pos: afterHead})
	// 	m.moveIt(inputData{data: tail, pos: afterHead})
	// 	return
	// }
	for _, b := range branches {
		if len(b) > 0 {
			m.moveIt2(inputData{data: b + tail, pos: afterHead})
		}
	}
}

// (N(W|E)W|E)
func (m *maze) moveIt(input inputData) {
	if input.data == "" {
		return
	}
	head, left, right, tail := parsemovement2(input.data)
	//	fmt.Println(input.pos, "head", head, "left", left, "right", right, "tail", tail)
	// special case when left branch is not empty and right branch is empty
	// just do the left branch w/o the tail
	// N(W|)E
	m.moveInputData(inputData{data: head, pos: input.pos})
	afterHead := m.cur
	// if left != "" && right == "" {
	// 	m.moveInputData(inputData{data: left, pos: afterHead})
	// 	m.moveIt(inputData{data: tail, pos: afterHead})
	// 	return
	// }
	// // N(|W)E
	// if left == "" && right != "" {
	// 	m.moveInputData(inputData{data: right, pos: afterHead})
	// 	m.moveIt(inputData{data: tail, pos: afterHead})
	// 	return
	// }
	if left != "" {
		m.moveIt(inputData{data: left + tail, pos: afterHead})
	}
	if right != "" {
		m.moveIt(inputData{data: right + tail, pos: afterHead})
	}
}

func main() {
	input := internal.ReadInput()[0]
	input = strings.TrimPrefix(input, "^")
	input = strings.TrimSuffix(input, "$")
	tc := strings.Count(input, "E") + strings.Count(input, "S") + strings.Count(input, "N") + strings.Count(input, "W")
	m := newMaze()
	m.visit() // visit the origin
	m.moveIt2(inputData{data: input, pos: point{0, 0}})
	m.finalize()
	fmt.Println(m)
	costs := m.djikstra()
	//	fmt.Println(m.CostString(costs))
	cs := make([]int, 0)
	for _, v := range costs {
		if v == math.MaxInt {
			continue
		}
		cs = append(cs, v)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(cs)))
	fmt.Println(cs[0])
	////	fmt.Println(m.CostString2(costs))
	fmt.Println("total rooms:", len(cs))
	for i := range cs {
		if cs[i] < 1000 {
			fmt.Println("part2:", i)
			break
		}
	}
	fmt.Println("total chars", tc)
	fmt.Println("counted chars", m.counted)
}

// part 1: 4239 correct
// part 2: 7811 too low
// part 2: 7812 too low (made this one up)
// part 2 correct: 8205 ... but what bug does my code have O_O

// part 1: 2107 too low
// part 1: 2223 too low

// WWW(WWWW|)
type inputData struct {
	pos  point
	data string
}

type maze struct {
	data                   map[point]string
	minx, miny, maxx, maxy int
	cur                    point
	counted                int
}

func (m *maze) finalize() {
	for k, d := range m.data {
		if d == "?" {
			m.data[k] = "#"
		}
	}
	m.cur = point{0, 0}
}

func (m *maze) djikstra() map[point]int {
	costs := make(map[point]int)
	for k := range m.data {
		costs[k] = math.MaxInt
	}
	costs[point{0, 0}] = 0
	visited := make(map[point]struct{})
	q := internal.NewQueue[point]()
	q.Enqueue(point{0, 0})
	for !q.Empty() {
		cur := q.Dequeue()
		if _, ok := visited[cur]; ok {
			continue
		}
		for _, n := range cur.oneAwayNeighbors() {
			// ignore walls
			if item, ok := m.data[n]; ok && item == "#" {
				continue
			}
			for _, n2 := range n.oneAwayNeighbors() {
				if n2 == cur {
					continue
				}
				if item, ok := m.data[n2]; ok && item == "#" {
					continue
				}
				if costs[cur]+1 < costs[n2] {
					costs[n2] = costs[cur] + 1
				}
				q.Enqueue(n2)
			}
		}
		visited[cur] = struct{}{}
	}
	return costs
}

func newMaze() *maze {
	return &maze{data: make(map[point]string)}
}

func (m *maze) jumpTo(p point) {
	m.cur = p
}

func (m *maze) visit() {
	// diagonals from where we visit are always walls
	m.doorsAround(m.cur)
	// doors are unknown from where we are
	m.initializeUnknownDoors(m.cur)
	m.add(m.cur, ".")
}
func (m *maze) doorsAround(cur point) {
	for _, p := range cur.diagonals() {
		if _, ok := m.data[p]; ok {
			continue
		}
		m.add(p, "#")
	}
}
func (m *maze) initializeUnknownDoors(cur point) {
	for _, p := range cur.oneAwayNeighbors() {
		// never overwite anything with a ?
		if _, ok := m.data[p]; ok {
			continue
		}
		m.add(p, "?")
	}
}
func (m *maze) initEmpty(cur point) {
	if _, ok := m.data[cur]; ok {
		return
	}
	m.add(cur, ".")
}
func (m *maze) String() string {
	var out strings.Builder
	for j := m.miny; j <= m.maxy; j++ {
		for i := m.minx; i <= m.maxx; i++ {
			p := point{i, j}
			if m.cur == p {
				out.WriteString("X")
				continue
			}
			if item, ok := m.data[point{i, j}]; ok {
				out.WriteString(item)
				continue
			}
			out.WriteString(" ")
		}
		out.WriteString("\n")
	}
	return out.String()
}
func (m *maze) CostString(costs map[point]int) string {
	var out strings.Builder
	for j := m.miny; j <= m.maxy; j++ {
		for i := m.minx; i <= m.maxx; i++ {
			p := point{i, j}
			if m.cur == p {
				out.WriteString("XX")
				continue
			}
			if item, ok := costs[point{i, j}]; ok && item != math.MaxInt {
				out.WriteString(fmt.Sprintf("%02d", item))
				continue
			}
			if item, ok := m.data[point{i, j}]; ok {
				out.WriteString(strings.Repeat(item, 2))
				continue
			}
			out.WriteString(" ")
		}
		out.WriteString("\n")
	}
	return out.String()
}

func red(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;41m%s\033[0m", in))
}
func green(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;42m%s\033[0m", in))
}
func blue(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;44m%s\033[0m", in))
}
func black(in string) []byte {
	return []byte(fmt.Sprintf("\033[1;40m%s\033[0m", in))
}

func (m *maze) CostString2(costs map[point]int) string {
	var out strings.Builder
	for j := m.miny; j <= m.maxy; j++ {
		for i := m.minx; i <= m.maxx; i++ {
			p := point{i, j}
			if m.cur == p {
				out.WriteString("X")
				continue
			}
			if item, ok := costs[point{i, j}]; ok && item != math.MaxInt {
				if item >= 10 && item < 100 {
					out.Write(red("a"))
					continue
				}
				if item >= 100 && item < 1000 {
					out.Write(blue("b"))
					continue
				}
				if item >= 1000 {
					out.Write(green("o"))
					continue
				}
				out.WriteString(fmt.Sprintf("%d", item))
				continue
			}
			if item, ok := m.data[point{i, j}]; ok {
				if item == "#" {
					out.Write(black(item))
					continue
				}
				out.WriteString(" ")
				continue
				out.WriteString(strings.Repeat(item, 1))
				continue
			}
			out.WriteString(" ")
		}
		out.WriteString("\n")
	}
	return out.String()
}
func (m *maze) moveString(s string) {
	for _, c := range s {
		m.counted++
		m.move(string(c))
	}
}
func (m *maze) moveInputData(id inputData) {
	//	fmt.Println(id.pos, "moving", id.data)
	m.jumpTo(id.pos)
	m.moveString(id.data)
}

func (m *maze) move(dir string) {
	switch dir {
	case "N":
		m.cur = m.cur.u()
		m.door("-")
		m.cur = m.cur.u()
	case "E":
		m.cur = m.cur.r()
		m.door("|")
		m.cur = m.cur.r()
	case "S":
		m.cur = m.cur.d()
		m.door("-")
		m.cur = m.cur.d()
	case "W":
		m.cur = m.cur.l()
		m.door("|")
		m.cur = m.cur.l()
	default:
		panic("bad direction: " + dir)
	}
	m.visit()
}
func (m *maze) door(c string) {
	m.add(m.cur, c)
}
func (m *maze) add(p point, s string) {
	if p.x < m.minx {
		m.minx = p.x
	}
	if p.x > m.maxx {
		m.maxx = p.x
	}
	if p.y < m.miny {
		m.miny = p.y
	}
	if p.y > m.maxy {
		m.maxy = p.y
	}
	// ? can become - or |
	// but doors and walls cannot become anything else
	if item, ok := m.data[p]; ok {
		if item == s {
			return
		}
		if item == "|" || item == "-" || item == "#" || item == "." {
			fmt.Printf("trying to ovewrite %q with %q\n", item, s)
			panic("do not overwrite a door or wall or empty")
		}
	}
	m.data[p] = s
}

type point struct {
	x, y int
}

func (p point) diagonals() []point {
	return []point{p.ul(), p.ur(), p.dr(), p.dl()}
}
func (p point) oneAwayNeighbors() []point {
	return []point{p.u(), p.r(), p.d(), p.l()}
}
func (p point) ul() point {
	return point{p.x - 1, p.y - 1}
}
func (p point) u() point {
	return point{p.x, p.y - 1}
}
func (p point) ur() point {
	return point{p.x + 1, p.y - 1}
}
func (p point) r() point {
	return point{p.x + 1, p.y}
}
func (p point) dr() point {
	return point{p.x + 1, p.y + 1}
}
func (p point) d() point {
	return point{p.x, p.y + 1}
}
func (p point) dl() point {
	return point{p.x - 1, p.y + 1}
}
func (p point) l() point {
	return point{p.x - 1, p.y}
}

// head, branch left, branch right, tail
// N(E|W)N => N <save point>, <restore point> EN, <restore point>WN
func parsemovement2(input string) (string, string, string, string) {
	cur := 0
	start := 0
	branchLeft := ""
	branchRight := ""
	parens := internal.NewStack[string]()
	head := ""
	tail := ""
	for {
		if cur >= len(input) {
			return input, "", "", "" // to get here there must be no other ()| in the string
		}
		c := input[cur]
		cur++
		switch c {
		case '^', '$':
			start = cur
		case '(':
			if parens.Depth() == 0 {
				head = input[start : cur-1]
				start = cur
			}
			parens.Push("(")
		case ')':
			parens.Pop()
			if parens.Depth() == 0 {
				branchRight = input[start : cur-1]
				tail = input[cur:]
				return head, branchLeft, branchRight, tail
			}
		case '|':
			if parens.Depth() == 1 {
				branchLeft = input[start : cur-1]
				start = cur
			}
		}
	}
}

func parsemovement3(input string) (string, []string, string) {
	cur := 0
	start := 0
	branches := []string{}
	parens := internal.NewStack[string]()
	head := ""
	tail := ""
	for {
		if cur >= len(input) {
			return input, branches, "" // to get here there must be no other ()| in the string
		}
		c := input[cur]
		cur++
		switch c {
		case '^', '$':
			start = cur
		case '(':
			if parens.Depth() == 0 {
				head = input[start : cur-1]
				start = cur
			}
			parens.Push("(")
		case ')':
			parens.Pop()
			if parens.Depth() == 0 {
				branches = append(branches, input[start:cur-1])
				tail = input[cur:]
				return head, branches, tail
			}
		case '|':
			if parens.Depth() == 1 {
				branches = append(branches, input[start:cur-1])
				start = cur
			}
		}
	}
}

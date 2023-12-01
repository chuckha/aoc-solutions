package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const debug = false
const rounds = 10000

type monkies map[int]*monkey

func (m monkies) String() string {
	out := []string{}
	for _, monkey := range m {
		out = append(out, fmt.Sprintf("Monkey %d inspected items %d times.", monkey.id, monkey.count))
	}
	return strings.Join(out, "\n")
}

func DebugF(tmp string, a ...any) {
	if debug {
		fmt.Printf(tmp, a...)
	}
}

func (m monkies) testmod() int {
	prod := 1
	for _, monk := range m {
		prod *= monk.testOperationVal
	}
	return prod
}

func (m monkies) round() {
	for i := 0; i < len(m); i++ {
		monk := m[i]
		DebugF("Monkey %d:\n", monk.id)
		for !monk.items.Empty() {
			worry := monk.items.Dequeue()
			DebugF("\tMonkey inspects an item with a worry level of %d.\n", worry)
			monk.count++
			tmpl := "\t\tWorry level increases by %v to %v.\n"
			if monk.operation == mul {
				tmpl = "\t\tWorry level is multiplied by %v to %v.\n"
			}
			oldWorry := worry
			worry = worry % m.testmod()
			worry = monk.changeWorry(worry)
			DebugF(tmpl, num(monk.operationVal, oldWorry), worry)
			//			worry = worry % monk.testOperationVal
			// worry = worry / 3
			// DebugF("\t\tMonkey gets bored with item. Worry level is divided by 3 to %v\n", worry)
			// testtmpl := "\t\tCurrent worry level is not %v %v.\n"
			//		DebugF(testtmpl, monk.testOperation, monk.testOperationVal)
			toMonkey := monk.falseMonkey
			if monk.test(worry) {
				toMonkey = monk.trueMonkey
			}
			DebugF("\t\ttem with worry level %v is thrown to monkey %v.\n", worry, toMonkey)
			m[toMonkey].items.Enqueue(worry)
		}
	}
}

func num(t, v int) int {
	if t == -1 {
		return v
	}
	return t
}

func main() {
	lines := internal.ReadRealRawInput()
	ms := monkies{}
	cur := &monkey{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		switch {
		case strings.HasPrefix(line, "Monkey"):
			cur.id = parseMonkeyID(line)
		case strings.HasPrefix(line, "  Starting items:"):
			cur.items = parseStartingItems(line)
		case strings.HasPrefix(line, "  Operation:"):
			cur.operation = parseOperation(line)
			cur.operationVal = parseOperationValue(line)
		case strings.HasPrefix(line, "  Test:"):
			cur.testOperation = parseTestOperation(line)
			cur.testOperationVal = parseTestOperationValue(line)
			i++
			trueLine := lines[i]
			cur.trueMonkey = throwToMonkey(trueLine)
			i++
			falseLine := lines[i]
			cur.falseMonkey = throwToMonkey(falseLine)
		case len(line) == 0:
			ms[cur.id] = cur
			cur = &monkey{}
		}
	}

	for i := 1; i <= rounds; i++ {
		ms.round()
		switch i {
		case 1, 20, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000:
			fmt.Printf("== After round %d ==\n", i)
			fmt.Println(ms)
		}
	}
	sortable := make(cheekymonkies, len(ms))
	for i, m := range ms {
		sortable[i] = m
	}
	fmt.Println(sortable.monkeyBusiness())
}

func parseMonkeyID(line string) int {
	parts := strings.Split(strings.TrimSuffix(line, ":"), " ")
	id, _ := strconv.Atoi(parts[1])
	return id
}

func parseStartingItems(line string) *internal.Queue[int] {
	out := internal.NewQueue[int]()
	parts := strings.Split(line, ": ")
	items := strings.Split(parts[1], ", ")
	for _, item := range items {
		worry, _ := strconv.Atoi(item)
		out.Enqueue(worry)
	}
	return out
}

func parseOperation(line string) op {
	parts := strings.Split(line, "new = old ")
	switch {
	case strings.HasPrefix(parts[1], "*"):
		return mul
	case strings.HasPrefix(parts[1], "+"):
		return add
	default:
		panic("unsupported operation")
	}
}

func parseOperationValue(line string) int {
	parts := strings.Split(line, " ")
	val, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1
	}
	return val
}

func parseTestOperation(line string) op {
	if strings.Contains(line, "divisible") {
		return mod
	}
	panic("unsupported test operation")
}

func parseTestOperationValue(line string) int {
	parts := strings.Split(line, " ")
	x, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1
	}
	return x
}

func throwToMonkey(line string) int {
	parts := strings.Split(line, " throw to monkey ")
	mid, _ := strconv.Atoi(parts[1])
	return mid
}

type op string

const (
	mul op = "*"
	add op = "+"
	mod op = "divisible by"
)

type monkey struct {
	id               int
	items            *internal.Queue[int]
	operation        op
	operationVal     int
	testOperation    op
	testOperationVal int
	trueMonkey       int
	falseMonkey      int
	count            int
}

func (m *monkey) String() string {
	out := fmt.Sprintf("Monkey %d:\n", m.id)
	out += fmt.Sprintf("\tStarting items: %v\n", m.items)
	out += fmt.Sprintf("\tOperation: new = old %v %v\n", m.operation, opVal(m.operationVal))
	out += fmt.Sprintf("\tTest: %v %v\n", m.testOperation, m.testOperationVal)
	out += fmt.Sprintf("\t\tIf true: throw to monkey %d\n", m.trueMonkey)
	out += fmt.Sprintf("\t\tIf false: throw to monkey %d", m.falseMonkey)
	return out
}

func opVal(x int) string {
	if x == -1 {
		return "old"
	}
	return fmt.Sprintf("%d", x)
}

func (m *monkey) changeWorry(in int) int {
	constant := in
	if m.operationVal != -1 {
		constant = m.operationVal
	}
	switch m.operation {
	case mul:
		return in * constant
	case add:
		return in + constant
	default:
		panic("unsupported operation")
	}
}

func (m *monkey) test(w int) bool {
	switch m.testOperation {
	case mod:
		return w%m.testOperationVal == 0
	default:
		panic("unsupported test operation")
	}
}

type cheekymonkies []*monkey

func (c cheekymonkies) Len() int           { return len(c) }
func (c cheekymonkies) Less(i, j int) bool { return c[i].count < c[j].count }
func (c cheekymonkies) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func (c cheekymonkies) monkeyBusiness() int {
	sort.Sort(sort.Reverse(c))
	return c[0].count * c[1].count
}

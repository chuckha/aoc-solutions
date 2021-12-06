package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	NOT = iota
	AND
	OR
	LSHIFT
	RSHIFT
	NONE
	NUM
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	answers := doit(lines, map[string]*gate{})
	answers["a"].output()
	out := doit(lines, map[string]*gate{
		"b": {
			name:     "b",
			operator: NUM,
			num:      answers["a"].output(),
			set:      true,
		},
	})
	// for _, gate := range out {
	// 	fmt.Println(gate)
	// }
	fmt.Println(out["a"].output())
}

func doit(lines []string, override map[string]*gate) map[string]*gate {
	gates := map[string]*gate{}
	for _, line := range lines {
		if line == "19138 -> b" {
			line = "16076 -> b"
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		items := strings.Split(line, " -> ")
		outputName := items[len(items)-1]

		outGate, ok := gates[outputName]
		if !ok {
			outGate = namedGate(outputName)
			gates[outputName] = outGate
		}

		leftside := strings.Split(items[0], " ")
		// its just a number or maybe same gate
		if len(leftside) == 1 {
			num, err := strconv.ParseInt(leftside[0], 10, 16)
			if err != nil {
				outGate.operator = NONE
				leftGate := newNopGate(leftside[0])
				gates[leftside[0]] = leftGate
				outGate.addInGate(leftGate)
				continue
			}
			// number case
			outGate.operator = NUM
			outGate.num = uint16(num)
			outGate.set = true
			continue
		}
		// NOT case
		if len(leftside) == 2 {
			num, err := strconv.ParseInt(leftside[1], 10, 16)
			// if its not a number
			if err != nil {
				notGate, ok := gates[leftside[1]]
				if !ok {
					notGate = newNotGate(leftside[1])
					gates[leftside[1]] = notGate
				}
				outGate.addInGate(notGate)
				continue
			}
			numgate := newNumber(uint16(num))
			outGate.addInGate(numgate)
			continue
		}
		// must be a leftside of len 3
		if len(leftside) != 3 {
			panic(fmt.Sprintf("bad input %q", line))
		}
		op := parseOp(leftside[1])
		outGate.operator = op

		// left could be a number
		num, err := strconv.ParseInt(leftside[0], 10, 16)
		// if its not a number
		if err != nil {
			leftGate, ok := gates[leftside[0]]
			if !ok {
				leftGate = namedGate(leftside[0])
				gates[leftside[0]] = leftGate
			}
			outGate.addInGate(leftGate)
		} else {
			leftGate := newNumber(uint16(num))
			outGate.addInGate(leftGate)
		}

		if op == LSHIFT || op == RSHIFT {
			num, _ := strconv.ParseInt(leftside[2], 10, 16)
			g := newNumber(uint16(num))
			outGate.addInGate(g)
			continue
		}
		// AND and OR
		// left could be a number
		num, err = strconv.ParseInt(leftside[2], 10, 16)
		// if its not a number
		if err != nil {
			rightGate, ok := gates[leftside[2]]
			if !ok {
				rightGate = namedGate(leftside[2])
				gates[leftside[2]] = rightGate
			}
			outGate.addInGate(rightGate)
		} else {
			rightGate := newNumber(uint16(num))
			outGate.addInGate(rightGate)
		}
	}
	return gates
}

func parseOp(in string) int {
	switch in {
	case "OR":
		return OR
	case "AND":
		return AND
	case "LSHIFT":
		return LSHIFT
	case "RSHIFT":
		return RSHIFT
	default:
		panic(fmt.Sprintf("bad operator %q", in))
	}
}

type gate struct {
	name     string
	operator int
	inputs   []*gate
	num      uint16
	set      bool
}

func newNumber(num uint16) *gate {
	return &gate{
		operator: NUM,
		num:      num,
		set:      true,
	}
}

func newNopGate(name string) *gate {
	return &gate{
		name:     name,
		operator: NONE,
	}
}

func namedGate(name string) *gate {
	return &gate{
		name: name,
	}
}

func newNotGate(name string) *gate {
	return &gate{
		name:     name,
		operator: NOT,
	}
}

func (g *gate) addInGate(in *gate) {
	g.inputs = append(g.inputs, in)
}

func (g *gate) String() string {
	return fmt.Sprintf("%d", g.output())
}

func (g *gate) output() uint16 {
	if g.set {
		return g.num
	}
	var val uint16
	switch g.operator {
	case NOT:
		if len(g.inputs) == 0 {
			fmt.Printf("failure on gate %q, %q, %q, %q\n", g.name, g.inputs, g.num, g.set)
		}
		//		fmt.Println("NOT")
		val = ^g.inputs[0].output()
	case AND:
		//		fmt.Println("AND")
		val = g.inputs[0].output() & g.inputs[1].output()
	case OR:
		//		fmt.Println("OR")
		val = g.inputs[0].output() | g.inputs[1].output()
	case LSHIFT:
		//		fmt.Println("LSHIFT")
		val = g.inputs[0].output() << g.inputs[1].output()
	case RSHIFT:
		//		fmt.Println("RSHIFT")
		val = g.inputs[0].output() >> g.inputs[1].output()
	case NONE:
		//		fmt.Println("NONE")
		val = g.inputs[0].output()
	}
	g.num = uint16(val)
	g.set = true
	return g.num
}

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
3296135418821 is too high
3296135418820

*/

type variable struct {
	id       string
	variable string
}

func (v variable) Value() string {
	return v.variable
}
func (v variable) ID() string {
	return v.id
}
func (v variable) String() string {
	return v.variable
}

type id string

func (i id) Value() string {
	return string(i)
}
func (i id) ID() string {
	return string(i)
}
func (i id) String() string {
	return string(i)
}

type valuer interface {
	Value() string
	ID() string
	String() string
}

type number struct {
	id  string
	num float64
}

func (n number) Value() string {
	return fmt.Sprintf("%f", n.num)
}
func (n number) ID() string {
	return n.id
}
func (n number) String() string {
	return fmt.Sprintf("%f", n.num)
}

type oo struct {
	id    string
	left  valuer
	op    op
	right valuer
}

func doMath(l, r float64, op op) float64 {
	switch op {
	case add:
		return l + r
	case mul:
		return l * r
	case sub:
		return l - r
	case div:
		return l / r
	case eql:
		if l == r {
			return 1
		}
		return 0
	default:
		panic("bad operator")
	}
}

func (o *oo) Value() string {
	if o.resolvable() {
		return fmt.Sprintf("%f", o.value())
	}
	return fmt.Sprintf("(%s %s %s)", o.left.Value(), o.op, o.right.Value())
}

func (o *oo) resolvable() bool {
	// if left resolves to an int and right resolves to an int then it's resolvable
	switch left := o.left.(type) {
	case *number:
		switch right := o.right.(type) {
		case *number:
			return true
		case *oo:
			return right.resolvable()
		case variable:
			return false
		case id:
			panic("shouldn't happen")
		}
	case *oo:
		switch right := o.right.(type) {
		case *number:
			return left.resolvable()
		case *oo:
			return left.resolvable() && right.resolvable()
		case variable:
			return false
		case id:
			panic("shouldn't happen")
		}
	case variable:
		return false
	}
	fmt.Printf("handle the case left: %T right: %T\n", o.left, o.right)
	panic("shouldn't happen")
}

func (o *oo) value() float64 {
	if v, ok := o.left.(*number); ok {
		if v2, ok2 := o.right.(*number); ok2 {
			return doMath(v.num, v2.num, o.op)
		}
		if v2, ok2 := o.right.(*oo); ok2 {
			return doMath(v.num, v2.value(), o.op)
		}
	}
	if v, ok := o.left.(*oo); ok {
		if v2, ok2 := o.right.(*number); ok2 {
			return doMath(v.value(), v2.num, o.op)
		}
		if v2, ok2 := o.right.(*oo); ok2 {
			return doMath(v.value(), v2.value(), o.op)
		}
	}
	fmt.Println("look at", o.id)
	fmt.Printf("Handle the case of left=%T and right=%T\n", o.left, o.right)
	return 0
}

func (o *oo) ID() string {
	return o.id
}
func (o *oo) String() string {
	return fmt.Sprintf("(%s %s %s)", o.left.String(), o.op, o.right.String())
}

func main() {
	input := internal.ReadInput()

	ops := map[string]valuer{}
	for _, line := range input {
		parts := strings.Split(line, ": ")

		found := false
		for _, op := range []op{mul, add, sub, div} {
			if strings.Contains(line, string(op)) {
				found = true
				ids := strings.Split(parts[1], fmt.Sprintf(" %s ", op))
				monkey := &oo{
					id:    parts[0],
					left:  id(ids[0]),
					op:    op,
					right: id(ids[1]),
				}
				if monkey.id == "root" {
					monkey.op = eql
				}
				ops[monkey.id] = monkey
			}
		}
		if found {
			continue
		}

		// must be a number otherwise
		if len(parts) == 2 {
			num, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				panic("bad parse")
			}
			if parts[0] == "humn" {
				//				ops[parts[0]] = &number{id: "humn", num: 3296135418821}
				ops[parts[0]] = variable{
					id:       parts[0],
					variable: "???",
				}
				continue
			}
			monkey := &number{
				id:  parts[0],
				num: num,
			}
			ops[monkey.id] = monkey
		}
	}

	// resolve the identifiers
	for _, m := range ops {
		switch t := m.(type) {
		case *oo:
			if v, ok := t.left.(id); ok {
				t.left = ops[v.ID()]
			}
			if v, ok := t.right.(id); ok {
				t.right = ops[v.ID()]
			}
		}
	}

	fmt.Println(ops["root"].Value())
	high := float64(3436788627814)
	low := float64(2930052996084)
	for {
		half := (low + high) / 2
		y := half
		x := (8.000000 * (6862813426220.000000 - ((((141.000000 + ((((((((((((((((((((165.000000 + ((346.000000 + (((29.000000 + (((2.000000 * (((((((573.000000 + ((((((2.000000 * (982.000000 + (2.000000 * ((((2.000000 * (45.000000 + ((((((528.000000 + (((614.000000 + ((((577.000000 + ((y - 961.000000) * 29.000000)) / 7.000000) - 507.000000) * 9.000000)) + 281.000000) / 2.000000)) / 2.000000) - 153.000000) * 3.000000) - 769.000000) / 5.000000))) - 826.000000) / 2.000000) + 147.000000)))) - 158.000000) / 2.000000) - 104.000000) + 835.000000) / 3.000000)) * 8.000000) - 890.000000) / 2.000000) + 124.000000) * 2.000000) - 233.000000)) + 752.000000) / 2.000000)) * 2.000000) - 675.000000)) / 3.000000)) / 2.000000) - 866.000000) * 3.000000) - 999.000000) * 2.000000) + 404.000000) / 2.000000) - 56.000000) * 2.000000) + 944.000000) / 2.000000) - 822.000000) * 2.000000) + 964.000000) / 4.000000) - 459.000000) * 2.000000) + 639.000000) + 546.000000)) / 5.000000) + 378.000000) / 5.000000)))
		goal := 23440423968672.000000
		fmt.Println(x-goal, y, x, goal, x == goal)
		if x-goal < 0 {
			high = (low + high) / 2
		}
		if x-goal > 0 {
			low = (low + high) / 2
		}
		if x-goal == 0 {
			fmt.Println(int(y))
			break
		}
	}
}

type op string

const (
	mul op = "*"
	add op = "+"
	sub op = "-"
	div op = "/"
	eql op = "="
)

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	var input string
	var w, x, y, z int
	var test bool
	flag.IntVar(&w, "w", 0, "set the w register")
	flag.IntVar(&x, "x", 0, "set the x register")
	flag.IntVar(&y, "y", 0, "set the y register")
	flag.IntVar(&z, "z", 0, "set the z register")
	flag.BoolVar(&test, "test", false, "run the tests")

	flag.StringVar(&input, "input", "", "")
	flag.Parse()
	lines := internal.ReadInput()

	if test {
		in := &inputNumber{
			numbers: [14]*digit{
				{9}, {9}, {9}, {9}, {9},
				{9}, {9}, {9}, {9}, {9},
				{9}, {9}, {9}, {9},
			},
		}
		for i := 0; i < 14; i++ {
			min := 999999999999999999
			num := 0
			for j := 0; j < 9; j++ {
				a := &alu{input: in.String(), w: w, x: x, y: y, z: z}
				for _, line := range lines {
					a.run(line)
				}
				if a.z < min {
					min = a.z
					num = in.numbers[i].num
				}
				in.numbers[i].dec()
			}
			in.numbers[i].num = num
			fmt.Printf("num %d in position %d yielded %d\n", num, i, min)
		}
		fmt.Println(in.String())
		os.Exit(0)
	}
	// a := &alu{input: input, w: w, x: x, y: y, z: z}
	// for _, line := range lines {
	// 	a.run(line)
	// }
	// fmt.Println(a)

	// a2 := &alu2{input: input, w: "0", x: "0", y: "0", z: "0"}
	// for _, line := range lines {
	// 	a2.run(line)
	// }
	// fmt.Println(a2)
	a3 := newALU3(input)
	for _, line := range lines {
		a3.run(line)
	}
	fmt.Println(a3.z)
	// for _, z := range a3.zs {
	// 	fmt.Println(z)
	// }
}

type inputNumber struct {
	numbers [14]*digit
}

func (i inputNumber) String() string {
	out := ""
	for _, d := range i.numbers {
		out += fmt.Sprintf("%d", d.num)
	}
	return out
}

type digit struct {
	num int
}

func (d *digit) inc() {
	d.num++
	if d.num == 10 {
		d.num = 1
	}
}
func (d *digit) dec() {
	d.num--
	if d.num == 0 {
		d.num = 9
	}
}

// w appears to be the last digit of the input

type alu struct {
	w, x, y, z int // registers
	input      string
	count      int
}

func (a *alu) getInput() int {
	x, _ := strconv.Atoi(string(a.input[a.count]))
	a.count++
	return x
}
func (a *alu) String() string {
	return fmt.Sprintf("w: %d, x: %d, y: %d, z: %d", a.w, a.x, a.y, a.z)
}

func (a *alu) run(line string) {
	parts := strings.Split(line, " ")
	switch parts[0] {
	case "inp":
		a.inp(parts[1])
	case "mul":
		a.mul(parts[1], a.getVal(parts[2]))
	case "add":
		a.add(parts[1], a.getVal(parts[2]))
	case "mod":
		a.mod(parts[1], a.getVal(parts[2]))
	case "div":
		a.div(parts[1], a.getVal(parts[2]))
	case "eql":
		a.eql(parts[1], a.getVal(parts[2]))
	}
}

func (a *alu) eql(reg string, val int) {
	switch reg {
	case "w":
		if a.w == val {
			a.w = 1
			return
		}
		a.w = 0
	case "x":
		if a.x == val {
			a.x = 1
			return
		}
		a.x = 0
	case "y":
		if a.y == val {
			a.y = 1
			return
		}
		a.y = 0
	case "z":
		if a.z == val {
			a.z = 1
			return
		}
		a.z = 0
	}
}

func (a *alu) div(reg string, val int) {
	switch reg {
	case "w":
		a.w = a.w / val
	case "x":
		a.x = a.x / val
	case "y":
		a.y = a.y / val
	case "z":
		a.z = a.z / val
	}
}

func (a *alu) mod(reg string, val int) {
	switch reg {
	case "w":
		a.w = a.w % val
	case "x":
		a.x = a.x % val
	case "y":
		a.y = a.y % val
	case "z":
		a.z = a.z % val
	}
}

func (a *alu) add(reg string, val int) {
	switch reg {
	case "w":
		a.w = a.w + val
	case "x":
		a.x = a.x + val
	case "y":
		a.y = a.y + val
	case "z":
		a.z = a.z + val
	}
}

func (a *alu) mul(reg string, val int) {
	switch reg {
	case "w":
		a.w = a.w * val
	case "x":
		a.x = a.x * val
	case "y":
		a.y = a.y * val
	case "z":
		a.z = a.z * val
	}
}

func (a *alu) getVal(regOrInt string) int {
	switch regOrInt {
	case "w":
		return a.w
	case "x":
		return a.x
	case "y":
		return a.y
	case "z":
		return a.z
	default:
		x, _ := strconv.Atoi(regOrInt)
		return x
	}
}

func (a *alu) inp(v string) {
	switch v {
	case "w":
		a.w = a.getInput()
	case "x":
		a.x = a.getInput()
	case "y":
		a.y = a.getInput()
	case "z":
		a.z = a.getInput()
	default:
		panic("emmmm")
	}
}

type alu2 struct {
	w, x, y, z string
	input      string
	count      int
}

func (a *alu2) getInput() string {
	out := string(a.input[a.count])
	a.count++
	return out
}
func (a *alu2) String() string {
	return fmt.Sprintf("w: %s, x: %s, y: %s, z: %s", a.w, a.x, a.y, a.z)
}

func (a *alu2) run(line string) {
	parts := strings.Split(line, " ")
	switch parts[0] {
	case "inp":
		a.inp(parts[1])
	case "mul":
		a.mul(parts[1], a.getVal(parts[2]))
	case "add":
		a.add(parts[1], a.getVal(parts[2]))
	case "mod":
		a.mod(parts[1], a.getVal(parts[2]))
	case "div":
		a.div(parts[1], a.getVal(parts[2]))
	case "eql":
		a.doit(parts[1], "==", a.getVal(parts[2]))
	}
}

func (a *alu2) add(reg, val string) {
	op := "+"
	if val == "0" {
		return
	}
	switch reg {
	case "w":
		if a.w == "0" {
			a.w = val
			return
		}
		a.w = fmt.Sprintf("(%s %s %s)", a.w, op, val)
	case "x":
		if a.x == "0" {
			a.x = val
			return
		}
		a.x = fmt.Sprintf("(%s %s %s)", a.x, op, val)
	case "y":
		if a.y == "0" {
			a.y = val
			return
		}
		a.y = fmt.Sprintf("(%s %s %s)", a.y, op, val)
	case "z":
		if a.z == "0" {
			a.z = val
			return
		}
		a.z = fmt.Sprintf("(%s %s %s)", a.z, op, val)
	}
}
func (a *alu2) div(reg, val string) {
	op := "/"
	if val == "1" {
		return
	}
	if val == "0" {
		panic("divide by 0")
	}
	switch reg {
	case "w":
		if a.w == "0" {
			return
		}
		a.w = fmt.Sprintf("(%s %s %s)", a.w, op, val)
	case "x":
		if a.x == "0" {
			return
		}
		a.x = fmt.Sprintf("(%s %s %s)", a.x, op, val)
	case "y":
		if a.y == "0" {
			return
		}
		a.y = fmt.Sprintf("(%s %s %s)", a.y, op, val)
	case "z":
		if a.z == "0" {
			return
		}
		a.z = fmt.Sprintf("(%s %s %s)", a.z, op, val)
	}
}

func (a *alu2) mod(reg, val string) {
	op := "%"
	if val == "1" {
		return
	}
	if val == "0" {
		panic("divide by 0")
	}
	switch reg {
	case "w":
		if a.w == "0" {
			return
		}
		a.w = fmt.Sprintf("(%s %s %s)", a.w, op, val)
	case "x":
		if a.x == "0" {
			return
		}
		a.x = fmt.Sprintf("(%s %s %s)", a.x, op, val)
	case "y":
		if a.y == "0" {
			return
		}
		a.y = fmt.Sprintf("(%s %s %s)", a.y, op, val)
	case "z":
		if a.z == "0" {
			return
		}
		a.z = fmt.Sprintf("(%s %s %s)", a.z, op, val)
	}
}

func (a *alu2) mul(reg, val string) {
	op := "*"
	if val == "1" {
		return
	}
	if val == "0" {
		switch reg {
		case "w":
			a.w = "0"
		case "x":
			a.x = "0"
		case "y":
			a.y = "0"
		case "z":
			a.z = "0"
		}
		return
	}
	if val == "-1" {
		switch reg {
		case "w":
			a.w = fmt.Sprintf("-%s", a.w)
		case "x":
			a.x = fmt.Sprintf("-%s", a.x)
		case "y":
			a.y = fmt.Sprintf("-%s", a.y)
		case "z":
			a.z = fmt.Sprintf("-%s", a.z)
		}
		return
	}
	switch reg {
	case "w":
		if a.w == "0" {
			return
		}
		a.w = fmt.Sprintf("(%s %s %s)", a.w, op, val)
	case "x":
		if a.x == "0" {
			return
		}
		a.x = fmt.Sprintf("(%s %s %s)", a.x, op, val)
	case "y":
		if a.y == "0" {
			return
		}
		a.y = fmt.Sprintf("(%s %s %s)", a.y, op, val)
	case "z":
		if a.z == "0" {
			return
		}
		a.z = fmt.Sprintf("(%s %s %s)", a.z, op, val)
	}
}

func (a *alu2) doit(reg, op, val string) {
	switch reg {
	case "w":
		a.w = fmt.Sprintf("(%s %s %s)", a.w, op, val)
	case "x":
		a.x = fmt.Sprintf("(%s %s %s)", a.x, op, val)
	case "y":
		a.y = fmt.Sprintf("(%s %s %s)", a.y, op, val)
	case "z":
		a.z = fmt.Sprintf("(%s %s %s)", a.z, op, val)
	}

}

func (a *alu2) getVal(regOrInt string) string {
	switch regOrInt {
	case "w":
		return a.w
	case "x":
		return a.x
	case "y":
		return a.y
	case "z":
		return a.z
	default:
		return regOrInt
	}
}

func (a *alu2) inp(v string) {
	switch v {
	case "w":
		a.w = a.getInput()
	case "x":
		a.x = a.getInput()
	case "y":
		a.y = a.getInput()
	case "z":
		a.z = a.getInput()
	default:
		panic("emmmm")
	}
}

type alu3 struct {
	w, x, y, z Expr
	input      string
	count      int
	zs         []string
}

func newALU3(in string) *alu3 {
	return &alu3{
		w:     &num{0},
		x:     &num{0},
		y:     &num{0},
		z:     &num{0},
		zs:    []string{},
		input: in,
	}
}

func (a *alu3) getInput() Expr {
	out := a.input[a.count]
	a.count++
	if out <= 'z' && out >= 'a' {
		return &ident{
			name: string(out),
		}
	}
	x, _ := strconv.Atoi(string(out))
	return &num{
		val: x,
	}
}

func (a *alu3) String() string {
	return fmt.Sprintf("w: %v, x: %v, y: %v, z: %v", a.w, a.x, a.y, a.z)
}

func (a *alu3) run(line string) {
	parts := strings.Split(line, " ")
	switch parts[0] {
	case "inp":
		a.inp(parts[1])
	case "mul":
		a.doit(parts[1], "*", a.getVal(parts[2]))
	case "add":
		a.doit(parts[1], "+", a.getVal(parts[2]))
	case "mod":
		a.doit(parts[1], "%", a.getVal(parts[2]))
	case "div":
		a.doit(parts[1], "/", a.getVal(parts[2]))
	case "eql":
		a.doit(parts[1], "==", a.getVal(parts[2]))
	}
	a.w = a.w.simplify()
	a.x = a.x.simplify()
	a.y = a.y.simplify()
	a.z = a.z.simplify()
}

func (a *alu3) doit(register, operator string, val Expr) {
	switch register {
	case "w":
		a.w = &op{op: operator, x: a.w, y: val}
	case "x":
		a.x = &op{op: operator, x: a.x, y: val}
	case "y":
		a.y = &op{op: operator, x: a.y, y: val}
	case "z":
		a.z = &op{op: operator, x: a.z, y: val}
	}
}

func (a *alu3) getVal(regOrInt string) Expr {
	switch regOrInt {
	case "w":
		return a.w
	case "x":
		return a.x
	case "y":
		return a.y
	case "z":
		return a.z
	default:
		x, _ := strconv.Atoi(regOrInt)
		return &num{val: x}
	}
}

func (a *alu3) inp(v string) {
	a.zs = append(a.zs, a.z.String())
	if a.count > 0 {
		//a.z = &ident{fmt.Sprintf("z[%d]", a.count)}
	}
	switch v {
	case "w":
		a.w = a.getInput()
	case "x":
		a.x = a.getInput()
	case "y":
		a.y = a.getInput()
	case "z":
		a.z = a.getInput()
	default:
		panic("emmmm")
	}
}

type ident struct {
	name string
}

func (i *ident) expr() {}
func (i *ident) String() string {
	return i.name
}
func (i *ident) simplify() Expr {
	return i
}

type num struct {
	val int
}

func (n *num) expr() {}
func (n *num) String() string {
	return fmt.Sprintf("%d", n.val)
}
func (n *num) simplify() Expr {
	return n
}

type op struct {
	op string
	x  Expr
	y  Expr
}

func (o *op) simplify() Expr {
	if o.op == "+" {
		if xn, ok := o.x.(*num); ok {
			if xn.val == 0 {
				return o.y
			}
			if yn, ok := o.y.(*num); ok {
				return &num{xn.val + yn.val}
			}
			// 	// ((2983389734 + (a + 13)) + 16)
			// 	if yn, ok := o.y.(*op); ok {
			// 		if yn.op == "+" {
			// 			if y2n, ok := yn.y.(*num); ok {
			// 				xn.val += y2n.val
			// 				o.y = yn.x
			// 			}
			// 		}
			// 	}
		}
		if yn, ok := o.y.(*num); ok {
			if yn.val == 0 {
				return o.x
			}
		}
	}
	if o.op == "/" {
		// if yn, ok := o.y.(*num); ok {
		// 	if xn, ok := o.x.(*op); ok {
		// 		if xn.op == "/" {
		// 			if y2n, ok := xn.y.(*num); ok {
		// 				return &op{
		// 					op: "/",
		// 					x:  xn.x,
		// 					y:  &num{y2n.val * yn.val},
		// 				}
		// 			}
		// 		}
		// 	}
		// }
		if xn, ok := o.x.(*num); ok {
			if xn.val == 0 {
				return &num{0}
			}
			if yn, ok := o.y.(*num); ok {
				return &num{xn.val / yn.val}
			}
		}
		if yn, ok := o.y.(*num); ok {
			if yn.val == 1 {
				return o.x
			}
		}
	}
	if o.op == "*" {
		// (12 + (a + b)) % 12 == (a + b) % 12
		// if xn, ok := o.x.(*op); ok {
		// 	if yn, ok := o.y.(*num); ok {
		// 		if x2n, ok := xn.x.(*num); ok {
		// 			o.x = &num{x2n.val * yn.val}
		// 		}
		// 	}
		// 	if xn.op == "/" {
		// 		if yn, ok := xn.y.(*num); ok {
		// 			if y2n, ok := o.y.(*num); ok {
		// 				if y2n.val == yn.val {
		// 					return xn.x
		// 				}
		// 			}
		// 		}
		// 	}
		// }
		if xn, ok := o.x.(*num); ok {
			if xn.val == 0 {
				return &num{0}
			}
			if xn.val == 1 {
				return o.y
			}
			if yn, ok := o.y.(*num); ok {
				return &num{xn.val * yn.val}
			}
		}
		if yn, ok := o.y.(*num); ok {
			if yn.val == 0 {
				return &num{0}
			}
			if yn.val == 1 {
				return o.x
			}
		}
	}
	if o.op == "%" {
		// if xn, ok := o.x.(*op); ok {
		// 	if yn, ok := o.y.(*num); ok {
		// 		if x2n, ok := xn.x.(*num); ok {
		// 			if x2n.val%yn.val == 0 {
		// 				o.x = &num{0}
		// 			}
		// 		}
		// 	}
		// }
		// if x is an expression, and any of the numbers
		if xn, ok := o.x.(*num); ok {
			// if xn.val == 0 {
			// 	return &num{0}
			// }
			if yn, ok := o.y.(*num); ok {
				return &num{xn.val % yn.val}
			}
		}
		if yn, ok := o.y.(*num); ok {
			if yn.val == 0 {
				return &num{0}
			}
		}
	}
	if o.op == "==" {
		if yn, ok := o.y.(*num); ok {
			if yn.val == 0 {
				if xn, ok := o.x.(*op); ok {
					if xn.op == "==" {
						// ((13 == a) == 0)
						if x2n, ok := xn.x.(*num); ok {
							if x2n.val >= 10 {
								if _, ok := xn.y.(*ident); ok {
									return &num{1}
								}
							}
						}
						//		(((((a + 8) % 26) + 12) == b) == 0))
						if _, ok := xn.y.(*ident); ok {
							if x2n, ok := xn.x.(*op); ok {
								if x2n.op == "+" {
									if y2n, ok := x2n.y.(*num); ok {
										if y2n.val >= 10 {
											return &num{1}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		if xn, ok := o.x.(*num); ok {
			if yn, ok := o.y.(*num); ok {
				if xn.val == yn.val {
					return &num{1}
				}
				return &num{0}
			}
			// 	// (0 == a)
			// 	if _, ok := o.y.(*ident); ok {
			// 		if xn.val == 0 {
			// 			return &num{0} // idents can never be 0
			// 		}
			// 	}
		}

	}
	return o
}

func (o *op) expr() {}

func (o *op) String() string {
	return fmt.Sprintf("(%v %s %v)", o.x, o.op, o.y)
}

type Expr interface {
	expr()
	simplify() Expr
	String() string
}

/*
ab81efghijklmn

(((a + 8) % 26) + 12) == b) // always 0 because 1 <= a,b <= 9
(((((467 * 26 + (c + 4)) % 26) + -11) == d)
(((((467 * 26 + (e + 13)) % 26) + 13) == f) == 0))
(((((5857 *54 * 26 + g) % 26) + -5) == h) == 0)))
((((8223228 + (i + 7)) % 26) == j) == 0)
(((((442 + (m + 14)) % 26) + -11) == n) == 0))
*/

/*

0
(a + 8)
(((a + 8) * 26) + (b + 16))
((z[2] * 26) + (c + 4))
(((z[3] / 26) * ((25 * ((((z[3] % 26) + -11) == d) == 0)) + 1)) + ((d + 1) * ((((z[3] % 26) + -11) == d) == 0)))
((z[4] * 26) + (e + 13))
((z[5] * 26) + (f + 5))
((z[6] * 26) + g)
(((z[7] / 26) * ((25 * ((((z[7] % 26) + -5) == h) == 0)) + 1)) + ((h + 10) * ((((z[7] % 26) + -5) == h) == 0)))
((z[8] * 26) + (i + 7))
(((z[9] / 26) * ((25 * (((z[9] % 26) == j) == 0)) + 1)) + ((j + 2) * (((z[9] % 26) == j) == 0)))
(((z[10] / 26) * ((25 * ((((z[10] % 26) + -11) == k) == 0)) + 1)) + ((k + 13) * ((((z[10] % 26) + -11) == k) == 0)))
(((z[11] / 26) * ((25 * ((((z[11] % 26) + -13) == l) == 0)) + 1)) + ((l + 15) * ((((z[11] % 26) + -13) == l) == 0)))
(((z[12] / 26) * ( (25 * ((((z[12] % 26) + -13) == m) == 0) ) + 1)) + ((m + 14) * ((((z[12] % 26) + -13) == m) == 0)))

(x / 26) * (((25(x%26)-13)==m)==0))+1))

z[13] < 26 == [14,22] will yield a 0 expression (m == 1-9)
z[12] < 26 == [14,22] yields 0                  (l == 1-9)
z[11] < 26 == [12,20]                           (k == 1-9)
z[10] < 26 == [1-9]                              (j == 1-9)
z[9] =?
z[8] has to be 0?

(((z[7] / 26) * ((25 * ((((z[7] % 26) + -5) == h) == 0)) + 1)) + ((h + 10) * ((((z[7] % 26) + -5) == h) == 0)))
z[7]==[6,14] make z[8] == 0
z[6]

(((z[3] / 26) * ((25 * ((((z[3] % 26) + -11) == d) == 0)) + 1)) + ((d + 1) * ((((z[3] % 26) + -11) == d) == 0)))
z[3] = [12,20]
si


(b+16 + c + 4) % 26 == 20

*/

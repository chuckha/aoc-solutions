package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

/*
162012 too high
116274 too low
123149 correct
*/

// horrific code.

var eof = errors.New("EOF")

const boxWidth = 50
const faceSize = 50

func main() {
	lines := internal.ReadRealRawInput()
	overview := lines[:len(lines)-2]
	longest := math.MinInt
	for _, line := range overview {
		if len(line) > longest {
			longest = len(line)
		}
	}
	grid := internal.NewGridV3[*cell]()
	cube := newCube()
	keys := []int{}
	for y, line := range overview {
		for x, c := range line {
			if c == ' ' {
				continue
			}
			fid := faceID(x, y, longest/faceSize)
			if _, ok := cube.faces[fid]; !ok {
				cube.faces[fid] = newFace(fid)
				keys = append(keys, fid)
			}
			pt := internal.Point{x, y}
			relativePoint := internal.Point{x - (x/faceSize)*faceSize, y - (y/faceSize)*faceSize}
			cube.faces[fid].data.Set(relativePoint, &cell{point: relativePoint, data: string(c)})
			grid.Set(pt, &cell{point: pt, data: string(c)})
		}
	}

	// sort.Ints(keys)
	// for _, k := range keys {
	// 	fmt.Println(cube.faces[k].data)
	// }
	grid.DefaultOutput = " "
	fmt.Println(grid)
	instr := lines[len(lines)-2]
	//	instr = "0R50R5L50R1L10R1L10R1R20"
	dir := &directions{ptr: 0, raw: instr}
	starting := findStartingPoint(grid)
	person := newPerson(starting)

	//	for i := 0; i < 200; i++ {
	for {
		act, err := dir.next()
		person.move(grid, act)
		// fmt.Println(act)
		// fmt.Println(grid.Layer(person.path))
		// time.Sleep(1 * time.Second)

		if err != nil {
			break
		}
	}
	person.path.Set(person.location, &cell{person.location, string(internal.Red(string(person.facing)))})
	fmt.Println(grid.Layer(person.path))
	fmt.Println(person.password())

}

func faceID(x, y, numPerRow int) int {
	return (y/faceSize)*numPerRow + (x / faceSize)
}

type cell struct {
	point internal.Point
	data  string
}

func (c *cell) String() string {
	return c.data
}

type directions struct {
	ptr int
	raw string
}

func (d *directions) next() ([]action, error) {
	if d.ptr == len(d.raw) {
		return nil, eof
	}
	out := []action{}
	x := d.peek()
	switch {
	case x >= 48 && x <= 57:
		forward, err := d.readInt()
		for i := 0; i < forward; i++ {
			out = append(out, fwd)
		}
		if err == eof {
			return out, eof
		}
	default:
		turn := d.readOne()
		out = append(out, actionFactory(turn))
	}
	return out, nil
}

func (d *directions) readInt() (int, error) {
	item := []byte{}
	cur := d.raw[d.ptr]
	end := false
	for cur >= 48 && cur <= 57 {
		item = append(item, cur)
		d.ptr++
		if d.ptr == len(d.raw) {
			end = true
			break
		}
		cur = d.raw[d.ptr]
	}
	i, err := strconv.Atoi(string(item))
	if err != nil {
		return -1, err
	}
	if end {
		return i, eof
	}
	return i, nil
}

func (d *directions) readOne() string {
	out := d.raw[d.ptr]
	d.ptr++

	return string(out)
}

func (d *directions) backup() {
	d.ptr--
}

func (d *directions) peek() byte {
	return d.raw[d.ptr]
}

type person struct {
	location internal.Point
	facing   dir
	face     int
	path     *internal.GridV3[*cell]
}

func newPerson(init internal.Point) *person {
	path := internal.NewGridV3[*cell]()
	person := &person{
		location: init,
		facing:   right,
		path:     path,
	}
	path.Set(init, &cell{init, string(person.facing)})
	return person
}

func (p *person) password() int {
	vals := map[dir]int{right: 0, down: 1, left: 2, up: 3}
	return 1000*(p.location.Y+1) + (p.location.X+1)*4 + vals[p.facing]
}

func (p *person) move(board *internal.GridV3[*cell], actions []action) {
	for _, action := range actions {
		if action == fwd {
			startbox := box(p.location)
			// walk forward
			p.forward()

			// Box 1 and 2
			if p.location.Y < board.Min.Y {
				if startbox == 1 {
					// end up in box 6
					pt := internal.Point{0, p.location.X + 100}
					if board.At(pt).data == "#" {
						p.backup()
						continue
					}
					p.location = pt
					p.facing = p.facing.turn(cw)
				}
				if startbox == 2 {
					// end up in box 6
					// no direction change
					pt := internal.Point{p.location.X - 100, 199}
					if board.At(pt).data == "#" {
						p.backup()
						continue
					}
					p.location = pt
				}
			}

			// Box 6
			if p.location.Y > board.Max.Y {
				pt := internal.Point{p.location.X + 100, 0}
				if board.At(pt).data == "#" {
					p.backup()
					continue
				}
				p.location = pt
			}

			// Box 4, 6
			if p.location.X < board.Min.X {
				if startbox == 4 {
					pt := internal.Point{50, 149 - p.location.Y}
					if board.At(pt).data == "#" {
						p.backup()
						continue
					}
					// do 180
					p.facing = p.facing.turn(cw)
					p.facing = p.facing.turn(cw)
					p.location = pt
				}
				if startbox == 6 {
					pt := internal.Point{p.location.Y - 100, 0}
					if board.At(pt).data == "#" {
						p.backup()
						continue
					}
					p.facing = p.facing.turn(ccw)
					p.location = pt
				}
			}

			// Box 2
			if p.location.X > board.Max.X {
				pt := internal.Point{99, 149 - p.location.Y}
				if board.At(pt).data == "#" {
					p.backup()
					continue
				}
				// do a 180
				p.facing = p.facing.turn(cw)
				p.facing = p.facing.turn(cw)
				p.location = pt
			}

			// now we do inbound walk off check
			if !board.In(p.location) || board.At(p.location).data == " " {
				switch p.facing {
				case up:
					// Box 4
					pt := internal.Point{50, p.location.X + 50}
					if board.At(pt).data == "#" {
						p.backup()
						continue
					}
					p.facing = p.facing.turn(cw)
					p.location = pt
				case right:
					// Box 3, 5, 6
					if startbox == 3 {
						pt := internal.Point{p.location.Y + 50, 49}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(ccw)
						p.location = pt
					}
					if startbox == 5 {
						pt := internal.Point{149, 149 - p.location.Y}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(cw)
						p.facing = p.facing.turn(cw)
						p.location = pt
					}
					if startbox == 6 {
						pt := internal.Point{p.location.Y - 100, 149}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(ccw)
						p.location = pt
					}
				case down:
					// Box 2, 5
					if startbox == 2 {
						pt := internal.Point{99, p.location.X - 50}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(cw)
						p.location = pt
					}
					if startbox == 5 {
						pt := internal.Point{49, p.location.X + 100}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(cw)
						p.location = pt
					}
				case left:
					// Box 1, 3
					if startbox == 1 {
						pt := internal.Point{0, 149 - p.location.Y}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(cw)
						p.facing = p.facing.turn(cw)
						p.location = pt
					}
					if startbox == 3 {
						pt := internal.Point{p.location.Y - 50, 100}
						if board.At(pt).data == "#" {
							p.backup()
							continue
						}
						p.facing = p.facing.turn(ccw)
						p.location = pt
					}
				}
			}
			if board.At(p.location).data == "#" {
				p.backup()
				continue
			}
			p.path.Set(p.location, &cell{p.location, string(internal.Red(string(p.facing)))})
			continue
		}
		p.facing = p.facing.turn(action)
		p.path.Set(p.location, &cell{p.location, string(internal.Red(string(p.facing)))})
	}
}

// find out which box i'm in
// build a box map, box 1

func findBottomOpposite(b *internal.GridV3[*cell], off internal.Point) internal.Point {
	for y := b.Max.Y; y > off.Y; y-- {
		pt := internal.Point{off.X, y}
		if b.In(pt) {
			if b.At(pt).data == " " {
				continue
			}
			if b.At(pt).data == "." {
				return pt
			}
			if b.At(pt).data == "#" {
				return pt
			}
		}
	}
	panic("couldn't find bottom opposite")
}

func findTopOpposite(b *internal.GridV3[*cell], off internal.Point) internal.Point {
	for y := b.Min.Y; y < off.Y; y++ {
		pt := internal.Point{off.X, y}
		if b.In(pt) {
			if b.At(pt).data == " " {
				continue
			}
			if b.At(pt).data == "." {
				return pt
			}
			if b.At(pt).data == "#" {
				return pt
			}
		}
	}
	panic("couldn't find top opposite")
}

func findRightOpposite(b *internal.GridV3[*cell], off internal.Point) internal.Point {
	for x := b.Max.X; x > off.X; x-- {
		pt := internal.Point{x, off.Y}
		if b.In(pt) {
			if b.At(pt).data == " " {
				continue
			}
			if b.At(pt).data == "." {
				return pt
			}
			if b.At(pt).data == "#" {
				return pt
			}
		}
	}
	panic("couldn't find right opposite")
}

func findLeftOpposite(b *internal.GridV3[*cell], off internal.Point) internal.Point {
	for x := b.Min.X; x < off.X; x++ {
		pt := internal.Point{x, off.Y}
		if b.In(pt) {
			if b.At(pt).data == " " {
				continue
			}
			if b.At(pt).data == "." {
				return pt
			}
			if b.At(pt).data == "#" {
				return pt
			}
		}
	}
	panic("couldn't find right opposite")
}

func (p *person) forward() {
	switch p.facing {
	case up:
		p.location = p.location.Up()
	case right:
		p.location = p.location.Right()
	case down:
		p.location = p.location.Down()
	case left:
		p.location = p.location.Left()
	}
}

func (p *person) backup() {
	switch p.facing {
	case up:
		p.location = p.location.Down()
	case right:
		p.location = p.location.Left()
	case down:
		p.location = p.location.Up()
	case left:
		p.location = p.location.Right()
	}
}

type dir string

const (
	up    dir = "^"
	right dir = ">"
	down  dir = "v"
	left  dir = "<"
)

var rightDirs = []dir{up, right, down, left}
var leftDirs = []dir{up, left, down, right}

func (d dir) turn(a action) dir {
	var dirs []dir
	switch a {
	case cw:
		dirs = rightDirs
	case ccw:
		dirs = leftDirs
	}
	idx := internal.Search(d, dirs)
	return dirs[(idx+1)%len(dirs)]
}

type action string

const (
	fwd action = "forward"
	cw  action = "clockwise"
	ccw action = "counter-clockwise"
)

func actionFactory(s string) action {
	switch s {
	case "L":
		return ccw
	case "R":
		return cw
	default:
		panic("bad action")
	}
}

func findStartingPoint(g *internal.GridV3[*cell]) internal.Point {
	for i := g.Min.X; i < g.Max.X; i++ {
		if !g.In(internal.Point{i, g.Min.Y}) {
			continue
		}
		if g.At(internal.Point{i, g.Min.Y}).data == "." {
			return internal.Point{i, g.Min.Y}
		}
	}
	panic("no starting point?")
}

func box(pt internal.Point) int {
	if pt.X >= 50 && pt.X <= 99 && pt.Y >= 0 && pt.Y <= 49 {
		return 1
	}
	if pt.X >= 100 && pt.X <= 149 && pt.Y >= 0 && pt.Y <= 49 {
		return 2
	}
	if pt.X >= 50 && pt.X <= 99 && pt.Y >= 50 && pt.Y <= 99 {
		return 3
	}
	if pt.X >= 0 && pt.X <= 49 && pt.Y >= 100 && pt.Y <= 149 {
		return 4
	}
	if pt.X >= 50 && pt.X <= 99 && pt.Y >= 100 && pt.Y <= 149 {
		return 5
	}
	if pt.X >= 0 && pt.X <= 49 && pt.Y >= 150 && pt.Y <= 199 {
		return 6
	}
	return -1
}

// idea to actually solve it
// build out a list of faces instead of hard coding the values for our one input
type face struct {
	id   int
	data *internal.GridV3[*cell]
}

func newFace(fid int) *face {
	return &face{
		id:   fid,
		data: internal.NewGridV3[*cell](),
	}
}

type cube struct {
	person *person
	faces  map[int]*face
}

func newCube() *cube {
	return &cube{
		person: newPerson(internal.Point{}),
		faces:  map[int]*face{},
	}
}
func (c *cube) addFace(f *face) {
	c.faces[f.id] = f
}

// walk along a cube

/*

manual input
if it's value is " " it's not in a box

0-49
50-99
100-149
150-199

box1
min: {50,0}
max: {99,49}
box2:
min {100,0}
max {149,49}
box3:
min: {50,50}
max: {99,99}
box4:
min: {0,100}
max: {49,149}
box5:
min: {50,100}
max: {99,149}
box6:
min: {0,200}
max: {49,249}

   6
 4 1 2 5
   3

1 -> 2, no direction change
1 -> 3, no direction change
1 -> 4, do a 180 ((go left) {50, 0-49} -> {0, 149-y})
1 -> 6, cw turn ((go up) {50, 0} -> {0, 150} / {51,0} -> {0, 151} / ... / {99, 0} -> {0,199})

   6
 1 2 5 4
   3

2 -> 5, do a 180 (go right {149, 0} -> {99, 149} / {149, 49} -> {99,100} / {149, 1} -> {99, 148}) (149 - y)
2 -> 3, cw turn (go down {100, 49} -> {99, 50} / {149,49} -> {99, 99}) (x = 99, y = x-50)
2 -> 1, no direction change
2 -> 6, no direction change ((go up) {100, 0} -> {0,199} / {149,0}->{49,199})

  1
4 3 2 6
  5

3 -> 2, ccw turn ((go right) {99, 50} -> {100, 49} / {99, 51} -> {101, 49}) (y + 50, 49)
3 -> 5, no direction change
3 -> 4, cw turn ((go left), {49, 50-99} -> {y-50,100})
3 -> 1, no direction change

  3
1 4 5 2
  6

4 -> 5, no direction change
4 -> 6, no direction change
4 -> 1, do a 180 ((go left) {0, 100} -> {50,49} / {0, 101} -> {50, 48} / {0, 149} -> {50, 0}) (x is always 50, y == 149-y)
4 -> 3, do a cw turn ((go up) {49, 100} -> {50, 99} / {0,100} -> {50, 50} ) x always 50, y = x + 50

  3
4 5 2 1
  6

5 -> 2, do a 180 ((go right) {99, 100} -> {149, 49} / {99, 149} -> {149, 0}) (x always 149, y = 149-y) (149 -> 0, 148 -> 1, 149 - (y))
5 -> 6, cw turn ((go down) {50-99, 149} -> {49, x+100})
5 -> 4, no direction change
5 -> 3, no direction change

  4
1 6 5 3
  2

6 -> 5, ccw turn ((go right) {49, 150} -> {50, 149})
6 -> 2, no direction change (go down {0,199} -> {100,0}, {49, 199} -> {149,0})
6 -> 1, ccw turn ((go left) {0,150} -> {50, 0} / {0, 151} -> {51, 0}, {0, 199} -> {99,0}})
6 -> 4, no direction change

*/

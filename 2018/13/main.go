package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

type direction string
type track string
type turn string

var turns = []turn{left, straight, right}

const (
	north direction = "north"
	east  direction = "east"
	south direction = "south"
	west  direction = "west"

	northSouth   track = "|"
	eastWest     track = "-"
	intersection track = "+"
	northToEast  track = "/"
	northToWest  track = `\`

	left     turn = "left"
	straight turn = "straight"
	right    turn = "right"
)

func main() {
	lines := internal.ReadRawInput()
	g := newGrid()
	for j, line := range lines {
		for i, c := range line {
			if string(c) == "" {
				continue
			}
			g.add(i, j, string(c))
		}
	}
	// part 1
	// for !g.hasCollision() {
	// 	g.moveCarts()
	// 	fmt.Println(g)
	// 	fmt.Println(strings.Repeat("$", 150))
	// }
	// for _, c := range g.carts {
	// 	if c.collided {
	// 		fmt.Println(c.pos)
	// 		break
	// 	}
	// }
	//part 2
	itrs := 0
	for g.cartCount() > 1 {
		g.moveCarts()
		// fmt.Println(strings.Repeat("$", 9), g.cartCount())
		// fmt.Println(g)
		fmt.Println(g.carts)
		itrs++
	}
	fmt.Println(itrs)
	for _, c := range g.carts {
		if c.collided {
			continue
		}
		fmt.Println(c.pos)
	}
}

type grid struct {
	data       map[point]track
	minx, maxx int
	miny, maxy int
	carts      []*Cart
	collisions []point
}

func (g *grid) hasCollision() bool {
	for _, c := range g.carts {
		if c.collided {
			return true
		}
	}
	return false
}
func (g *grid) cartCount() int {
	count := 0
	for _, c := range g.carts {
		if !c.collided {
			count++
		}
	}
	return count
}

func (g *grid) moveCarts() {
	// check if the next move of each cart will collide with another cart
	// then check if the next moves

	for _, cart := range g.carts {
		if cart.collided {
			continue
		}
		nextPos := cart.pos.next(cart.direction)
		cart.pos = nextPos

		for _, c := range g.carts {
			if c.collided {
				continue
			}
			if c == cart {
				continue
			}
			if c.pos.x == cart.pos.x && c.pos.y == cart.pos.y {
				// update the current cart pos to where the collision took place
				// the other cart doesn't exist now as a collision is 2 carts
				// remove both from g.carts
				cart.collided = true
				c.collided = true
				fmt.Printf("%d carts remain\n", g.cartCount())
				break
			}
		}

		// update cart direction
		switch g.data[nextPos] {
		case eastWest:
			if cart.direction == north || cart.direction == south {
				panic("weird track")
			}
		case northSouth:
			if cart.direction == east || cart.direction == west {
				panic("very weird track")
			}
		case intersection:
			cart.turn()
		case northToEast: // '/
			switch cart.direction {
			case north:
				cart.direction = east
			case west:
				cart.direction = south
			case east:
				cart.direction = north
			case south:
				cart.direction = west
			}
		case northToWest:
			switch cart.direction {
			case north:
				cart.direction = west
			case west:
				cart.direction = north
			case east:
				cart.direction = south
			case south:
				cart.direction = east
			}
		}
	}
	sort.Sort(carts(g.carts))
	for _, c := range g.carts {
		for _, c2 := range g.carts {
			if c2.collided || c.collided {
				continue
			}
			if c == c2 {
				continue
			}
			if c.pos.x == c2.pos.x && c.pos.y == c2.pos.y {
				c.collided = true
				c2.collided = true
			}
		}
	}
}

type carts []*Cart

func (c carts) Len() int      { return len(c) }
func (c carts) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c carts) Less(i, j int) bool {
	if c[i].pos.y < c[j].pos.y {
		return true
	}
	if c[i].pos.y > c[j].pos.y {
		return false
	}
	return c[i].pos.x < c[j].pos.x
}

type Cart struct {
	pos       point
	direction direction
	collided  bool
	lastTurn  turn
}

func newCart(i, j int, d string) *Cart {
	return &Cart{
		pos:       point{i, j},
		direction: newCartDirection(d),
		lastTurn:  right, // init so the first turn is left
	}
}

func (c *Cart) turn() {
	switch c.lastTurn {
	case left:
		c.lastTurn = straight
		// direction unaffected
	case straight:
		c.lastTurn = right
	case right:
		c.lastTurn = left
	}
	c.direction = c.direction.turn(c.lastTurn)
}

func (c *Cart) String() string {
	if c.collided {
		return "X"
	}
	switch c.direction {
	case north:
		return "^"
	case south:
		return "v"
	case east:
		return ">"
	case west:
		return "<"
	default:
		panic("invalid direction for cart")
	}
}

func newGrid() *grid {
	return &grid{
		data: make(map[point]track),
		minx: 0, maxx: 0, miny: 0, maxy: 0,
	}
}
func (g *grid) add(i, j int, d string) {
	if d == "" || d == " " {
		return
	}
	if i < g.minx {
		g.minx = i
	}
	if i > g.maxx {
		g.maxx = i
	}
	if j < g.miny {
		g.miny = j
	}
	if j > g.maxy {
		g.maxy = j
	}
	track := tt(d)
	g.data[point{i, j}] = track
	switch d {
	case "^", ">", "v", "<":
		g.carts = append(g.carts, newCart(i, j, d))
	}
}

func tt(d string) track {
	switch d {
	case ">":
		return eastWest
	case "<":
		return eastWest
	case "^":
		return northSouth
	case "v":
		return northSouth
	case "|":
		return northSouth
	case "-":
		return eastWest
	case "+":
		return intersection
	case `\`:
		return northToWest
	case "/":
		return northToEast
	default:
		panic("ahhhhhhh")
	}
}

func (g *grid) String() string {
	var out strings.Builder
	for j := g.miny; j <= g.maxy; j++ {
		for i := g.minx; i <= g.maxx; i++ {
			item, ok := g.data[point{i, j}]
			if !ok {
				out.WriteString(" ")
				continue
			}
			printTrack := true
			for _, c := range g.carts {
				if c.pos.x == i && c.pos.y == j {
					if c.collided {
						continue
					}
					printTrack = false
					out.WriteString(string(c.String()))
					break
				}
			}
			if printTrack {
				out.WriteString(string(item))
			}
		}
		out.WriteString("\n")
	}
	return out.String()

}

type point struct {
	x, y int
}

func (p point) next(direction direction) point {
	switch direction {
	case north:
		return p.north()
	case east:
		return p.east()
	case south:
		return p.south()
	case west:
		return p.west()
	default:
		panic("invalid direction for point")
	}
}
func (p point) north() point {
	return point{p.x, p.y - 1}
}
func (p point) east() point {
	return point{p.x + 1, p.y}
}
func (p point) south() point {
	return point{p.x, p.y + 1}
}
func (p point) west() point {
	return point{p.x - 1, p.y}
}

func (d direction) turn(turn turn) direction {
	switch turn {
	case left:
		switch d {
		case north:
			return west
		case east:
			return north
		case south:
			return east
		case west:
			return south
		}
	case straight:
		return d
	case right:
		switch d {
		case north:
			return east
		case east:
			return south
		case south:
			return west
		case west:
			return north
		}
	}
	panic(fmt.Sprintf("unreachable: %s", turn)) // why does this need to be here...
}
func newCartDirection(s string) direction {
	switch s {
	case "^":
		return north
	case ">":
		return east
	case "v":
		return south
	case "<":
		return west
	}
	panic("unknown direction:" + s)
}

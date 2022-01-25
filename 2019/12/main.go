package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
	"github.com/chuckha/aoc-solutions/internal/input"
)

// 4686774924
// 7030162386
//  781129154

// part 2 248367250691276 too low

type coords struct {
	a, b, c, d int
}

func main() {
	lines := input.GetInput(2019, 12)
	ogPlanets := []planet{}
	for i, line := range lines {
		ogPlanets = append(ogPlanets, newPlanetFromLine(i, line))
	}
	ogxcoords := coords{
		ogPlanets[0].pos.x,
		ogPlanets[1].pos.x,
		ogPlanets[2].pos.x,
		ogPlanets[3].pos.x,
	}
	ogycoords := coords{
		ogPlanets[0].pos.y,
		ogPlanets[1].pos.y,
		ogPlanets[2].pos.y,
		ogPlanets[3].pos.y,
	}
	ogzcoords := coords{
		ogPlanets[0].pos.z,
		ogPlanets[1].pos.z,
		ogPlanets[2].pos.z,
		ogPlanets[3].pos.z,
	}
	//	fmt.Println(ogPlanets)
	s := &system{
		planets: make([]planet, 0),
	}
	for i, line := range lines {
		s.planets = append(s.planets, newPlanetFromLine(i, line))
	}
	// xcoords := map[coords]int{}
	// ycoords := map[coords]int{}
	// zcoords := map[coords]int{}

	inc := 1
	one := false
	two := false
	three := false
	xo, yo, zo := 0, 0, 0
	zeroVel := coords{0, 0, 0, 0}
	for j := 0; j < 1000000; j += inc {
		// xcoords[s.xcoords()] = j
		// ycoords[s.ycoords()] = j
		// zcoords[s.zcoords()] = j
		s.applyGravity()
		// if _, ok := xcoords[s.xcoords()]; ok && !one {
		// 	xo = j + 1
		// 	fmt.Println("x returns to itself on", j+1)
		// 	one = true
		// }
		// if _, ok := ycoords[s.ycoords()]; ok && !two {
		// 	yo = j + 1
		// 	fmt.Println("y returns to itself on", j+1)
		// 	two = true
		// }
		// if _, ok := zcoords[s.zcoords()]; ok && !three {
		// 	zo = j + 1
		// 	fmt.Println("z returns to itself on", j+1)
		// 	three = true
		// }
		if s.xcoords() == ogxcoords && s.xvels() == zeroVel && !one {
			xo = j + 1
			fmt.Println("x returns to itself on", j+1)
			one = true
		}
		if s.ycoords() == ogycoords && s.yvels() == zeroVel && !two {
			yo = j + 1
			fmt.Println("y returns to itself on", j+1)
			two = true
		}
		if s.zcoords() == ogzcoords && s.zvels() == zeroVel && !three {
			zo = j + 1
			fmt.Println("z returns to itself on", j+1)
			three = true
		}
		if one && two && three {
			break
		}
	}
	fmt.Println(xo, yo, zo, xo*yo*zo)
	a := gcd(xo, yo)
	b := gcd(yo, zo)
	c := gcd(xo, zo)
	fmt.Println("a", a, "b", b, "c", c)
	g := gcd(gcd(xo, yo), zo)
	//	g2 := gcd(xo, gcd(yo, zo))
	if g > 0 {
		fmt.Println(xo/g, yo/g, zo/g, (xo/g)*(yo/g)*(zo/g))
		fmt.Println(xo/max(a, c), yo, zo/max(b, c), (xo/max(a, c))*(yo)*(zo/max(b, c)))
	}
	fmt.Println(math.MaxInt)
	//	fmt.Println(2772 % 43)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type system struct {
	planets []planet
	steps   int
}

func (s *system) xequal(ps []planet) bool {
	for i, p := range s.planets {
		if !p.xequal(ps[i]) {
			return false
		}
	}
	return true
}
func (s *system) yequal(ps []planet) bool {
	for i, p := range s.planets {
		if !p.yequal(ps[i]) {
			return false
		}
	}
	return true
}

func (s *system) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("After step %d\n", s.steps))
	for _, p := range s.planets {
		out.WriteString(p.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (s *system) applyGravity() {
	for i := 0; i < len(s.planets); i++ {
		for j := 0; j < len(s.planets); j++ {
			if i == j {
				continue
			}
			cur := s.planets[i]
			vx, vy, vz := cur.gravity(s.planets[j])
			cur.v.x += vx
			cur.v.y += vy
			cur.v.z += vz
			s.planets[i] = cur
		}
	}
	for i, p := range s.planets {
		p.pos = p.pos.addVeloicty(p.v)
		s.planets[i] = p
	}
	s.steps++
}
func (s *system) xcoords() coords {
	return coords{
		s.planets[0].pos.x,
		s.planets[1].pos.x,
		s.planets[2].pos.x,
		s.planets[3].pos.x,
	}
}
func (s *system) xvels() coords {
	return coords{
		s.planets[0].v.x,
		s.planets[1].v.x,
		s.planets[2].v.x,
		s.planets[3].v.x,
	}
}
func (s *system) yvels() coords {
	return coords{
		s.planets[0].v.y,
		s.planets[1].v.y,
		s.planets[2].v.y,
		s.planets[3].v.y,
	}
}
func (s *system) zvels() coords {
	return coords{
		s.planets[0].v.z,
		s.planets[1].v.z,
		s.planets[2].v.z,
		s.planets[3].v.z,
	}
}
func (s *system) ycoords() coords {
	return coords{
		s.planets[0].pos.y,
		s.planets[1].pos.y,
		s.planets[2].pos.y,
		s.planets[3].pos.y,
	}
}
func (s *system) zcoords() coords {
	return coords{
		s.planets[0].pos.z,
		s.planets[1].pos.z,
		s.planets[2].pos.z,
		s.planets[3].pos.z,
	}
}
func (s *system) totalEnergySum() int {
	sum := 0
	for _, p := range s.planets {
		sum += p.totalEnergy()
	}
	return sum
}

type planet struct {
	id  int
	pos point
	v   veloicty
}

func (p planet) xequal(p2 planet) bool {
	return p.pos.x == p2.pos.x && p.v.x == p2.v.x
}

func (p planet) yequal(p2 planet) bool {
	return p.pos.y == p2.pos.y && p.v.y == p2.v.y
}
func (p planet) equalPosition(p2 planet) bool {
	return p.pos.x == p2.pos.x && p.pos.y == p2.pos.y && p.pos.z == p2.pos.z
}
func (p planet) potentialEnergy() int {
	return internal.Abs(p.pos.x) + internal.Abs(p.pos.y) + internal.Abs(p.pos.z)
}
func (p planet) kineticEnergy() int {
	return internal.Abs(p.v.x) + internal.Abs(p.v.y) + internal.Abs(p.v.z)
}
func (p planet) totalEnergy() int {
	return p.potentialEnergy() * p.kineticEnergy()
}

func (p planet) String() string {
	return fmt.Sprintf("(%d) pos=<%v> vel=<%v>", p.id, p.pos, p.v)
}

func (p planet) gravity(p2 planet) (int, int, int) {
	x, y, z := 0, 0, 0
	if p.pos.x > p2.pos.x {
		x = -1
	}
	if p.pos.x < p2.pos.x {
		x = 1
	}
	if p.pos.y > p2.pos.y {
		y = -1
	}
	if p.pos.y < p2.pos.y {
		y = 1
	}
	if p.pos.z > p2.pos.z {
		z = -1
	}
	if p.pos.z < p2.pos.z {
		z = 1
	}
	return x, y, z
}

func newPlanetFromLine(id int, line string) planet {
	line = strings.Trim(line, "<>")
	coords := strings.Split(line, ", ")
	c := make([]int, 3)
	for i, coord := range coords {
		words := strings.Split(coord, "=")
		c[i], _ = strconv.Atoi(words[1])
	}
	return planet{
		id:  id,
		pos: point{x: c[0], y: c[1], z: c[2]},
	}
}

type point struct {
	x, y, z int
}

func (p point) addVeloicty(v veloicty) point {
	return point{p.x + v.x, p.y + v.y, p.z + v.z}
}

func (p point) String() string {
	return fmt.Sprintf("x=%d,y=%d,z=%d", p.x, p.y, p.z)
}

type veloicty struct {
	x, y, z int
}

func (v veloicty) String() string {
	return fmt.Sprintf("x=%d,y=%d,z=%d", v.x, v.y, v.z)
}

/*


p1 pos, vel
p1pos = p1.x + (p2.x>p1.x ? 1 : -1)
p2

<x=1, y=-0, z=0>
<x=10, y=0, z=0>
<x=100, y=0, z=0>


<x=-8, y=-10, z=0>
<x=5, y=5, z=10>
<x=2, y=-7, z=3>
<x=9, y=-8, z=-3>

x = -8, vx = 0
x = 5, vx = 0

x = -8, vx = 1
x = 5, vx = -1

x = -7 vx = 2
x = 4 vx = -2

x = -5 vx = 3
x = 2 vx = -3

x = -2 vx = 4
x = -1 vx = -4

x = 2 vx = 3
x = -5 vx = -3

x = 5 vx = 2
x = -8 vx = -2

x = 7 vx = 1
x = -10 vx = -1

x = 8 vx = 0
x = -11 vx = 0


sqrt of x distance is the number of steps until it repeats

so if the distance is -8 and 5, that's a distance of 13 which means
the sqrt (int) is 3.xxx, which means between 3 and 4 they pass each other.
they reach their apex at 3.xxx * 2 (in this case ~8)

they get back to their original position after the sqrt 4 times, so 3.xxx * 4 ~ 16

twice the sqrt gets them to switch coordinates

start -> apex from start -> end

3x 1x




*/

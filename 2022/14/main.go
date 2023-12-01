package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	input := internal.ReadInput()
	grid := internal.NewGridV3[*Cell]()
	sim := NewSimulation(grid)
	for _, line := range input {
		corners := pts(line)
		for i := 0; i < len(corners)-1; i++ {
			for _, p := range corners[i].Until(corners[i+1]) {
				sim.AddRock(p)
			}
		}
	}

	for sim.Tick() {
		//		fmt.Println(sim)
	}
	fmt.Println(sim)
	fmt.Println("sand that has come to rest", sim.SandCount())
}

func pts(line string) []internal.Point {
	parts := strings.Split(line, " -> ")
	out := make([]internal.Point, len(parts))
	for i, v := range parts {
		parts1 := strings.Split(v, ",")
		x, _ := strconv.Atoi(parts1[0])
		y, _ := strconv.Atoi(parts1[1])
		out[i] = internal.Point{X: x, Y: y}
	}
	return out
}

type Cell struct {
	SandOrigin bool
	HasSand    bool
}

func (c *Cell) String() string {
	if c.HasSand {
		return "o"
	}
	if c.SandOrigin {
		return "+"
	}
	return "#"
}

type Simulation struct {
	Rocks       *internal.GridV3[*Cell]
	Sand        *internal.GridV3[*Cell]
	SandOrigins *internal.GridV3[*Cell]
	Floor       int
}

func NewSimulation(rocks *internal.GridV3[*Cell]) *Simulation {
	sim := &Simulation{
		Rocks:       rocks,
		Sand:        internal.NewGridV3[*Cell](),
		SandOrigins: internal.NewGridV3[*Cell](),
		Floor:       rocks.Max.Y + 2,
	}
	sim.SandOrigins.Set(internal.Point{X: 500, Y: 0}, &Cell{SandOrigin: true})
	return sim
}

func (s *Simulation) Tick() bool {
	// generate a sand below sand origin
	for p := range s.SandOrigins.Data {
		s.Sand.Set(p, &Cell{HasSand: true})
	}
	// move sand down until there is something below it or it is 2 below the largest Y value in the grid
	for sand := range s.Sand.Data {
		moving := sand
		for s.SandCanMove(moving) && !s.OutOfBounds(moving) {
			for s.SandCanMoveDown(moving) {
				moving.Y = moving.Y + 1
			}
			if s.SandCanMoveDownLeft(moving) {
				moving.Y = moving.Y + 1
				moving.X = moving.X - 1
				continue
			}
			if s.SandCanMoveDownRight(moving) {
				moving.Y = moving.Y + 1
				moving.X = moving.X + 1
				continue
			}
		}
		s.Sand.Clear(sand)
		s.AddSand(moving)
	}
	return !s.Rocks.In(internal.Point{500, 0})
}

func (s *Simulation) String() string {
	return s.SandOrigins.Layer(s.Rocks).String()
}

func (s *Simulation) AddRock(pt internal.Point) {
	s.Rocks.Set(pt, &Cell{})
	s.Floor = s.Rocks.Max.Y + 2
}

func (s *Simulation) AddSand(pt internal.Point) {
	s.Rocks.Set(pt, &Cell{HasSand: true})
}

func (s *Simulation) SandCanMove(p internal.Point) bool {
	return s.SandCanMoveDown(p) || s.SandCanMoveDownLeft(p) || s.SandCanMoveDownRight(p)
}

func (s *Simulation) SandCanMoveDown(p internal.Point) bool {
	return !s.OutOfBounds(p.Down()) && !s.Rocks.In(p.Down())
}

func (s *Simulation) SandCanMoveDownLeft(p internal.Point) bool {
	return !s.OutOfBounds(p.Down().Left()) && !s.Rocks.In(p.Down().Left())
}

func (s *Simulation) SandCanMoveDownRight(p internal.Point) bool {
	return !s.OutOfBounds(p.Down().Right()) && !s.Rocks.In(p.Down().Right())
}

func (s *Simulation) OutOfBounds(p internal.Point) bool {
	return p.Y == s.Floor
}

func (s *Simulation) SandCount() int {
	count := 0
	for _, c := range s.Rocks.Data {
		if c.HasSand {
			count++
		}
	}
	return count
}

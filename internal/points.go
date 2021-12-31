package internal

type Point struct {
	X, Y int
}

func (p Point) Neighbors() []Point {
	return []Point{
		{X: p.X, Y: p.Y - 1},
		{X: p.X + 1, Y: p.Y},
		{X: p.X, Y: p.Y + 1},
		{X: p.X - 1, Y: p.Y},
	}
}

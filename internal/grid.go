package internal

type Grid struct {
	Data           map[Point]string
	Length, Height int
}

func (g *Grid) String() string {
	out := ""
	for j := 0; j < g.Height; j++ {
		for i := 0; i < g.Length; i++ {
			out += g.Data[Point{i, j}]
		}
		out += "\n"
	}
	return out
}

func NewGrid(length, height int, defaultChar string) *Grid {
	g := &Grid{
		Data:   map[Point]string{},
		Length: length,
		Height: height,
	}
	for j := 0; j < g.Height; j++ {
		for i := 0; i < g.Length; i++ {
			g.Data[Point{i, j}] = defaultChar
		}
	}
	return g
}

// RotateRow shifts a row to the right at position j by 'by' places
func (g *Grid) RotateRow(j, by int) {
	newPoints := map[Point]string{}
	for i := 0; i < g.Length; i++ {
		newI := (i + by) % g.Length
		newPoints[Point{newI, j}] = g.Data[Point{i, j}]
	}
	for k, v := range newPoints {
		g.Data[k] = v
	}
}

// RotateCol shifts a column down at position i by 'by' places
func (g *Grid) RotateCol(i, by int) {
	newPoints := map[Point]string{}
	for j := 0; j < g.Height; j++ {
		newJ := (j + by) % g.Height
		newPoints[Point{i, newJ}] = g.Data[Point{i, j}]
	}
	for k, v := range newPoints {
		g.Data[k] = v
	}
}

// Rect fills in a rectangle starting wiht the top left with the "on" character
func (g *Grid) Rect(x, y int) {
	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			g.Data[Point{i, j}] = "#"
		}
	}
}

// On counts the number of grid points that are on
func (g *Grid) On() int {
	count := 0
	for j := 0; j < g.Height; j++ {
		for i := 0; i < g.Length; i++ {
			if g.Data[Point{i, j}] == "#" {
				count++
			}
		}
	}
	return count

}

type Point struct {
	X, Y int
}

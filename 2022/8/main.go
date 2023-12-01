package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/chuckha/aoc-solutions/internal"
)

type Height int

func (h Height) String() string {
	return fmt.Sprintf("%d", h)
}

func main() {
	input := internal.ReadInput()
	trees := internal.NewGridV3[*Tree]()

	x := 0
	y := 0
	//	maxx := len(input[0]) - 1
	//	maxy := len(input) - 1
	for _, line := range input {
		for _, h := range line {
			treeHeight, err := strconv.Atoi(string(h))
			if err != nil {
				panic(err)
			}
			t := NewTree(treeHeight)
			trees.Set(internal.Point{X: x, Y: y}, t)
			x++
		}
		x = 0
		y++
	}
	scenicScore(trees)
	all := []*Tree{}
	trees.Each(func(t *Tree) {
		all = append(all, t)
	})
	sort.Sort(tz(all))
	fmt.Println(all[0], all[1])
}

type tz []*Tree

func (t tz) Len() int           { return len(t) }
func (t tz) Less(i, j int) bool { return t[i].sceneicScore > t[j].sceneicScore }
func (t tz) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func part1(trees *internal.GridV3[*Tree]) {
	visible(trees)
	fmt.Println(trees)
	count := 0
	trees.Each(func(t *Tree) {
		if t.Visible {
			count++
		}
	})
	fmt.Println(count)
}

type Tree struct {
	Height       int
	Visible      bool
	sceneicScore int
}

func (t *Tree) String() string {
	if t.Visible {
		return fmt.Sprintf("%d", t.Height)
	}
	return fmt.Sprintf("%d", t.sceneicScore)
}

func NewTree(h int) *Tree {
	return &Tree{Height: h}
}

func visible(g *internal.GridV3[*Tree]) {
	// for each row find the left visible one and stop when one is not visible
	for j := 0; j <= g.Max.Y; j++ {
		for i := 0; i <= g.Max.X; i++ {
			pt := internal.Point{i, j}
			cur := g.At(pt)
			if cur.Visible {
				continue
			}
			if i == 0 || j == 0 || i == g.Max.X || j == g.Max.Y {
				cur.Visible = true
				continue
			}
			// check up
			vis := true
			for k := pt.Y - 1; k >= 0; k-- {
				x := g.At(internal.Point{pt.X, k})
				if x.Height >= cur.Height {
					vis = false
					break
				}
			}
			if vis {
				cur.Visible = true
				continue
			}
			// check right
			vis = true
			for k := pt.X + 1; k <= g.Max.X; k++ {
				x := g.At(internal.Point{k, pt.Y})
				if x.Height >= cur.Height {
					vis = false
					break
				}
			}
			if vis {
				cur.Visible = true
				continue
			}
			// check down
			vis = true
			for k := pt.Y + 1; k <= g.Max.Y; k++ {
				x := g.At(internal.Point{pt.X, k})
				if x.Height >= cur.Height {
					vis = false
					break
				}
			}
			if vis {
				cur.Visible = true
				continue
			}
			// check left
			vis = true
			for k := pt.X - 1; k >= 0; k-- {
				x := g.At(internal.Point{k, pt.Y})
				if x.Height >= cur.Height {
					vis = false
					break
				}
			}
			if vis {
				cur.Visible = true
				continue
			}
		}
	}
}

func scenicScore(g *internal.GridV3[*Tree]) {
	// for each row find the left visible one and stop when one is not visible
	for j := 0; j <= g.Max.Y; j++ {
		for i := 0; i <= g.Max.X; i++ {
			pt := internal.Point{i, j}
			cur := g.At(pt)
			// check up
			up := 0
			for k := pt.Y - 1; k >= 0; k-- {
				x := g.At(internal.Point{pt.X, k})
				up++
				if x.Height >= cur.Height {
					break
				}
			}
			// check right
			right := 0
			for k := pt.X + 1; k <= g.Max.X; k++ {
				x := g.At(internal.Point{k, pt.Y})
				right++
				if x.Height >= cur.Height {
					break
				}
			}
			// check down
			down := 0
			for k := pt.Y + 1; k <= g.Max.Y; k++ {
				x := g.At(internal.Point{pt.X, k})
				down++
				if x.Height >= cur.Height {
					break
				}
			}
			// check left
			left := 0
			for k := pt.X - 1; k >= 0; k-- {
				x := g.At(internal.Point{k, pt.Y})
				left++
				if x.Height >= cur.Height {
					break
				}
			}
			cur.sceneicScore = up * right * down * left
		}
	}
}

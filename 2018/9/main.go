package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	line := internal.ReadInput()[0]
	words := strings.Split(line, " ")
	numPlayers, _ := strconv.Atoi(words[0])
	numMarbles, _ := strconv.Atoi(words[6])
	// part 2
	numMarbles *= 100
	players := make([]*player, numPlayers)
	for i := range players {
		players[i] = &player{}
	}
	root := internal.NewCircularLinkedList(0)
	game := &game{
		marbles:      root,
		players:      players,
		curPlayerIdx: 1,
		curMarbleVal: 1,
	}
	for game.curMarbleVal < numMarbles {
		//		fmt.Printf("[%d] %v\n", game.curPlayerIdx+1, root)
		game.turn()
	}
	scores := []int{}
	for _, p := range players {
		scores = append(scores, p.score)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	fmt.Println(scores[0])
}

type player struct {
	score int
}

type game struct {
	marbles      *internal.CircularLinkedList[int]
	players      []*player
	curPlayerIdx int
	curMarbleVal int
}

func (g *game) turn() {
	if g.curMarbleVal%23 == 0 {
		g.players[g.curPlayerIdx].score += g.curMarbleVal
		for i := 0; i < 7; i++ {
			g.marbles = g.marbles.Prev
		}
		g.players[g.curPlayerIdx].score += g.marbles.Data
		g.marbles = g.marbles.Remove()
		g.curMarbleVal++
		g.curPlayerIdx = (g.curPlayerIdx + 1) % len(g.players)
		return
	}
	g.marbles = g.marbles.Next.InsertAfter(g.curMarbleVal)
	g.curPlayerIdx = (g.curPlayerIdx + 1) % len(g.players)
	g.curMarbleVal++
}

// too low: 251876

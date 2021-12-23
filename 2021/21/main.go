package main

import (
	"fmt"
)

func main() {
	o := baseLookup()
	for k := range o {
		fmt.Println(k)
	}

	// log := &internal.Logger{false}
	// lines := internal.ReadInput()
	// w := strings.Split(lines[0], " ")
	// pos1, _ := strconv.Atoi(w[len(w)-1])
	// w2 := strings.Split(lines[1], " ")
	// pos2, _ := strconv.Atoi(w2[len(w2)-1])
	// player1 := &player{pos: pos1}
	// player2 := &player{pos: pos2}
	// cur := player1
	// die := &die{}
	// for {
	// 	fmt.Println("p1:", player1, "p2:", player2)
	// 	if player1.score >= 1000 {
	// 		log.Println("player 1 wins")
	// 		break
	// 	}
	// 	if player2.score >= 1000 {
	// 		log.Println("player 2 wins")
	// 		break
	// 	}
	// 	firstRoll := die.roll()
	// 	log.Println("first", firstRoll)
	// 	secondRoll := die.roll()
	// 	log.Println("second", secondRoll)
	// 	thirdRoll := die.roll()
	// 	log.Println("third", thirdRoll)
	// 	cur.update(firstRoll + secondRoll + thirdRoll)
	// 	if cur == player1 {
	// 		cur = player2
	// 		continue
	// 	}
	// 	cur = player1
	// }
	// if player1.score >= 1000 {
	// 	fmt.Println("p1 wins")
	// 	return
	// }
	// fmt.Println("p2 wins")
}

func previousStates(k key) key {
	curTurn := ((k.playerturn) % 2) + 1
	for i := 3; i <= 9; i++ {
		if k.playerturn == 1 {
			score1 := k.score1
			score2 := k.score2 - k.pos2
			pos1 := k.pos1
			pos2 := k.pos2 - i
		}
	}
	return key{
		playerturn: curTurn,
		score1:     0,
		score2:     0,
		pos1:       0,
		pos2:       0,
	}

}

type game struct {
	player1 *player
	player2 *player
}

type player struct {
	pos   int
	score int
}

func (p *player) String() string {
	return fmt.Sprintf("{%d %d}", p.pos, p.score)
}

func (p *player) update(spaces int) {
	newPos := (p.pos + spaces) % 10
	if newPos == 0 {
		newPos = 10
	}
	p.pos = newPos
	p.score += newPos
}

type die struct{}

func (d *die) roll() int {
	return 1
}

type lookup map[key]*wins

type key struct {
	playerturn     int
	score1, score2 int
	pos1, pos2     int
}

func (k key) String() string {
	if k.playerturn == 1 {
		return fmt.Sprintf("P1 %d@%d | P2 %d@%d", k.score1, k.pos1, k.score2, k.pos2)
	}
	return fmt.Sprintf("P2 %d@%d | P1 %d@%d", k.score2, k.pos2, k.score1, k.pos1)
}

type wins struct {
	p1 int
	p2 int
}

func (w *wins) String() string {
	return fmt.Sprintf("wins; p1: %d p2: %d", w.p1, w.p2)
}

func baseLookup() lookup {
	countMap := map[int]int{
		3: 1,
		4: 3,
		5: 6,
		6: 7,
		7: 6,
		8: 3,
		9: 1,
	}
	out := make(lookup)
	for score1 := 12; score1 < 20; score1++ {
		for score2 := 12; score2 < 20; score2++ {
			for pos1 := 1; pos1 <= 10; pos1++ {
				for pos2 := 1; pos2 <= 10; pos2++ {
					for turn := 1; turn <= 2; turn++ {
						out[key{
							playerturn: turn,
							score1:     score1,
							score2:     score2,
							pos1:       pos1,
							pos2:       pos2,
						}] = &wins{}
					}
				}
			}
		}
	}
	for k := range out {
		// simulate a roll
		for i := 3; i <= 9; i++ {
			if k.playerturn == 1 {
				if k.score1+i >= 21 {
					out[k].p1 += countMap[i]
				}
				continue
			}
			if k.score2+i >= 21 {
				out[k].p2 += countMap[i]
			}
		}
	}
	return out
}

/**

p1: 4
p2: 8

turn 1:
{7,8,9,10,1,2,3} / 8
{7,8,9,10,1,2,3} / {1,2,3,4,5,6,7}
n = 4
next turn {n+3 ... n-1} (all mod 10)
{n+3 ... n-1} / {}



990->999
pos 1-10
x 2 (current turn)


{
	pos
}

unique outcomes:
3, 4, 5, 6, 7, 8, 9 == 7 unique rolls per turn
which means
3 = 1 universe
4 = 3 universes
5 = 6 universes
6 = 7 universes
7 = 6 universes
8 = 3 universes
9 = 1 universe

all possible rolls
{1,1,1} = 3

{1,2,1} = 4
{1,1,2} = 4
{2,1,1} = 4

{1,1,3} = 5
{1,2,2} = 5
{1,3,1} = 5
{2,1,2} = 5
{2,2,1} = 5
{3,1,1} = 5

{1,2,3} = 6
{1,3,2} = 6
{2,1,3} = 6
{2,2,2} = 6
{2,3,1} = 6
{3,1,2} = 6
{3,2,1} = 6

{2,3,2} = 7
{2,2,3} = 7
{1,3,3} = 7
{3,1,3} = 7
{3,2,2} = 7
{3,3,1} = 7

{2,3,3} = 8
{3,3,2} = 8
{3,2,3} = 8

{3,3,3} = 9
*/

/*
p1: 4 {0}
p2: 8 {0}
p1 turn

from the beginning the number of wins for player 1 =

the number of wins for player 1 when he rolls a 3
+
the number of wins for player 1 when he rolls a 4
+
the number of wins for player 1 when he rolls a 5
+
the number of wins for player 1 when he rolls a 6
+
the number of wins for player 1 when he rolls a 7
+
the number of wins for player 1 when he rolls a 8
+
the number of wins for player 1 when he rolls a 9


*/

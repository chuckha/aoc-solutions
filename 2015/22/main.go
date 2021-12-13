package main

import (
	"fmt"
	"sort"

	"github.com/chuckha/aoc-solutions/internal"
)

const (
	playerTurn = iota
	bossTurn
)

func main() {
	me := &player{
		hitPoints: 50,
		mana:      500,
		manaSpent: 0,
	}
	spells := []spell{
		{
			name: "magic missle",
			cost: 53,
		},
		{
			name: "drain",
			cost: 73,
		},
		{
			name: "shield",
			cost: 113,
		},
		{
			name: "recharge",
			cost: 229,
		},
		{
			name: "poison",
			cost: 173,
		},
	}
	boss := &bossChar{
		hitPoints: 58,
		damage:    9,
	}
	q := internal.NewStack[*state]()
	//	q := internal.NewQueue[*state]()
	// sample game
	// me = &player{
	// 	hitPoints: 10,
	// 	mana:      250,
	// }
	// boss = &bossChar{
	// 	hitPoints: 14,
	// 	damage:    8,
	// }
	minimum := 99999999999999

	s := &state{
		player:  me,
		boss:    boss,
		effects: make(map[string]*effect),
		turn:    playerTurn,
		//		debug:   true,
	}
	manaSpends := []int{}
	q.Push(s)
	//	q.Enqueue(s)
	for !q.Empty() {
		//		turn := q.Dequeue()
		turn := q.Pop()
		if turn.player.manaSpent > minimum {
			continue
		}
		if turn.debug && turn.lastSpell != "" {
			fmt.Printf("Player casts %s\n", turn.lastSpell)
			turn.lastSpell = ""
		}

		if turn.debug {
			fmt.Println(turn)
		}
		newTurn := turn.copy()
		newTurn.player.hitPoints -= 1
		if newTurn.player.hitPoints <= 0 {
			if newTurn.debug {
				fmt.Println("dead human")
				continue
			}
		}
		newTurn.runEffects()
		if newTurn.boss.hitPoints <= 0 {
			if newTurn.player.manaSpent < minimum {
				minimum = newTurn.player.manaSpent
			}
			manaSpends = append(manaSpends, newTurn.player.manaSpent)
			if newTurn.debug {
				fmt.Println("dead boss")
			}
			continue
		}
		newTurn.removeEndedEffects()

		if turn.turn == bossTurn {
			newTurn.bossAttack()
			if newTurn.player.hitPoints <= 0 {
				if newTurn.debug {
					fmt.Println("dead human")
				}
				continue
			}
			newTurn.turn = playerTurn
			//			q.Enqueue(newTurn)
			q.Push(newTurn)
			continue
		}

		possibleSpells := newTurn.validSpells(spells)
		// we can't cast any spells so the game is over
		if len(possibleSpells) == 0 {
			continue
		}

		for _, s := range possibleSpells {
			t := newTurn.copy()
			t.cast(s)
			if t.boss.hitPoints <= 0 {
				if t.player.manaSpent < minimum {
					minimum = t.player.manaSpent
				}
				manaSpends = append(manaSpends, t.player.manaSpent)
				if t.debug {
					fmt.Println("dead boss")
				}
				continue
			}
			if t.player.hitPoints <= 0 {
				if t.debug {
					fmt.Println("dead human")
				}
				continue
			}
			t.lastSpell = s.name
			t.turn = bossTurn
			q.Push(t)
			//			q.Enqueue(t)
		}
	}
	///	fmt.Println(minimum)
	//	fmt.Println(manaSpends)
	sort.Sort(sort.IntSlice(manaSpends))
	fmt.Println(manaSpends[0])
	// queue a state
	// while the queue isn't empty
	// if this game state has spent more mana than the current minimum, just move on
	// game state is player, boss, turn count, effect

}

type player struct {
	hitPoints int
	defense   int
	mana      int
	manaSpent int
}

func (p *player) copy() *player {
	return &player{
		hitPoints: p.hitPoints,
		defense:   p.defense,
		mana:      p.mana,
		manaSpent: p.manaSpent,
	}
}

type bossChar struct {
	hitPoints int
	damage    int
}

func (b *bossChar) copy() *bossChar {
	return &bossChar{
		hitPoints: b.hitPoints,
		damage:    b.damage,
	}
}

type state struct {
	player    *player
	boss      *bossChar
	effects   map[string]*effect
	turn      int
	lastSpell string
	debug     bool
}

func (s *state) bossAttack() {
	dmg := s.boss.damage - s.player.defense
	if s.debug {
		fmt.Printf("boss attacks for %d damage\n", dmg)
	}
	s.player.hitPoints -= dmg
}

func (s *state) copy() *state {
	effects := make(map[string]*effect)
	for k := range s.effects {
		effects[k] = s.effects[k].copy()
	}
	return &state{
		player:    s.player.copy(),
		boss:      s.boss.copy(),
		effects:   effects,
		turn:      s.turn,
		lastSpell: s.lastSpell,
		debug:     s.debug,
	}
}

func (s *state) runEffects() {
	for name, effect := range s.effects {
		switch name {
		case "shield":
			if s.debug {
				fmt.Printf("Shield's timer is now %d\n", effect.turns-1)
			}
		case "poison":
			if s.debug {
				fmt.Printf("Poison deals 3 damage; new timer is %d\n", effect.turns-1)
			}
			s.boss.hitPoints -= 3
		case "recharge":
			if s.debug {
				fmt.Printf("Recharge provides 101 mana, new timer is %d\n", effect.turns-1)
			}
			s.player.mana += 101
		default:
			panic("very bad")
		}
		effect.turns -= 1
	}
}

func (s *state) removeEndedEffects() {
	removeKeys := []string{}
	for k, e := range s.effects {
		if e.turns == 0 {
			removeKeys = append(removeKeys, k)
		}
	}
	for _, k := range removeKeys {
		if k == "shield" {
			s.player.defense -= 7
		}
		delete(s.effects, k)
	}
}

func (s *state) cast(spell spell) {
	switch spell.name {
	case "magic missle":
		s.player.mana -= 53
		s.player.manaSpent += 53
		s.boss.hitPoints -= 4
	case "drain":
		s.player.mana -= 73
		s.player.manaSpent += 73
		s.boss.hitPoints -= 2
		s.player.hitPoints += 2
	case "shield":
		s.player.mana -= 113
		s.player.manaSpent += 113
		s.player.defense += 7
		s.effects["shield"] = &effect{
			name:  "sheild",
			turns: 6,
		}
	case "recharge":
		s.player.mana -= 229
		s.player.manaSpent += 229
		s.effects["recharge"] = &effect{
			name:  "recharge",
			turns: 5,
		}
	case "poison":
		s.player.mana -= 173
		s.player.manaSpent += 173
		s.effects["poison"] = &effect{
			name:  "poison",
			turns: 6,
		}
	default:
		panic(fmt.Sprintf("unknown spell %q", spell.name))
	}

}

func (s *state) validSpells(spells []spell) []spell {
	out := make([]spell, 0)
	for _, spell := range spells {
		// make sure the spell is affordable
		if spell.cost > s.player.mana {
			continue
		}
		// cannot cast a spell that's already active
		if _, ok := s.effects[spell.name]; ok {
			continue
		}
		out = append(out, spell)
	}
	return out
}

func (s *state) String() string {
	turnText := "\n-- Player turn --"
	if s.turn == bossTurn {
		turnText = "\n-- Boss turn --"
	}
	out := fmt.Sprintf("%s\n", turnText)
	out += fmt.Sprintf("- Player: hit points: %d, armor: %d, mana: %d\n", s.player.hitPoints, s.player.defense, s.player.mana)
	out += fmt.Sprintf("- Boss: hit points: %d", s.boss.hitPoints)
	return out
}

type effect struct {
	name  string
	turns int
}

func (e *effect) copy() *effect {
	return &effect{
		name:  e.name,
		turns: e.turns,
	}
}

type spell struct {
	name string
	cost int
}

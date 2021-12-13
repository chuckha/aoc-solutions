package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()

	boss := &player{}
	for _, line := range lines {
		words := strings.Split(line, ": ")
		val, _ := strconv.Atoi(words[1])
		switch words[0] {
		case "Hit Points":
			boss.hitPoints = val
		case "Damage":
			boss.damage = val
		case "Armor":
			boss.defense = val
		default:
			panic("unkonwn input: " + words[0])
		}
	}
	s := newStore()
	me := &player{
		hitPoints: 100,
		spent:     0,
	}
	spent := []int{}
	players := []player{}
	for _, w := range s.weapons {
		// copy player
		newMe := me.copy()
		newMe.addWeapon(w)
		if victory(newMe, boss) {
			players = append(players, *newMe)
			spent = append(spent, newMe.spent)
		}
		for _, a := range s.armor {
			weaponized := newMe.copy()
			weaponized.addArmor(a)
			if victory(weaponized, boss) {
				players = append(players, *weaponized)
				spent = append(spent, weaponized.spent)
			}
			for _, r1 := range s.rings {
				armored := weaponized.copy()
				armored.addRing(r1)
				if victory(armored, boss) {
					players = append(players, *armored)
					spent = append(spent, armored.spent)
				}
				for _, r2 := range s.rings {
					if r1 == r2 {
						continue
					}
					oneRing := armored.copy()
					oneRing.addRing(r2)
					if victory(oneRing, boss) {
						players = append(players, *oneRing)
						spent = append(spent, oneRing.spent)
					}
				}
			}
		}
		// no armor case
		for _, r1 := range s.rings {
			weaponized := newMe.copy()
			weaponized.addRing(r1)
			if victory(weaponized, boss) {
				players = append(players, *weaponized)
				spent = append(spent, weaponized.spent)
			}
			for _, r2 := range s.rings {
				if r1 == r2 {
					continue
				}
				oneRing := weaponized.copy()
				oneRing.addRing(r2)
				if victory(oneRing, boss) {
					players = append(players, *oneRing)
					spent = append(spent, oneRing.spent)
				}
			}
		}
	}
	sort.Sort(sort.Reverse(playerVictories(players)))
	// for _, p := range players {
	// 	fmt.Println(p)
	// }
	fmt.Println(players[0].spent)
}

type playerVictories []player

func (p playerVictories) Len() int { return len(p) }
func (p playerVictories) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p playerVictories) Less(i, j int) bool {
	return p[i].spent < p[j].spent
}

type player struct {
	hitPoints int
	damage    int
	defense   int
	spent     int
	weapon    weapon
	armor     armor
	rings     []ring
}

func (p *player) copy() *player {
	return &player{
		hitPoints: p.hitPoints,
		damage:    p.damage,
		defense:   p.defense,
		spent:     p.spent,
		weapon:    p.weapon,
		armor:     p.armor,
		rings:     p.rings,
	}
}

func (p *player) addWeapon(w weapon) {
	p.damage += w.damage
	p.spent += w.cost
	p.weapon = w
}
func (p *player) addArmor(a armor) {
	p.defense += a.armor
	p.spent += a.cost
	p.armor = a
}
func (p *player) addRing(r ring) {
	if len(p.rings) == 2 {
		return
	}
	p.spent += r.cost
	p.damage += r.damage
	p.defense += r.armor
	p.rings = append(p.rings, r)
}

func victory(me, boss *player) bool {
	if boss.damage-me.defense <= 0 {
		return false
	}
	if me.damage-boss.defense <= 0 {
		return true
	}
	myDeath := me.hitPoints / (boss.damage - me.defense)
	bossDeath := boss.hitPoints / (me.damage - boss.defense)

	return myDeath < bossDeath
}

type store struct {
	weapons []weapon
	armor   []armor
	rings   []ring
}

type weapon struct {
	name   string
	cost   int
	damage int
	armor  int
}
type armor struct {
	name   string
	cost   int
	damage int
	armor  int
}
type ring struct {
	name   string
	cost   int
	damage int
	armor  int
}

func newStore() *store {
	data := `Weapons:    Cost  Damage  Armor
Dagger        8     4       0
Shortsword   10     5       0
Warhammer    25     6       0
Longsword    40     7       0
Greataxe     74     8       0

Armor:      Cost  Damage  Armor
Leather      13     0       1
Chainmail    31     0       2
Splintmail   53     0       3
Bandedmail   75     0       4
Platemail   102     0       5

Rings:      Cost  Damage  Armor
Damage+1    25     1       0
Damage+2    50     2       0
Damage+3   100     3       0
Defense+1   20     0       1
Defense+2   40     0       2
Defense+3   80     0       3`
	lines := strings.Split(data, "\n")
	thing := 0
	s := &store{
		weapons: []weapon{},
		armor:   []armor{},
		rings:   []ring{},
	}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.HasPrefix(line, "Weapons") {
			continue
		}
		if strings.HasPrefix(line, "Armor") {
			thing++
			continue
		}
		if strings.HasPrefix(line, "Rings") {
			thing++
			continue
		}
		data := strings.Fields(line)
		cost, _ := strconv.Atoi(data[1])
		damage, _ := strconv.Atoi(data[2])
		a, _ := strconv.Atoi(data[3])
		switch thing {
		case 0:
			s.weapons = append(s.weapons, weapon{data[0], cost, damage, a})
		case 1:
			s.armor = append(s.armor, armor{data[0], cost, damage, a})
		case 2:
			s.rings = append(s.rings, ring{data[0], cost, damage, a})
		default:
			panic("bad thing")
		}
	}
	return s
}

/*
hp8, d5, a5
hp12, d7, a2

y1 = 8-2x1 player health
y2 = 12-3x2 boss health
win scenario: y1 = 0; x1 >= x2
lose scenario: y1 = 0; x1 < x2

0 = 8 - 2x1
-8 = -2x1
4 = x1

0 = 12 - 3x2
3x2 = 12
4 = x2

*/

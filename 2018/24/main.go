package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

const (
	Infection    = "Infection"
	ImmuneSystem = "Immune System"
)

type results struct {
	boost  int
	winner string
}

func lastOppositeWin(in []results, curwinner string) results {
	for i := len(in) - 1; i >= 0; i-- {
		if in[i].winner == curwinner {
			return in[i]
		}
	}
	return results{}
}

func main() {
	lines := internal.ReadRealRawInput()
	boosts := []results{}
	curBoost := 1
	finalBoost := 0
	totallyDone := false
	for !totallyDone {
		fmt.Println("Boosting by", curBoost)
		armies := newArmiesFromInput(lines)
		armies[ImmuneSystem].boost(curBoost)
		done := false
		restart := false
		for {
			// for _, army := range armies {
			// 	fmt.Println(army.name)
			// 	for _, g := range army.groups {
			// 		fmt.Println(g.DetailString())
			// 	}
			// }

			out := fight(armies)
			if len(out) == 0 {
				fmt.Println("no targets possible")
				// for _, army := range armies {
				// 	fmt.Println(army.name)
				// 	for _, g := range army.groups {
				// 		fmt.Println(g.DetailString())
				// 	}
				// }
				curBoost++
				restart = true
				break
			}
			// for _, v := range out {
			// 	fmt.Printf("attacking group: %v, defending group: %v, (%d)\n", v.attacker, v.defender, v.attacker.wouldDealDamage(v.defender))
			// }
			sort.Sort(attacks(out))

			for _, o := range out {
				_ = o.attacker.attack(o.defender)
				// fmt.Println(o.attacker, "is attacking", o.defender, "killing", dead, "units")
			}
			for _, army := range armies {
				if army.beaten() {
					done = true
				}
				for _, g := range army.groups {
					g.resetRound()
				}
			}
			if done {
				break
			}
		}
		if restart {
			continue
		}

		for _, army := range armies {
			sum := 0
			for _, g := range army.groups {
				sum += g.units
			}
			if sum > 0 {
				if army.name == ImmuneSystem {
					result := results{boost: curBoost, winner: army.name}
					fmt.Println(army.name, "wins with", sum, "units", result)
					boosts = append(boosts, result)
					if len(boosts) == 1 {
						fmt.Println("finished with a single point boost")
						fmt.Println(army.name, "has", sum, "units")
						break
					}
					lastBoost := lastOppositeWin(boosts, Infection)
					//					fmt.Println(lastBoost.boost, curBoost, lastBoost.boost-curBoost == 1)
					if abs(lastBoost.boost-curBoost) == 1 {
						finalBoost = curBoost
						totallyDone = true
					}
					//					fmt.Println("reducing (", lastBoost.boost, "+", curBoost, ")/2")
					//	curBoost = (curBoost + lastBoost.boost) / 2
					curBoost++
				}
				if army.name == Infection {
					result := results{boost: curBoost, winner: army.name}
					fmt.Println(army.name, "wins with", sum, "units", result)
					boosts = append(boosts, result)
					if len(boosts) == 1 {
						//						curBoost = curBoost * 2
						curBoost++
						continue
					}
					last := lastOppositeWin(boosts, ImmuneSystem)
					if last.winner == "" {
						//						curBoost *= 2
						curBoost++
						continue
					}
					if abs(last.boost-curBoost) == 1 {
						finalBoost = last.boost
						totallyDone = true
					}
					//					fmt.Println("reducing (", last.boost, "+", curBoost, ")/2")
					//curBoost = (last.boost + curBoost) / 2
					curBoost++
				}
				if totallyDone {
					fmt.Println("part 2:", finalBoost)
				}
			}
		}

	}
	// TODO: clear boolena flags on groups (targeted)
}

func abs(a int) int {
	if a < 0 {
		return a * -1
	}
	return a
}

type attacks []*targetingPair

func (a attacks) Len() int      { return len(a) }
func (a attacks) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a attacks) Less(i, j int) bool {
	return a[i].attacker.initiative > a[j].attacker.initiative
}
func fight(armies map[string]*army) []*targetingPair {
	pairings := []*targetingPair{}
	for _, army := range armies {
		allGroups := []*group{}
		for _, g := range army.groups {
			if !g.active() {
				continue
			}
			allGroups = append(allGroups, g)
		}
		sort.Sort(groups(allGroups))
		for _, g := range allGroups {
			// blorg
			switch g.army {
			case Infection:
				defender := armies[ImmuneSystem].defender(g)
				if defender == nil {
					continue
				}
				pairings = append(pairings, &targetingPair{
					attacker: g,
					defender: defender,
				})
			case ImmuneSystem:
				defender := armies[Infection].defender(g)
				if defender == nil {
					continue
				}
				pairings = append(pairings, &targetingPair{
					attacker: g,
					defender: defender,
				})
			default:
				panic("bad army")
			}
		}
	}
	return pairings
}

type army struct {
	name   string
	groups map[int]*group
}

func (a *army) beaten() bool {
	count := 0
	for _, g := range a.groups {
		count += g.units
	}
	return count == 0
}
func (a *army) boost(i int) {
	for _, g := range a.groups {
		g.damage += i
	}
}
func weaknessAndImmunities(item string) ([]string, []string) {
	weaknesses := []string{}
	immunities := []string{}
	parts := strings.Split(item, "; ")
	words1 := strings.Split(parts[0], " ")
	if words1[0] == "immune" {
		immunities = strings.Split(strings.Split(parts[0], " to ")[1], ", ")
		if len(parts) > 1 {
			weaknesses = strings.Split(strings.Split(parts[1], " to ")[1], ", ")
		}
	}
	if words1[0] == "weak" {
		weaknesses = strings.Split(strings.Split(parts[0], " to ")[1], ", ")
		if len(parts) > 1 {
			immunities = strings.Split(strings.Split(parts[1], " to ")[1], ", ")
		}
	}
	return weaknesses, immunities
}

func newArmy() *army {
	return &army{
		groups: make(map[int]*group),
	}
}

func (a *army) defender(attacking *group) *group {
	possibilities := []*group{}
	for _, g := range a.groups {
		// ignore already targted groups
		if g.targeted || !g.active() {
			continue
		}
		possibilities = append(possibilities, g)
	}
	dg := &defenderGroups{groups: possibilities, attacking: attacking}
	sort.Sort(dg)
	if len(dg.groups) == 0 {
		return nil
	}
	if attacking.wouldDealDamage(dg.groups[0]) == 0 {
		return nil
	}
	dg.groups[0].targeted = true
	return dg.groups[0]
}

type defenderGroups struct {
	groups    []*group
	attacking *group
}

func (d *defenderGroups) Len() int      { return len(d.groups) }
func (d *defenderGroups) Swap(i, j int) { d.groups[i], d.groups[j] = d.groups[j], d.groups[i] }
func (d *defenderGroups) Less(i, j int) bool {
	if d.attacking.wouldDealDamage(d.groups[i]) > d.attacking.wouldDealDamage(d.groups[j]) {
		return true
	}
	if d.attacking.wouldDealDamage(d.groups[i]) < d.attacking.wouldDealDamage(d.groups[j]) {
		return false
	}
	if d.groups[i].effectivePower() > d.groups[j].effectivePower() {
		return true
	}
	if d.groups[i].effectivePower() < d.groups[j].effectivePower() {
		return false
	}
	return d.groups[i].initiative > d.groups[j].initiative
}

type targetingPair struct {
	attacker *group
	defender *group
}

type groups []*group

func (g groups) Len() int      { return len(g) }
func (g groups) Swap(i, j int) { g[i], g[j] = g[j], g[i] }
func (g groups) Less(i, j int) bool {
	if g[i].effectivePower() > g[j].effectivePower() {
		return true
	}
	if g[i].effectivePower() < g[j].effectivePower() {
		return false
	}
	return g[i].initiative > g[j].initiative
}

type group struct {
	id         int
	army       string
	units      int
	hitPoints  int
	damage     int
	weaknesses []string
	immunities []string
	attackType string
	initiative int
	targeted   bool
}

func (g *group) immune(g2 *group) bool {
	for _, im := range g.immunities {
		if im == g2.attackType {
			return true
		}
	}
	return false
}
func (g *group) weakTo(g2 *group) bool {
	for _, w := range g.weaknesses {
		if w == g2.attackType {
			return true
		}
	}
	return false
}

func (g *group) wouldDealDamage(defender *group) int {
	base := g.effectivePower()
	if defender.immune(g) {
		//		fmt.Println(g, "would deal 0 damage to", defender)
		return 0
	}
	if defender.weakTo(g) {
		///		fmt.Println(g, "would deal double (", 2*base, ") damage to", defender)
		return 2 * base
	}
	//	fmt.Println(g, "would deal regular (", base, ") damage to", defender)
	return base
}

func (g *group) effectivePower() int {
	return g.units * g.damage
}
func (g *group) active() bool {
	return g.units > 0
}
func (g *group) String() string {
	return fmt.Sprintf("%s group %d", g.army, g.id)
}
func (g *group) DetailString() string {
	return fmt.Sprintf("%s group %d UNITS: %d HP: %d WEK: %v (%d) IMM: %v (%d) ATK: %d (%s) INIT: %d EP: (%d)", g.army, g.id, g.units, g.hitPoints, g.weaknesses, len(g.weaknesses), g.immunities, len(g.immunities), g.damage, g.attackType, g.initiative, g.effectivePower())
}
func (g *group) attack(defender *group) int {
	if !g.active() || defender == nil {
		return -1
	}
	dmg := g.wouldDealDamage(defender)
	willKill := dmg / defender.hitPoints
	if willKill > defender.units {
		willKill = defender.units
	}
	defender.units -= willKill
	return willKill
}
func (g *group) resetRound() {
	g.targeted = false
}

func newArmiesFromInput(lines []string) map[string]*army {
	armies := map[string]*army{}
	currentArmy := newArmy()
	identifier := 1
	for _, line := range lines {
		if strings.HasSuffix(line, ":") {
			currentArmy.name = strings.TrimSuffix(line, ":")
			continue
		}
		if len(line) == 0 {
			identifier = 1
			armies[currentArmy.name] = currentArmy
			currentArmy = newArmy()
			continue
		}
		w := []string{}
		i := []string{}
		//2313 units each with 6792 hit points (weak to fire, radiation; immune to cold) with an attack that does 29 bludgeoning damage at initiative 9
		// 1117 units each with 5042 hit points with an attack that does 44 fire damage at initiative 15
		openParenIdx := strings.Index(line, "(")
		closeParenIdx := strings.Index(line, ")")
		cleaned := line
		if openParenIdx != -1 {
			cleaned = line[:openParenIdx] + line[closeParenIdx+2:]
			w, i = weaknessAndImmunities(line[openParenIdx+1 : closeParenIdx])
		}
		words := strings.Split(cleaned, " ")
		units, _ := strconv.Atoi(words[0])
		hitPoints, _ := strconv.Atoi(words[4])
		damage, _ := strconv.Atoi(words[12])
		attackType := words[13]
		initiative, _ := strconv.Atoi(words[17])
		currentArmy.groups[identifier] = &group{
			id:         identifier,
			army:       currentArmy.name,
			units:      units,
			hitPoints:  hitPoints,
			damage:     damage,
			weaknesses: w,
			immunities: i,
			attackType: attackType,
			initiative: initiative,
		}
		identifier++
	}
	return armies
}

/*
Immune System:
17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 4507 fire damage at initiative 2
989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 25 slashing damage at initiative 3

Infection:
801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1
4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4

*/

/*
Immune System:
17 units each with 5390 hit points (weak to radiation, bludgeoning) with an attack that does 6077 fire damage at initiative 2
989 units each with 1274 hit points (immune to fire; weak to bludgeoning, slashing) with an attack that does 1595 slashing damage at initiative 3

Infection:
801 units each with 4706 hit points (weak to radiation) with an attack that does 116 bludgeoning damage at initiative 1
4485 units each with 2961 hit points (immune to radiation; weak to fire, cold) with an attack that does 12 slashing damage at initiative 4
*/

package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

/*

 ready chan
workers grab items off ready chan
when they finish the write them to a done chan


when the updater reads an item off done, it checks to see if any other steps are ready
if they are ready it puts it on the ready chan

something needs to coordinate tickers

*/
func exampletime(n string) int {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return strings.Index(letters, n) + 1
}
func main() {
	numWorkers := 5
	lines := internal.ReadInput()
	allSteps := make(steps, 0)

	for _, line := range lines {
		stepName, dependencyName := parseLine(line)
		s := allSteps.get(stepName)
		if s == nil {
			s = &step{name: stepName, workTime: doTime(stepName)}
			// s = &step{name: stepName, workTime: exampletime(stepName)}
			allSteps = append(allSteps, s)
		}
		dependency := allSteps.get(dependencyName)
		if dependency == nil {
			dependency = &step{name: dependencyName, workTime: doTime(dependencyName)}
			// dependency = &step{name: dependencyName, workTime: exampletime(dependencyName)}
			allSteps = append(allSteps, dependency)
		}
		s.dependencies = append(s.dependencies, dependency)
	}
	ticks := 0
	wkrs := make([]*worker, numWorkers)
	for i := range wkrs {
		wkrs[i] = &worker{}
	}
	for len(findReady(allSteps)) != 0 || !workers(wkrs).noneActive() {
		ready := findReady(allSteps)

		for _, w := range wkrs {
			if w.active() {
				continue
			}
			for _, r := range ready {
				w.s = r
				w.s.working = true
				break
			}
			ready = findReady(allSteps)
		}
		// for _, w := range wkrs {
		// 	fmt.Print("active?", w.active())
		// 	if w.active() {
		// 		fmt.Print(w.s)
		// 	}
		// 	fmt.Println()
		// }
		workers(wkrs).tick()
		ticks++
		//		fmt.Println(workers(wkrs).noneActive())

		for _, w := range wkrs {
			if !w.active() {
				continue
			}
			if w.done() {
				fmt.Println(ticks, w.s.name, "is done")
				w.s.done = true
				w.s = nil
			}
		}
	}
	fmt.Println(ticks)

	// for _, v := range allSteps {
	// 	fmt.Println(v)
	// }
	// //	done := []string{}
	// ticks := 0
	// active := make(steps, workers-1)
	// for len(findReady(allSteps)) != 0 {
	// 	ready := findReady(allSteps)
	// 	for _, r := range ready {
	// 		fmt.Println(ticks, r)
	// 	}
	// 	cur := ready[0]
	// 	ready = ready[1:]
	// 	for j := 0; j < len(active); j++ {
	// 		if active[j] == nil && j < len(ready) {
	// 			active[j] = ready[j]
	// 			ready[j].working = true
	// 		}
	// 	}
	// 	ticks += cur.workTime
	// 	cur.done = true
	// 	fmt.Println("hello", findReady(allSteps))
	// 	fmt.Println(ticks, cur.name, "is done")
	// 	for j := 0; j < cur.workTime; j++ {
	// 		for i := 0; i < len(active); i++ {
	// 			if active[i] == nil {
	// 				continue
	// 			}
	// 			active[i].workTime -= 1
	// 		}
	// 		for i := 0; i < len(active); i++ {
	// 			if active[i] == nil {
	// 				continue
	// 			}

	// 			if active[i].workTime == 0 {
	// 				active[i].done = true
	// 				fmt.Println(ticks, active[i].name, "is donee")
	// 				active[i] = nil
	// 				if len(ready) > 0 {
	// 					active[i] = ready[0]
	// 					ready[0].working = true
	// 					ready = ready[1:]
	// 				}
	// 				if len(ready) == 0 {
	// 					ready = findReady(allSteps)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// for i := 0; i < len(ready); i++ {
	// 	ticks += ready[i].workTime
	// 	fmt.Println(ready[i].workTime, "has elapsed", "length", len(ready))
	// 	ready[i].done = true
	// 	fmt.Println(ready[i].name, "is done")
	// 	// go from i+1 -> i+workers or len(ready)
	// 	for j := i + 1; j < min(i+workers, len(ready)); j++ {
	// 		fmt.Println("removing", ready[i].workTime, "ticks from", ready[j].name)
	// 		ready[j].workTime -= ready[i].workTime
	// 	}
	// 	for j := i + 1; j < min(i+workers, len(ready)); j++ {
	// 		if ready[j].workTime == 0 {
	// 			ready[j].done = true
	// 			fmt.Println(ready[j].name, "is done")
	// 		}
	// 	}
	// 	ready = findReady(allSteps)
	// 	for _, r := range ready {
	// 		fmt.Println(ticks, "ready", r)
	// 	}
	// 	i = -1
	// }
	fmt.Println(ticks)
	// find available
}

type active struct {
	steps map[string]*step
}

func (a *active) addStep(step *step) {
	a.steps[step.name] = step
}
func (a *active) tick() {
	for _, s := range a.steps {
		s.workTime -= 1
	}
	a.checkForDone()
}
func (a *active) checkForDone() {
	for k, s := range a.steps {
		if s.workTime == 0 {
			delete(a.steps, k)
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func findReady(s steps) steps {
	out := make(steps, 0)
	for _, step := range s {
		if step.ready() {
			out = append(out, step)
		}
	}
	sort.Sort(out)
	return out
}

func doTime(name string) int {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return strings.Index(letters, name) + 61
}

func run(available []string, allSteps steps) {
	for i := 0; i < len(available); i++ {
		step := allSteps.get(available[i])
		if step.done {
			continue
		}
		step.do()
		// get all the ones that need that first one
		allThatNeedStep := allSteps.getAllThatNeed(step.name)
		for _, a2 := range allThatNeedStep {
			s := allSteps.get(a2)
			if s.done {
				continue
			}
			if s.ready() {
				available = append(available, s.name)
			}
		}
		sort.Sort(sort.StringSlice(available))
		i = -1
	}
}

type step struct {
	name         string
	done         bool
	dependencies []*step
	workTime     int
	working      bool
}

func (s *step) String() string {
	out := []string{}
	for _, st := range s.dependencies {
		out = append(out, st.name)
	}
	done := "not done"
	if s.done {
		done = "done"
	}
	return fmt.Sprintf("%s %s (%d) needs %v", done, s.name, s.workTime, strings.Join(out, " "))
}

func (s *step) ready() bool {
	if s.done || s.working {
		return false
	}
	if len(s.dependencies) == 0 {
		return true
	}
	for _, dep := range s.dependencies {
		if !dep.done {
			return false
		}
	}
	return true
}

func (s *step) do() {
	if s.done {
		panic("shouldn't call do twice:" + s.name)
	}
	for _, dep := range s.dependencies {
		if !dep.done {
			panic("not ready")
		}
	}
	fmt.Printf(s.name)
	s.done = true
}

// returns a depends on b
func parseLine(line string) (string, string) {
	words := strings.Split(line, " ")
	return words[7], words[1]
}

type steps []*step

func (s steps) Len() int           { return len(s) }
func (s steps) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s steps) Less(i, j int) bool { return s[i].workTime < s[j].workTime }
func (s steps) get(name string) *step {
	for _, step := range s {
		if step.name == name {
			return step
		}
	}
	return nil
}

func (s steps) getAllThatNeed(e string) []string {
	out := []string{}
	for _, step := range s {
		if step.done {
			continue
		}
		for _, d := range step.dependencies {
			if d.name == e {
				out = append(out, step.name)
			}
		}
	}
	return out
}

// wrong guesses
// EUGJKYFQWSCLTXNIZMAPVORDBH
// correct
// EUGJKYFQSCLTWXNIZMAPVORDBH

/*
time    w1   w2   w3   w4   w5  done
0       E 65 U 80 Y 84
...
64      E    U   Y
65           U   Y               E
...
79           U   Y				E
80      G    K   Y               EU
*/

// 787 is too low
// 818 is too low
// 1149 is too high
// 1133 is wrong
// 1002 is wrong
// 1000 is wrong
// 1014?

type workers []*worker

func (w workers) tick() {
	for _, worker := range w {
		if worker.active() {
			worker.tick()
		}
	}
}
func (w workers) noneActive() bool {
	for _, wkr := range w {
		if wkr.active() {
			return false
		}
	}
	return true
}

type worker struct {
	s *step
}

func (w *worker) tick() {
	w.s.workTime -= 1
}

func (w *worker) active() bool { return w.s != nil }
func (w *worker) done() bool   { return w.s.workTime == 0 }

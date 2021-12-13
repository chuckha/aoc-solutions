package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	subs := map[string][]string{}
	reverseSubs := []*substitution{}
	for _, line := range lines[:len(lines)-1] {
		words := strings.Split(line, " => ")
		if _, ok := subs[words[0]]; ok {
			subs[words[0]] = append(subs[words[0]], words[1])
		} else {
			subs[words[0]] = []string{words[1]}
		}
		if words[0] == "e" {
			continue
		}
		reverseSubs = append(reverseSubs, &substitution{words[1], words[0]})
	}
	// sort.Sort(sort.Reverse(subst(reverseSubs)))
	// fmt.Println(reverseSubs)
	input := lines[len(lines)-1]
	count := 0
	internal.Shuffle(reverseSubs)
	newInput := ""

	for newInput != "e" {
		for _, r := range reverseSubs {
			for {
				newInput = strings.Replace(input, r.from, r.to, 1)
				if newInput == input {
					break
				}
				count++
				input = newInput
			}
			fmt.Println(input, count)
		}
	}
	// pieces := strings.Split(input, "Ar")
	// for _, piece := range pieces {
	// 	fullPiece := piece + "Ar"
	// 	doit(fullPiece, reverseSubs)
	// }
	// thing(input, subs)
	//doit(input, reverseSubs)
	// for key, replacements := range subs {
	// 	//		fmt.Println("===", key, "===")
	// 	parts := strings.Split(input, key)
	// 	for i := 0; i < len(parts)-1; i++ {
	// 		for _, replacement := range replacements {
	// 			//				fmt.Println("***", replacement, "***")
	// 			muck := make([]string, len(parts))
	// 			copy(muck, parts)
	// 			//				fmt.Println(muck, len(muck))
	// 			muck[i] = muck[i] + replacement + muck[i+1]
	// 			muck = append(muck[:i+1], muck[i+2:]...)
	// 			r := strings.Join(muck, key)
	// 			//				fmt.Println(r)
	// 			if _, ok := substitutionSet[r]; !ok {
	// 				substitutionSet[r] = struct{}{}
	// 			}
	// 		}
	// 	}
	// }
	// //	fmt.Println(substitutionSet)
	// fmt.Println(len(substitutionSet))
}

func doit(item string, subs []*substitution) {
	/*
	   for each substitution, reverse substitute
	   if you run into a situation where you can't go backwards...?
	*/
	q := internal.NewQueue[*state]()
	q.Enqueue(&state{item, 0})
	for !q.Empty() {
		w := q.Dequeue()
		fmt.Println("number of substitutions:", w.numberOfSubstitutions)
		fmt.Println(w.data)
		// for reach substitutoin, run one and add it to the queue
		for _, sub := range subs {
			replaced := strings.Replace(w.data, sub.from, sub.to, 1)
			if len(replaced) == 1 {
				fmt.Println("solved")
				return
			}
			if replaced == w.data {
				continue
			}

			// fmt.Println("replacing", sub.from, "with", sub.to)
			// fmt.Println("result")
			// fmt.Println(replaced)
			//			time.Sleep(100 * time.Millisecond)
			// some analytics
			q.Enqueue(&state{replaced, w.numberOfSubstitutions + 1})
		}
	}
}

type state struct {
	data                  string
	numberOfSubstitutions int
}

type nodeData struct {
	data        string
	prefixCount int
}

func (n *nodeData) String() string {
	return n.data
}

func thing(goal string, subs map[string][]string) {
	rootNode := &nodeData{"e", 0}
	t := internal.NewTree(rootNode, 0)
	q := internal.NewQueue[*internal.Tree[*nodeData]]()
	q.Enqueue(t)
	goalParts := splitOnUpperCase(goal)
	count := 0
	for !q.Empty() {
		if count == 1000 {
			fmt.Println(t)
			return
		}
		count++
		node := q.Dequeue()
		if node.Data.data == goal {
			fmt.Println("DONE", node.Data)
			os.Exit(1)
		}
		els := splitOnUpperCase(node.Data.data)
		// fmt.Println("node data", node.Data, els)
		for i, el := range els {
			// fmt.Println("el", el)
			// fmt.Println("subs", subs[el])
			// fmt.Println("els", els)
			for _, sub := range subs[el] {
				substitution := strings.Join(els[:i], "") + sub + strings.Join(els[i+1:], "")
				subsPart := splitOnUpperCase(substitution)
				pfxCount := prefixCount(goalParts, subsPart)
				// don't go backwards (this might result in not finding a minimum)
				// if pfxCount <= node.Data.prefixCount && node.Depth > 1 {
				// 	continue
				// }

				child := internal.NewTree(&nodeData{substitution, pfxCount}, node.Depth+1)
				node.Children = append(node.Children, child)
				if len(substitution) > len(goal) {
					continue
				}
				if strings.HasSuffix(substitution, "Ar") {
					continue
				}
				if strings.HasPrefix(substitution, "CRnTh") {
					continue
				}
				if strings.HasPrefix(substitution, "CRnAl") {
					continue
				}
				q.Enqueue(child)
			}
		}
	}
	fmt.Println(t)
}

func splitOnUpperCase(in string) []string {
	out := []string{}
	hodl := ""
	for _, c := range in {
		if c >= 56 && c <= 90 && len(hodl) == 1 {
			out = append(out, hodl)
			hodl = string(c)
			continue
		}
		hodl += string(c)
		if c >= 97 && c <= 122 {
			out = append(out, hodl)
			hodl = ""
		}
	}
	if hodl != "" {
		out = append(out, hodl)
	}
	return out
}

func prefixCount(goal, sub []string) int {
	for i, g := range goal {
		if sub[i] == g {
			continue
		}
		return i
	}
	return -1
}

type subst []*substitution

func (s subst) Len() int {
	return len(s)
}
func (s subst) Less(i, j int) bool {
	return len(s[i].from) < len(s[j].from)
}

// Swap swaps the elements with indexes i and j.
func (s subst) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type substitution struct {
	from string
	to   string
}

func (s *substitution) String() string {
	return fmt.Sprintf("%s->%s", s.from, s.to)
}

// e
// H, O
// HO, OH, HH
// HOO, OHO, HHH, OOH, HOH, OHH, HHO, HOH
/*

       e
   H        O




*/

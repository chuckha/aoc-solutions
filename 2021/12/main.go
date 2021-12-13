package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	graph := map[string][]string{}
	for _, line := range lines {
		nodes := strings.Split(line, "-")
		if _, ok := graph[nodes[0]]; ok {
			graph[nodes[0]] = append(graph[nodes[0]], nodes[1])
		} else {
			graph[nodes[0]] = []string{nodes[1]}
		}
		if _, ok := graph[nodes[1]]; ok {
			graph[nodes[1]] = append(graph[nodes[1]], nodes[0])
		} else {
			graph[nodes[1]] = []string{nodes[0]}
		}
	}
	fmt.Println(len(bfs(graph)))
	// for _, path := range bfs(graph) {
	// 	fmt.Println(strings.Join(path, ","))
	// }
}

func bfs(graph map[string][]string) [][]string {
	q := internal.NewQueue[*state]()
	q.Enqueue(&state{
		node: "start",
		path: []string{"start"},
	})
	paths := [][]string{}
	for !q.Empty() {
		s := q.Dequeue()
		for _, way := range graph[s.node] {
			hdv := s.hasDoubleVisit

			if way == "start" {
				continue
			}
			if isSmallCave(way) && internal.Search(way, s.path) != -1 {
				if s.hasDoubleVisit {
					continue
				}
				hdv = true
			}

			newPath := make([]string, len(s.path))
			copy(newPath, s.path)
			newPath = append(newPath, way)
			if way == "end" {
				paths = append(paths, newPath)
				continue
			}
			//			fmt.Println("enqueueing", way, newPath, hdv)
			q.Enqueue(&state{
				node:           way,
				path:           newPath,
				hasDoubleVisit: hdv,
			})
		}
	}
	return paths
}

type state struct {
	node           string
	path           []string
	hasDoubleVisit bool
}

func isSmallCave(name string) bool {
	return name[0] >= 97 && name[0] <= 122
}

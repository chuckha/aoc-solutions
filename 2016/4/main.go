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
	s := 0
	valid := []*room{}
	for _, line := range lines {
		parts := strings.Split(line, "-")
		lastIndex := findLastIndex(parts)
		lastPart := strings.Split(parts[lastIndex], "[")
		id, _ := strconv.Atoi(lastPart[0])
		r := &room{
			code:             strings.Join(parts[:lastIndex], ""),
			id:               id,
			providedChecksum: strings.TrimSuffix(lastPart[1], "]"),
			name:             strings.Join(parts[:lastIndex], " "),
		}
		if r.checksum() == r.providedChecksum {
			s += r.id
		}
		valid = append(valid, r)
	}
	for _, room := range valid {
		fmt.Println(room.uncipher(), room.id)
	}
}

func findLastIndex(parts []string) int {
	for i, part := range parts {
		if part[0] >= 48 && part[0] <= 57 {
			return i
		}
	}
	return -1
}

type room struct {
	code             string
	id               int
	name             string
	providedChecksum string
}

func (r *room) String() string {
	return fmt.Sprintf("%s %s vs %s", r.code, r.providedChecksum, r.checksum())
}

type count struct {
	letter string
	count  int
}

func (r *room) checksum() string {
	counts := map[string]*count{}
	for _, l := range r.code {
		if c, ok := counts[string(l)]; ok {
			c.count++
			continue
		}
		counts[string(l)] = &count{letter: string(l), count: 1}
	}
	cs := make([]*count, 0)
	for _, v := range counts {
		cs = append(cs, v)
	}
	sort.Sort(sort.Reverse(cSlice(cs)))
	out := ""
	for _, c := range cs[:5] {
		out += c.letter
	}
	return out
}

func (r *room) uncipher() string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	readableName := ""
	for _, letter := range r.name {
		if string(letter) == " " {
			readableName += " "
			continue
		}
		idx := strings.Index(letters, string(letter))
		idx += r.id
		idx %= 26
		readableName += string(letters[idx])
	}
	return readableName
}

type cSlice []*count

func (c cSlice) Len() int { return len(c) }
func (c cSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c cSlice) Less(i, j int) bool {
	if c[i].count < c[j].count {
		return true
	}
	if c[i].count == c[j].count {
		return c[i].letter > c[j].letter
	}
	return false
}

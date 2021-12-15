package main

import (
	"fmt"
	"strings"

	"github.com/chuckha/aoc-solutions/internal"
)

func main() {
	lines := internal.ReadInput()
	count := 0
	for _, line := range lines {
		if supportsSSL(line) {
			count++
		}
	}
	fmt.Println(count)
}

func supportsSSL(line string) bool {
	s := &state2{}
	b := bracketedaba{
		ignoreInput: true,
		state2:      &state2{},
	}
	for _, c := range line {
		s.addLetter(string(c))
		b.addLetter(string(c))
	}
	return hasAbaAndBab(s.abas, b.state2.abas)
}

func findABBA(line string) bool {
	s := &state{}
	b := &bracketedAbba{
		state: &state{},
	}
	for _, c := range line {
		s.addLetter(string(c))
		b.addLetter(string(c))
		//		fmt.Println(s.data)
	}
	fmt.Println(s.found, !b.found)
	return s.found && !b.found
}

type bracketedAbba struct {
	found             bool
	possiblyInBracket bool
	*state
}

func (b *bracketedAbba) addLetter(letter string) {
	if b.found {
		return
	}
	if letter == "[" {
		b.possiblyInBracket = true
		return
	}
	if letter == "]" && b.possiblyInBracket {
		b.possiblyInBracket = false
		return
	}
	if b.possiblyInBracket {
		b.state.addLetter(letter)
	}
	b.found = b.state.found
}

type state struct {
	data  []string
	found bool
}

func (s *state) addLetter(letter string) {
	if s.found {
		return
	}
	if letter == "[" || letter == "]" {
		s.data = make([]string, 0)
		return
	}
	s.data = append(s.data, letter)
	if !canBe(s.data) {
		s.data = s.data[1:]
	}
	if len(s.data) == 4 {
		s.found = abba(s.data)
	}
}

func abba(in []string) bool {
	return in[0] == in[3] && in[1] == in[2] && in[0] != in[1]
}

func canBe(in []string) bool {
	switch len(in) {
	case 0:
		return true
	case 1:
		return true
	case 2:
		return in[0] != in[1]
	case 3:
		return in[0] != in[1] && in[1] == in[2]
	case 4:
		return abba(in)
	default:
		panic("very bad")
	}
}

/*
	find all aba outside of brackets
	find all aba inside brackets
	for each aba outside, is there a cooresponding bab inside?

	zazbz[bzb]cdb
	zaz, zbz || bzb


*/
type bracketedaba struct {
	ignoreInput bool
	*state2
}

func (b *bracketedaba) addLetter(letter string) {
	if letter == "[" {
		b.ignoreInput = false
		return
	}
	if letter == "]" {
		b.ignoreInput = true
	}
	if b.ignoreInput {
		return
	}
	b.state2.addLetter(letter)
}

type state2 struct {
	abas        []string
	data        []string
	ignoreInput bool
}

func (s *state2) addLetter(letter string) {
	if letter == "]" {
		s.ignoreInput = false
	}
	if s.ignoreInput {
		return
	}
	if letter == "[" {
		s.ignoreInput = true
		s.data = make([]string, 0)
		return
	}
	s.data = append(s.data, letter)
	if !canBeAba(s.data) {
		s.data = s.data[1:]
	}
	if len(s.data) == 3 && aba(s.data) {
		s.abas = append(s.abas, strings.Join(s.data, ""))
		s.data = s.data[1:]
	}
}

func aba(in []string) bool {
	return in[0] == in[2] && in[1] != in[0]
}

func canBeAba(in []string) bool {
	switch len(in) {
	case 0:
		return true
	case 1:
		return true
	case 2:
		return in[0] != in[1]
	case 3:
		return aba(in)
	default:
		panic("very not good")
	}
}

func hasAbaAndBab(abas, babs []string) bool {
	for _, aba := range abas {
		for _, bab := range babs {
			if aba[0] == bab[1] && aba[1] == bab[0] {
				return true
			}
		}
	}
	return false
}

/*
aba[bab]xyz
xyx[xyx]xyx
aaa[kek]eke
zazbz[bzb]cdb
*/

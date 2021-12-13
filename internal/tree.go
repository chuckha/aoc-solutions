package internal

import (
	"fmt"
	"strings"
)

type Tree[T any] struct {
	Data     T
	Depth    int
	Children []*Tree[T]
}

func NewTree[T any](data T, depth int) *Tree[T] {
	return &Tree[T]{
		Data:     data,
		Depth:    depth,
		Children: make([]*Tree[T], 0),
	}
}

func (t *Tree[T]) AddChild(data T) {
	t.Children = append(t.Children, NewTree(data, t.Depth+1))
}

func (t *Tree[T]) String() string {
	out := fmt.Sprintf("%v\n", t.Data)
	for _, c := range t.Children {
		out += fmt.Sprintf("%s%v\n", strings.Repeat("\t", t.Depth), c)
	}
	return out
}

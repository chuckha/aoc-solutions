package internal

import (
	"fmt"
	"strings"
)

type Node[T any] struct {
	Next *Node[T]
	Prev *Node[T]
	Data T
}

func NewNode[T any](data T) *Node[T] {
	c := &Node[T]{
		Data: data,
	}
	c.Next = c
	c.Prev = c
	return c
}

func (n *Node[T]) Forward(num int) *Node[T] {
	cur := n
	for i := 0; i < num; i++ {
		cur = cur.Next
	}
	return cur
}

func (n *Node[T]) Len() int {
	count := 1
	root := n
	cur := n
	for cur.Next != root {
		cur = cur.Next
		count++
	}
	return count
}

func (n *Node[T]) Insert(item T) *Node[T] {
	newNode := NewNode(item)
	next := n.Next
	n.Next = newNode
	newNode.Prev = n
	newNode.Next = next
	next.Prev = newNode
	return newNode
}

func (n *Node[T]) Remove() *Node[T] {
	next := n.Next
	prev := n.Prev
	prev.Next = next
	next.Prev = prev
	return prev
}

func (n *Node[T]) String() string {
	root := n
	cur := n
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%v -> ", cur.Data))
	cur = cur.Next
	for cur != root {
		out.WriteString(fmt.Sprintf("%v -> ", cur.Data))
		cur = cur.Next
	}
	return out.String()
}

func (n *Node[T]) Swap(i, j int) {
	root := n
	first := root.Forward(i)
	second := root.Forward(j)
	firstPrev := first.Remove()
	firstPrev.Insert(second.Data)
	secondPrev := second.Remove()
	secondPrev.Insert(first.Data)
}

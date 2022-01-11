package internal

import (
	"fmt"
	"strings"
)

type CircularLinkedList[T any] struct {
	Prev *CircularLinkedList[T]
	Next *CircularLinkedList[T]
	Data T
}

func NewCircularLinkedList[T any](data T) *CircularLinkedList[T] {
	c := &CircularLinkedList[T]{
		Data: data,
	}
	c.Prev = c
	c.Next = c
	return c
}

func (c *CircularLinkedList[T]) Forward(n int) *CircularLinkedList[T] {
	next := c
	for i := 0; i < n; i++ {
		next = next.Next
	}
	return next
}

func (c *CircularLinkedList[T]) InsertAfter(data T) *CircularLinkedList[T] {
	n := NewCircularLinkedList(data)
	n.Next = c.Next
	c.Next.Prev = n
	n.Prev = c
	c.Next = n
	if c.Prev == c {
		c.Prev = n
	}
	return n
}

// 0 <- 0 -> 0
//
// 0 <- 1 -> 0
// 1 <- 0 -> 1
//
// 2
// 1 <- 0 -> 1
// 0 <- 1 -> 0
// 1 <- 2 -> 0

// 0 -> 1 -> 2 -> 3 -> 4
// c = 2
// prev = 1
// next = 3
// 1 -> 3
// 1 <- 3

func (c *CircularLinkedList[T]) Remove() *CircularLinkedList[T] {
	prev := c.Prev
	next := c.Next
	prev.Next = next
	next.Prev = prev
	return next
}

func (c *CircularLinkedList[T]) String() string {
	cur := c
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%v", cur.Data))
	for cur.Next != c {
		cur = cur.Next
		out.WriteString(fmt.Sprintf("->%v", cur.Data))
	}
	return out.String()
}

func (c *CircularLinkedList[T]) PrintFormat(c1, c2 *CircularLinkedList[T]) string {
	cur := c
	var out strings.Builder
	if cur == c1 {
		out.WriteString(fmt.Sprintf("(%v)", cur.Data))
	} else if cur == c2 {
		out.WriteString(fmt.Sprintf("[%v]", cur.Data))
	} else {
		out.WriteString(fmt.Sprintf("%v", cur.Data))
	}
	for cur.Next != c {
		cur = cur.Next
		if cur == c1 {
			out.WriteString(fmt.Sprintf("->(%v)", cur.Data))
			continue
		}
		if cur == c2 {
			out.WriteString(fmt.Sprintf("->[%v]", cur.Data))
			continue
		}
		out.WriteString(fmt.Sprintf("->%v", cur.Data))
	}
	return out.String()

}

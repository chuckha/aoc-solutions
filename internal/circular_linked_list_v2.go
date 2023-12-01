package internal

import (
	"fmt"
	"strings"
)

type sortable[T any] interface {
	Equal(a T) bool
	Data() T
}

type CircularLinkedListV2[T sortable[T]] struct {
	Prev *CircularLinkedListV2[T]
	Next *CircularLinkedListV2[T]
	Data T
}

func NewCircularLinkedListV2[T sortable[T]](data T) *CircularLinkedListV2[T] {
	c := &CircularLinkedListV2[T]{
		Data: data,
	}
	c.Prev = c
	c.Next = c
	return c
}

func (c *CircularLinkedListV2[T]) Forward(n int) *CircularLinkedListV2[T] {
	next := c
	for i := 0; i < n; i++ {
		next = next.Next
	}
	return next
}

func (c *CircularLinkedListV2[T]) Backwards(n int) *CircularLinkedListV2[T] {
	prev := c
	for i := 0; i < n; i++ {
		prev = prev.Prev
	}
	return prev
}

func (c *CircularLinkedListV2[T]) InsertAfter(data T) *CircularLinkedListV2[T] {
	n := NewCircularLinkedListV2(data)
	n.Next = c.Next
	c.Next.Prev = n
	n.Prev = c
	c.Next = n
	if c.Prev == c {
		c.Prev = n
	}
	return n
}

func (c *CircularLinkedListV2[T]) InsertBefore(data T) *CircularLinkedListV2[T] {
	n := NewCircularLinkedListV2(data)
	n.Next = c
	c.Prev.Next = n
	n.Prev = c.Prev
	c.Prev = n
	if c.Next == c {
		c.Next = n
	}
	return n
}

func (c *CircularLinkedListV2[T]) Find(x T) *CircularLinkedListV2[T] {
	for ; c.Next != c; c = c.Next {
		if c.Data.Equal(x.Data()) {
			return c
		}
	}
	return nil
}

func (c *CircularLinkedListV2[T]) Remove() *CircularLinkedListV2[T] {
	prev := c.Prev
	next := c.Next
	prev.Next = next
	next.Prev = prev
	c = nil
	return next
}

func (c *CircularLinkedListV2[T]) String() string {
	cur := c
	var out strings.Builder
	out.WriteString(fmt.Sprintf("%v", cur.Data))
	for cur.Next != c {
		cur = cur.Next
		out.WriteString(fmt.Sprintf("->%v", cur.Data))
	}
	return out.String()
}

func (c *CircularLinkedListV2[T]) Verify() bool {
	cur := c
	for cur.Next != c {
		if cur != cur.Next.Prev {
			return false
		}
		if cur != cur.Prev.Next {
			return false
		}
		cur = cur.Next
	}
	return true
}

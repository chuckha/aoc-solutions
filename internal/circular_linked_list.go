package internal

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

type CircularLinkedList[T constraints.Ordered] struct {
	Prev *CircularLinkedList[T]
	Next *CircularLinkedList[T]
	Data T
}

func NewCircularLinkedList[T constraints.Ordered](data T) *CircularLinkedList[T] {
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

func (c *CircularLinkedList[T]) Backwards(n int) *CircularLinkedList[T] {
	prev := c
	for i := 0; i < n; i++ {
		prev = prev.Prev
	}
	return prev
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

func (c *CircularLinkedList[T]) InsertBefore(data T) *CircularLinkedList[T] {
	n := NewCircularLinkedList(data)
	n.Next = c
	c.Prev.Next = n
	n.Prev = c.Prev
	c.Prev = n
	if c.Next == c {
		c.Next = n
	}
	return n
}

func (c *CircularLinkedList[T]) Find(x T) *CircularLinkedList[T] {
	for ; c.Next != c; c = c.Next {
		if c.Data == x {
			return c
		}
	}
	return nil
}

func (c *CircularLinkedList[T]) Remove() *CircularLinkedList[T] {
	prev := c.Prev
	next := c.Next
	prev.Next = next
	next.Prev = prev
	c = nil
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

func (c *CircularLinkedList[T]) Verify() bool {
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

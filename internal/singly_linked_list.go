package internal

import (
	"fmt"
	"strings"
)

type LinkedList[T comparable] struct {
	Data T
	Next *LinkedList[T]
}

func NewLL[T comparable](data T) *LinkedList[T] {
	return &LinkedList[T]{
		Data: data,
	}
}
func (l *LinkedList[T]) Forward(n int) *LinkedList[T] {
	for i := 0; i < n; i++ {
		l = l.Next
	}
	return l
}

func (l *LinkedList[T]) Insert(data T) *LinkedList[T] {
	tmp := l.Next
	l.Next = NewLL(data)
	l.Next.Next = tmp
	return l.Next
}

func (l *LinkedList[T]) String() string {
	var out strings.Builder
	for l != nil {
		out.WriteString(fmt.Sprintf("%v->", l.Data))
		l = l.Next
	}
	return out.String()
}

func (l *LinkedList[T]) Reverse() *LinkedList[T] {
	cur := l.Next
	l.Next = nil
	prev := l
	tmp := cur.Next
	for tmp != nil {
		cur.Next = prev
		prev = cur
		cur = tmp
		tmp = cur.Next
	}
	cur.Next = prev
	return cur
}

func (l *LinkedList[T]) Split(start, end int) (*LinkedList[T], *LinkedList[T], *LinkedList[T]) {
	startNode := l.Forward(start)
	endNode := l.Forward(end)
	l.Next = nil
	startTail := endNode.Next
	endNode.Next = nil
	return l, startNode, startTail
}

func (l *LinkedList[T]) Len() int {
	count := 0
	for l != nil {
		l = l.Next
		count++
	}
	return count
}

func (l *LinkedList[T]) Find(item T) *LinkedList[T] {
	for l != nil && l.Data != item {
		l = l.Next
	}
	return l
}

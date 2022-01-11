package internal

import (
	"fmt"
	"strings"
)

type DLL[T any] struct {
	Next *DLL[T]
	Prev *DLL[T]
	Data T
}

func NewDLL[T any](data T) *DLL[T] {
	return &DLL[T]{
		Data: data,
	}
}
func (n *DLL[T]) InsertAfter(d T) *DLL[T] {
	newNode := NewDLL(d)
	next := n.Next
	if next != nil {
		next.Prev = newNode
	}
	newNode.Next = next
	newNode.Prev = n
	n.Next = newNode
	return newNode
}
func (n *DLL[T]) InsertBefore(d T) *DLL[T] {
	newDLL := NewDLL(d)
	newDLL.Next = n
	prev := n.Prev
	newDLL.Prev = prev
	if prev != nil {
		prev.Next = newDLL
	}
	n.Prev = newDLL
	return newDLL
}

func (d *DLL[T]) Remove() *DLL[T] {
	next := d.Next
	prev := d.Prev
	if prev != nil {
		prev.Next = next
	}
	if next != nil {
		next.Prev = prev
	}
	return prev

}
func (d *DLL[T]) String() string {
	cur := d
	var out strings.Builder
	for cur != nil {
		out.WriteString(fmt.Sprintf("%v -> ", cur.Data))
		cur = cur.Next
	}
	return out.String()
}
func (d *DLL[T]) Len() int {
	out := 0
	for d != nil {
		out++
		d = d.Next
	}
	return out
}

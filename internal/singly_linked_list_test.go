package internal

import (
	"fmt"
	"testing"
)

func TestSinglyLinkedList(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		root := NewLL(0)
		root.Insert(1).Insert(2)
		fmt.Println(root)
		root.Insert(3)
		fmt.Println(root)
	})
	t.Run("reverse should work", func(t *testing.T) {
		root := NewLL(0)
		root.Insert(1).Insert(2).Insert(3).Insert(4)
		root = root.Reverse()
		fmt.Println(root)
	})
	t.Run("forward", func(t *testing.T) {
		root := NewLL(0)
		root.Insert(1).Insert(2).Insert(3).Insert(4)
		out := root.Forward(2)
		if out.Data != 2 {
			t.Fatal("expected 2 but got", out.Data)
		}
	})
	t.Run("len", func(t *testing.T) {
		root := NewLL(0)
		root.Insert(1).Insert(2).Insert(3).Insert(4)
		if root.Len() != 5 {
			t.Fatal("expected 5 but got", root.Len())
		}
	})
	t.Run("Split", func(t *testing.T) {
		root := NewLL(0)
		root.Insert(1).Insert(2).Insert(3).Insert(4)
		a, b, c := root.Split(1, 3)
		fmt.Println("head", a)
		fmt.Println("body", b)
		fmt.Println("tail", c)
	})
	t.Run("Find", func(t *testing.T) {
		root := NewLL(0)
		root.Insert(1).Insert(2).Insert(3).Insert(4)
		if root.Find(3) == nil {
			t.Fatal("did not find the right guy")
		}
		if root.Find(3).Data != 3 {
			t.Fatal("real weird")
		}
	})
}

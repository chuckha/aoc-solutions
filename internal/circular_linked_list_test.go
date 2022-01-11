package internal

import "testing"

func TestCircularLinkedList(t *testing.T) {
	t.Run("test insert after", func(t *testing.T) {
		c := NewCircularLinkedList(1)
		c2 := c.InsertAfter(2)
		if c.Next != c2 {
			t.Fatal("c next is not c2")
		}
		if c.Prev != c2 {
			t.Fatal("c prev is not c2")
		}
		if c2.Next != c {
			t.Fatal("c2 next is not c")
		}
		if c2.Prev != c {
			t.Fatal("c2 prev is not c")
		}
	})
	t.Run("remove a node", func(t *testing.T) {
		c := NewCircularLinkedList(0)
		c.InsertAfter(1).InsertAfter(2).InsertAfter(3).InsertAfter(4).InsertAfter(5)
		c = c.Remove()
		if c.Data != 1 {
			t.Fatal("did not remove a node properly")
		}
	})
	t.Run("print", func(t *testing.T) {
		c := NewCircularLinkedList(0)
		c.Next.InsertAfter(6).InsertAfter(9).InsertAfter(2)
		if c.String() != "0->6->9->2" {
			t.Fatal("expected 0 -> 6 -> 9 -> 2 but got" + c.String())
		}
	})
}

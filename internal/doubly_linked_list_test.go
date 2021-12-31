package internal

import (
	"fmt"
	"testing"
)

func TestLinkedList(t *testing.T) {
	root := NewNode(0)
	cur := root
	for i := 1; i < 100; i++ {
		cur = cur.Insert(i)
	}
	length := cur.Len()
	if length != 100 {
		t.Fatal("length is wrong", length, 100)
	}

	t.Run("swap should swap two elements", func(t *testing.T) {
		root := NewNode(0)
		root.Insert(1).Insert(2).Insert(3).Insert(4)
		root.Swap(2, 4)
		fmt.Println(root)
	})

}

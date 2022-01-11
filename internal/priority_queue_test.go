package internal

import "testing"

func TestPriorityQueue(t *testing.T) {
	t.Run("should extract cheapest item on pull", func(t *testing.T) {
		pq := NewPriorityQueue[int]()
		for i := 0; i < 10; i++ {
			pq.Add(i, i)
		}
		if pq.Size() != 10 {
			t.Fatal("pq size is wrong")
		}
		if pq.Pull() != 0 {
			t.Fatal("expected 0 but got something else")
		}
		if pq.Size() != 9 {
			t.Fatal("pull did not remove an item")
		}
	})
}
